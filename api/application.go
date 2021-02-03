package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/jinzhu/copier"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

const (
	ConfigDir                   = "/etc/baetyl"
	ProgramConfigDir            = "/var/lib/baetyl/bin"
	FunctionConfigPrefix        = "baetyl-function-config"
	FunctionProgramConfigPrefix = "baetyl-function-program-config"
	FunctionCodePrefix          = "baetyl-function-code"
	FunctionDefaultConfigFile   = "conf.yml"
)

// GetApplication get a application
func (api *API) GetApplication(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	app, err := api.App.Get(ns, n, "")
	if err != nil {
		return nil, err
	}

	// sys app: core、init、function is not visible
	if common.ValidIsInvisible(app.Labels) {
		return nil, common.Error(common.ErrResourceInvisible, common.Field("type", common.APP), common.Field("name", app.Name))
	}
	return api.ToApplicationView(app)
}

// ListApplication list application
func (api *API) ListApplication(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	params, err := api.parseListOptionsAppendSystemLabel(c)
	if err != nil {
		return nil, err
	}
	apps, err := api.App.List(ns, params)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	return apps, err
}

// CreateApplication create one application
func (api *API) CreateApplication(c *common.Context) (interface{}, error) {
	appView, err := api.parseApplication(c)
	if err != nil {
		return nil, err
	}
	ns, name := c.GetNamespace(), appView.Name

	err = api.validApplication(ns, appView)
	if err != nil {
		return nil, err
	}

	// TODO: remove get method, return error inside service instead
	oldApp, err := api.App.Get(ns, name, "")
	if err != nil {
		if e, ok := err.(errors.Coder); !ok || e.Code() != common.ErrResourceNotFound {
			return nil, err
		}
	}
	if oldApp != nil {
		return nil, common.Error(common.ErrResourceHasBeenUsed,
			common.Field("error", "this name is already in use"))
	}

	baseApp, err := api.getBaseAppIfSet(c)
	if err != nil {
		return nil, err
	}
	if baseApp != nil && baseApp.Type != appView.Type {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "the type of baseApp is conflicted"))
	}

	app, configs, err := api.ToApplication(appView, nil)
	if err != nil {
		return nil, err
	}

	err = api.updateGenConfigsOfFunctionApp(ns, configs)
	if err != nil {
		return nil, err
	}

	app, err = api.App.CreateWithBase(ns, app, baseApp)
	if err != nil {
		return nil, err
	}

	err = api.UpdateNodeAndAppIndex(ns, app)
	if err != nil {
		return nil, err
	}

	return api.ToApplicationView(app)
}

// UpdateApplication update the application
func (api *API) UpdateApplication(c *common.Context) (interface{}, error) {
	appView, err := api.parseApplication(c)
	if err != nil {
		return nil, err
	}

	ns, name := c.GetNamespace(), c.GetNameFromParam()

	err = api.validApplication(ns, appView)
	if err != nil {
		return nil, err
	}

	oldApp, err := api.App.Get(ns, name, "")
	if err != nil {
		return nil, err
	}

	// sys app: core、init、function is not visible
	if common.ValidIsInvisible(oldApp.Labels) {
		return nil, common.Error(common.ErrResourceInvisible, common.Field("type", common.APP), common.Field("name", oldApp.Name))
	}

	// labels and Selector can't be modified of sys apps
	if checkIsSysResources(oldApp.Labels) &&
		(oldApp.Selector != appView.Selector || !reflect.DeepEqual(oldApp.Labels, appView.Labels) || !appView.System) {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "selector，labels or system field can't be modified of sys apps"))
	}

	appView.Version = oldApp.Version
	app, configs, err := api.ToApplication(appView, oldApp)
	if err != nil {
		return nil, err
	}

	err = api.updateGenConfigsOfFunctionApp(ns, configs)
	if err != nil {
		return nil, err
	}

	app, err = api.App.Update(ns, app)
	if err != nil {
		return nil, err
	}

	if oldApp != nil && oldApp.Selector != app.Selector {
		// delete old nodes
		if err := api.DeleteNodeAndAppIndex(ns, oldApp); err != nil {
			return nil, err
		}
	}

	// update nodes
	if err := api.UpdateNodeAndAppIndex(ns, app); err != nil {
		return nil, err
	}

	api.cleanGenConfigsOfFunctionApp(configs, oldApp)

	return api.ToApplicationView(app)
}

// DeleteApplication delete the application
func (api *API) DeleteApplication(c *common.Context) (interface{}, error) {
	ns, name := c.GetNamespace(), c.GetNameFromParam()
	app, err := api.App.Get(ns, name, "")
	if err != nil {
		if e, ok := err.(errors.Coder); ok && e.Code() == common.ErrResourceNotFound {
			return nil, nil
		}
		return nil, err
	}

	if canDelete, err := api.isAppCanDelete(ns, name); err != nil {
		return nil, err
	} else if !canDelete {
		return nil, common.Error(common.ErrAppReferencedByNode, common.Field("name", name))
	}

	if err := api.App.Delete(ns, c.GetNameFromParam(), ""); err != nil {
		return nil, err
	}

	//delete the app from node
	if err := api.DeleteNodeAndAppIndex(ns, app); err != nil {
		return nil, err
	}

	api.cleanGenConfigsOfFunctionApp(nil, app)
	return nil, nil
}

func (api *API) GetSysAppConfigs(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	_, err := api.App.Get(ns, n, "")
	if err != nil {
		return nil, err
	}

	ops := &models.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", BaetylAppNameKey, n),
	}

	list, err := api.Config.List(ns, ops)
	if err != nil {
		log.L().Error("failed to list configs", log.Error(err))
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	for i := range list.Items {
		list.Items[i].Data = nil
	}
	return list, err
}

func (api *API) GetSysAppSecrets(c *common.Context) (interface{}, error) {
	res, err := api.getNodeSysAppSecretLikedResources(c)
	if err != nil {
		return nil, err
	}
	return api.ToSecretViewList(res), nil
}

func (api *API) GetSysAppCertificates(c *common.Context) (interface{}, error) {
	res, err := api.getNodeSysAppSecretLikedResources(c)
	if err != nil {
		return nil, err
	}

	return api.ToCertificateViewList(res), nil
}

func (api *API) GetSysAppRegistries(c *common.Context) (interface{}, error) {
	res, err := api.getNodeSysAppSecretLikedResources(c)
	if err != nil {
		return nil, err
	}
	return api.ToRegistryViewList(res), nil
}

func (api *API) parseApplication(c *common.Context) (*models.ApplicationView, error) {
	app := new(models.ApplicationView)
	app.Name = c.GetNameFromParam()
	app.Namespace = c.GetNamespace()
	err := c.LoadBody(app)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	if name := c.GetNameFromParam(); name != "" {
		app.Name = name
	}
	if app.Name == "" {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "name is required"))
	}

	if app.Type == common.ContainerApp {
		for _, v := range app.Services {
			if v.FunctionConfig != nil || v.Functions != nil {
				return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "add function info in container app"))
			}
			if v.Image == "" {
				return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "image is required in container app"))
			}
		}
	} else if app.Type == common.FunctionApp {
		for _, v := range app.Services {
			if v.FunctionConfig == nil {
				return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "function config can't be empty in function app"))
			}
		}
		if len(app.Registries) != 0 {
			return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "registries should be be empty in function app"))
		}
	} else {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "type is invalid"))
	}
	return app, nil
}

func (api *API) getBaseAppIfSet(c *common.Context) (*specV1.Application, error) {
	if base, ok := c.GetQuery("base"); ok {
		namespace := c.GetNamespace()
		baseApp, err := api.App.Get(namespace, base, "")
		if err != nil {
			return nil, err
		}
		return baseApp, nil
	}
	return nil, nil
}

func (api *API) parseListOptions(c *common.Context) (*models.ListOptions, error) {
	params := &models.ListOptions{}
	if err := c.Bind(params); err != nil {
		return nil, err
	}
	return params, nil
}

func (api *API) parseListOptionsAppendSystemLabel(c *common.Context) (*models.ListOptions, error) {
	opt, err := api.parseListOptions(c)
	if err != nil {
		return nil, err
	}

	ls := opt.LabelSelector
	if !strings.Contains(ls, common.LabelSystem) {
		if len(strings.TrimSpace(ls)) > 0 {
			ls += ","
		}
		ls += "!" + common.LabelSystem
	}

	opt.LabelSelector = ls
	return opt, nil
}

func (api *API) UpdateNodeAndAppIndex(namespace string, app *specV1.Application) error {
	nodes, err := api.Node.UpdateNodeAppVersion(namespace, app)
	if err != nil {
		return err
	}
	return api.Index.RefreshNodesIndexByApp(namespace, app.Name, nodes)
}

func (api *API) DeleteNodeAndAppIndex(namespace string, app *specV1.Application) error {
	_, err := api.Node.DeleteNodeAppVersion(namespace, app)
	if err != nil {
		return err
	}

	return api.Index.RefreshNodesIndexByApp(namespace, app.Name, make([]string, 0))
}

func (api *API) ToApplicationView(app *specV1.Application) (*models.ApplicationView, error) {
	appView := &models.ApplicationView{}
	copier.Copy(appView, app)

	err := api.translateSecretsToSecretLikedResources(appView)
	if err != nil {
		return nil, err
	}

	if app.Type != common.FunctionApp {
		return appView, nil
	}
	for index := range appView.Services {
		service := &appView.Services[index]
		generatedConfigName, err := getGenConfigNameOfFunctionService(app, service.Name)
		if err != nil {
			return nil, err
		}

		generatedProgramConfigName, err := getGenProgramNameOfFunctionService(app, service.Name)
		if err != nil {
			return nil, err
		}

		config, err := api.Config.Get(appView.Namespace, generatedConfigName, "")
		if err != nil {
			return nil, err
		}
		if data, ok := config.Data[FunctionDefaultConfigFile]; ok {
			serviceFunctions := new(models.ServiceFunction)
			err := json.Unmarshal([]byte(data), serviceFunctions)
			if err != nil {
				return nil, err
			}
			service.Functions = serviceFunctions.Functions
		}

		_, err = api.Config.Get(appView.Namespace, generatedProgramConfigName, "")
		if err != nil {
			return nil, err
		}

		populateFunctionVolumeMount(service)
	}
	return appView, nil
}

func (api *API) ToApplication(appView *models.ApplicationView, oldApp *specV1.Application) (*specV1.Application, []specV1.Configuration, error) {
	app := new(specV1.Application)
	copier.Copy(app, appView)

	translateSecretLikedModelsToSecrets(appView, app)

	if app.Type != common.FunctionApp {
		return app, nil, nil
	}
	oldServices := map[string]bool{}
	if oldApp != nil {
		for _, service := range oldApp.Services {
			oldServices[service.Name] = true
		}
	}

	var configs []specV1.Configuration
	for index := range app.Services {
		service := &app.Services[index]
		config, err := generateConfigOfFunctionService(service, app)
		if err != nil {
			return nil, nil, err
		}
		configs = append(configs, *config)

		programConfig, err := generateProgramOfFunctionService(api, service, app)
		if err != nil {
			return nil, nil, err
		}
		configs = append(configs, *programConfig)

		if _, ok := oldServices[service.Name]; !ok {
			volumeMount, volume := generateVmAndMount(service.Name, config.Name)
			service.VolumeMounts = append(service.VolumeMounts, volumeMount)
			app.Volumes = append(app.Volumes, volume)

			volumeMountPrpgram, volumeProgram := genProgramVmAndMount(service.Name, programConfig.Name)
			service.VolumeMounts = append(service.VolumeMounts, volumeMountPrpgram)
			app.Volumes = append(app.Volumes, volumeProgram)
		}

		runtimes, err := api.Func.ListRuntimes()
		if err != nil {
			return nil, nil, err
		}
		image, ok := runtimes[service.FunctionConfig.Runtime]
		if !ok {
			return nil, nil, common.Error(common.ErrResourceNotFound,
				common.Field("type", "runtime"),
				common.Field("name", service.FunctionConfig.Runtime))
		}
		service.Image = image

		service.Ports = []specV1.ContainerPort{
			{
				ContainerPort: 80,
				Protocol:      "TCP",
			},
		}
	}
	return app, configs, nil
}

func translateSecretLikedModelsToSecrets(appView *models.ApplicationView, app *specV1.Application) {
	for k, v := range appView.Volumes {
		if v.Certificate == nil {
			continue
		}
		app.Volumes[k].Secret = &specV1.ObjectReference{
			Name: appView.Volumes[k].Certificate.Name,
		}
	}

	for _, reg := range appView.Registries {
		secretVolume := specV1.Volume{
			Name: reg.Name,
			VolumeSource: specV1.VolumeSource{
				Secret: &specV1.ObjectReference{
					Name: reg.Name,
				},
			},
		}
		app.Volumes = append(app.Volumes, secretVolume)
	}
}

func (api *API) translateSecretsToSecretLikedResources(appView *models.ApplicationView) error {
	appView.Registries = make([]models.RegistryView, 0)
	volumes := make([]models.VolumeView, 0)
	for _, volume := range appView.Volumes {
		if volume.Secret != nil {
			secret, err := api.Secret.Get(appView.Namespace, volume.Secret.Name, "")
			if err != nil {
				if e, ok := err.(errors.Coder); ok && e.Code() == common.ErrResourceNotFound {
					continue
				}
				return err
			}

			if label, ok := secret.Labels[specV1.SecretLabel]; ok && label == specV1.SecretRegistry {
				registry := models.FromSecretToRegistry(secret, false)
				appView.Registries = append(appView.Registries, models.RegistryView{
					Name:     registry.Name,
					Address:  registry.Address,
					Username: registry.Username,
				})
				continue
			}

			if label, ok := secret.Labels[specV1.SecretLabel]; ok && label == specV1.SecretCertificate {
				volume = models.VolumeView{
					Name: volume.Name,
					Certificate: &specV1.ObjectReference{
						Name:    volume.Secret.Name,
						Version: volume.Secret.Version,
					},
				}
			}
		}
		volumes = append(volumes, volume)
	}

	appView.Volumes = volumes

	return nil
}

func (api *API) validApplication(namesapce string, app *models.ApplicationView) error {
	for _, v := range app.Volumes {
		if v.Config != nil {
			_, err := api.Config.Get(namesapce, v.Config.Name, "")
			if err != nil {
				return err
			}
		}
		if v.Secret != nil {
			_, err := api.Secret.Get(namesapce, v.Secret.Name, "")
			if err != nil {
				return err
			}
		}
		if v.Certificate != nil {
			_, err := api.Secret.Get(namesapce, v.Certificate.Name, "")
			if err != nil {
				return err
			}
		}
	}

	for _, r := range app.Registries {
		_, err := api.Secret.Get(namesapce, r.Name, "")
		if err != nil {
			return err
		}
	}
	return nil
}

func (api *API) isAppCanDelete(namesapce, name string) (bool, error) {
	// TODO: improve
	if strings.HasPrefix(name, "baetyl-") {
		nodeNames, err := api.Index.ListNodesByApp(namesapce, name)
		if err != nil {
			return false, err
		}

		if len(nodeNames) > 0 {
			return false, nil
		}
	}
	return true, nil
}

func (api *API) updateGenConfigsOfFunctionApp(namespace string, configs []specV1.Configuration) error {
	for _, config := range configs {
		_, err := api.Config.Upsert(namespace, &config)
		if err != nil {
			return err
		}
	}
	return nil
}

func (api *API) cleanGenConfigsOfFunctionApp(configs []specV1.Configuration, oldApp *specV1.Application) {
	m := map[string]bool{}
	for _, config := range configs {
		m[config.Name] = true
	}

	for _, v := range oldApp.Volumes {
		if v.VolumeSource.Config != nil {
			if _, ok := m[v.VolumeSource.Config.Name]; !ok && (strings.HasPrefix(v.VolumeSource.Config.Name, FunctionConfigPrefix) ||
				strings.HasPrefix(v.VolumeSource.Config.Name, FunctionProgramConfigPrefix)) {
				err := api.Config.Delete(oldApp.Namespace, v.VolumeSource.Config.Name)
				if err != nil {
					common.LogDirtyData(err,
						log.Any("type", common.Config),
						log.Any(common.KeyContextNamespace, oldApp.Namespace),
						log.Any("name", v.VolumeSource.Config.Name))
					continue
				}
			}
		}
	}
}

func getGenConfigNameOfFunctionService(app *specV1.Application, serviceName string) (string, error) {
	volumeMountName := getNameOfFunctionConfigVolumeMount(serviceName)
	for _, v := range app.Volumes {
		if v.Name == volumeMountName {
			if v.VolumeSource.Config == nil {
				return "", common.Error(common.ErrVolumeType, common.Field("name", v.Name), common.Field("type", common.Config))
			}
			return v.VolumeSource.Config.Name, nil
		}
	}
	return strings.ToLower(fmt.Sprintf("%s-%s-%s-%s", FunctionConfigPrefix, app.Name, serviceName, common.RandString(9))), nil
}

func getGenProgramNameOfFunctionService(app *specV1.Application, serviceName string) (string, error) {
	volumeMountName := getNameOfFunctionProgramVmMount(serviceName)
	for _, v := range app.Volumes {
		if v.Name != volumeMountName {
			continue
		}
		if v.VolumeSource.Config == nil {
			return "", common.Error(common.ErrVolumeType, common.Field("name", v.Name), common.Field("type", common.Config))
		}
		return v.VolumeSource.Config.Name, nil
	}
	return strings.ToLower(fmt.Sprintf("%s-%s-%s-%s", FunctionProgramConfigPrefix, app.Name, serviceName, common.RandString(9))), nil
}

func generateConfigOfFunctionService(service *specV1.Service, app *specV1.Application) (*specV1.Configuration, error) {
	serviceFunctions := models.ServiceFunction{
		Functions: service.Functions,
	}

	data, err := json.Marshal(serviceFunctions)
	if err != nil {
		return nil, err
	}

	generatedConfigName, err := getGenConfigNameOfFunctionService(app, service.Name)
	if err != nil {
		return nil, err
	}

	config := &specV1.Configuration{
		Name:      generatedConfigName,
		Namespace: app.Namespace,
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Data: map[string]string{
			FunctionDefaultConfigFile: string(data),
		},
	}
	return config, nil
}

func generateProgramOfFunctionService(api *API, service *specV1.Service, app *specV1.Application) (*specV1.Configuration, error) {
	generatedConfigName, err := getGenProgramNameOfFunctionService(app, service.Name)
	if err != nil {
		return nil, err
	}

	config := &specV1.Configuration{}
	tempalteName := fmt.Sprintf("baetyl-%s-program.yml", service.FunctionConfig.Runtime)
	params := map[string]interface{}{
		"Namespace":  app.Namespace,
		"ConfigName": generatedConfigName,
	}

	err = api.Template.UnmarshalTemplate(tempalteName, params, config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return config, nil
}

func generateVmAndMount(serviceName, configName string) (specV1.VolumeMount, specV1.Volume) {
	volumeMount := specV1.VolumeMount{
		Name:      getNameOfFunctionConfigVolumeMount(serviceName),
		MountPath: ConfigDir,
		ReadOnly:  true,
	}

	generatedConfigVolume := specV1.Volume{
		Name: volumeMount.Name,
		VolumeSource: specV1.VolumeSource{
			Config: &specV1.ObjectReference{
				Name: configName,
			},
		},
	}
	return volumeMount, generatedConfigVolume
}

func genProgramVmAndMount(serviceName, configName string) (specV1.VolumeMount, specV1.Volume) {
	volumeMount := specV1.VolumeMount{
		Name:      getNameOfFunctionProgramVmMount(serviceName),
		MountPath: ProgramConfigDir,
		ReadOnly:  true,
	}

	generatedConfigVolume := specV1.Volume{
		Name: volumeMount.Name,
		VolumeSource: specV1.VolumeSource{
			Config: &specV1.ObjectReference{
				Name: configName,
			},
		},
	}
	return volumeMount, generatedConfigVolume
}

func populateFunctionVolumeMount(service *specV1.Service) {
	codeVm := getNameOfFunctionCodeVolumeMount(service.Name)
	confVm := getNameOfFunctionConfigVolumeMount(service.Name)
	programConfVm := getNameOfFunctionProgramVmMount(service.Name)

	for i := range service.VolumeMounts {
		mount := &service.VolumeMounts[i]
		if mount.Name == codeVm || mount.Name == confVm || mount.Name == programConfVm {
			mount.Immutable = true
		}
	}
}

func getNameOfFunctionConfigVolumeMount(serviceName string) string {
	return fmt.Sprintf("%s-%s", FunctionConfigPrefix, serviceName)
}

func getNameOfFunctionProgramVmMount(serviceName string) string {
	return fmt.Sprintf("%s-%s", FunctionProgramConfigPrefix, serviceName)
}

func getNameOfFunctionCodeVolumeMount(serviceName string) string {
	return fmt.Sprintf("%s-%s", FunctionCodePrefix, serviceName)
}
