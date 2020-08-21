package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/jinzhu/copier"
)

const (
	ConfigDir                 = "/etc/baetyl"
	FunctionConfigPrefix      = "baetyl-function-config"
	FunctionCodePrefix        = "baetyl-function-code"
	FunctionDefaultConfigFile = "service.yml"
)

// GetApplication get a application
func (api *API) GetApplication(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	app, err := api.applicationService.Get(ns, n, "")
	if err != nil {
		return nil, err
	}

	return api.toApplicationView(app)
}

// ListApplication list application
func (api *API) ListApplication(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	apps, err := api.applicationService.List(ns, api.parseListOptionsAppendSystemLabel(c))
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
	oldApp, err := api.applicationService.Get(ns, name, "")
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

	app, configs, err := api.toAppliation(appView, nil)
	if err != nil {
		return nil, err
	}

	err = api.updateGeneratedConfigsOfFunctionApp(ns, configs)
	if err != nil {
		return nil, err
	}

	app, err = api.applicationService.CreateWithBase(ns, app, baseApp)
	if err != nil {
		return nil, err
	}

	err = api.updateNodeAndAppIndex(ns, app)
	if err != nil {
		return nil, err
	}

	return api.toApplicationView(app)
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

	oldApp, err := api.applicationService.Get(ns, name, "")
	if err != nil {
		return nil, err
	}

	appView.Version = oldApp.Version
	app, configs, err := api.toAppliation(appView, oldApp)
	if err != nil {
		return nil, err
	}

	err = api.updateGeneratedConfigsOfFunctionApp(ns, configs)
	if err != nil {
		return nil, err
	}

	app, err = api.applicationService.Update(ns, app)
	if err != nil {
		return nil, err
	}

	if oldApp != nil && oldApp.Selector != app.Selector {
		// delete old nodes
		if err := api.deleteNodeAndAppIndex(ns, oldApp); err != nil {
			return nil, err
		}
	}

	// update nodes
	if err := api.updateNodeAndAppIndex(ns, app); err != nil {
		return nil, err
	}

	api.cleanGeneratedConfigsOfFunctionApp(configs, oldApp)

	return api.toApplicationView(app)
}

// DeleteApplication delete the application
func (api *API) DeleteApplication(c *common.Context) (interface{}, error) {
	ns, name := c.GetNamespace(), c.GetNameFromParam()
	app, err := api.applicationService.Get(ns, name, "")
	if err != nil {
		return nil, err
	}

	if canDelete, err := api.isAppCanDelete(ns, name); err != nil {
		return nil, err
	} else if !canDelete {
		return nil, common.Error(common.ErrAppReferencedByNode, common.Field("name", name))
	}

	if err := api.applicationService.Delete(ns, c.GetNameFromParam(), ""); err != nil {
		return nil, err
	}

	//delete the app from node
	if err := api.deleteNodeAndAppIndex(ns, app); err != nil {
		return nil, err
	}

	api.cleanGeneratedConfigsOfFunctionApp(nil, app)
	return nil, nil
}

func (api *API) parseApplication(c *common.Context) (*models.ApplicationView, error) {
	app := new(models.ApplicationView)
	app.Name = c.GetNameFromParam()
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
		baseApp, err := api.applicationService.Get(namespace, base, "")
		if err != nil {
			return nil, err
		}
		return baseApp, nil
	}
	return nil, nil
}

func (api *API) parseListOptions(c *common.Context) *models.ListOptions {
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 64)
	lp := &models.ListOptions{
		LabelSelector: c.Query("selector"),
		FieldSelector: c.Query("fieldSelector"),
		Limit:         limit,
		Continue:      c.Query("continue"),
	}
	return lp
}

func (api *API) parseListOptionsAppendSystemLabel(c *common.Context) *models.ListOptions {
	opt := api.parseListOptions(c)

	ls := opt.LabelSelector
	if !strings.Contains(ls, common.LabelSystem) {
		if len(strings.TrimSpace(ls)) > 0 {
			ls += ","
		}
		ls += "!" + common.LabelSystem
	}

	opt.LabelSelector = ls
	return opt
}

func (api *API) updateNodeAndAppIndex(namespace string, app *specV1.Application) error {
	nodes, err := api.nodeService.UpdateNodeAppVersion(namespace, app)
	if err != nil {
		return err
	}
	return api.indexService.RefreshNodesIndexByApp(namespace, app.Name, nodes)
}

func (api *API) deleteNodeAndAppIndex(namespace string, app *specV1.Application) error {
	_, err := api.nodeService.DeleteNodeAppVersion(namespace, app)
	if err != nil {
		return err
	}

	return api.indexService.RefreshNodesIndexByApp(namespace, app.Name, make([]string, 0))
}

func (api *API) toApplicationView(app *specV1.Application) (*models.ApplicationView, error) {
	appView := &models.ApplicationView{}
	copier.Copy(appView, app)

	err := api.translateSecretsToRegistries(appView)
	if err != nil {
		return nil, err
	}

	if app.Type != common.FunctionApp {
		return appView, nil
	}
	for index := range appView.Services {
		service := &appView.Services[index]
		generatedConfigName, err := getGeneratedConfigNameOfFunctionService(app, service.Name)
		if err != nil {
			return nil, err
		}

		config, err := api.configService.Get(appView.Namespace, generatedConfigName, "")
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

		populateFunctionVolumeMount(service)
	}
	return appView, nil
}

func (api *API) toAppliation(appView *models.ApplicationView, oldApp *specV1.Application) (*specV1.Application, []specV1.Configuration, error) {
	app := new(specV1.Application)
	copier.Copy(app, appView)

	translateReistriesToSecrets(appView, app)

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

		if _, ok := oldServices[service.Name]; !ok {
			volumeMount, volume := generateVolumeAndVolumeMount(service.Name, config.Name)
			service.VolumeMounts = append(service.VolumeMounts, volumeMount)
			app.Volumes = append(app.Volumes, volume)
		}

		image, err := api.getFunctionImageByRuntime(service.FunctionConfig.Runtime)
		if err != nil {
			return nil, nil, err
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

func translateReistriesToSecrets(appView *models.ApplicationView, app *specV1.Application) {
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

func (api *API) translateSecretsToRegistries(appView *models.ApplicationView) error {
	appView.Registries = make([]models.RegistryView, 0)
	volumes := make([]specV1.Volume, 0)
	for _, volume := range appView.Volumes {
		if volume.Secret != nil {
			secret, err := api.secretService.Get(appView.Namespace, volume.Secret.Name, "")
			if err != nil {
				if e, ok := err.(errors.Coder); ok && e.Code() == common.ErrResourceNotFound {
					continue
				}
				return err
			}

			if label, ok := secret.Labels[specV1.SecretLabel]; ok && label == specV1.SecretRegistry {
				registry := models.FromSecret(secret)
				appView.Registries = append(appView.Registries, models.RegistryView{
					Name:     registry.Name,
					Address:  registry.Address,
					Username: registry.Username,
				})
				continue
			}
		}
		volumes = append(volumes, volume)
	}

	appView.Volumes = volumes

	return nil
}

func (api *API) validApplication(namesapce string, app *models.ApplicationView) error {
	for _, v := range app.Volumes {
		if v.VolumeSource.Config != nil {
			_, err := api.configService.Get(namesapce, v.VolumeSource.Config.Name, "")
			if err != nil {
				return err
			}
		}
		if v.VolumeSource.Secret != nil {
			_, err := api.secretService.Get(namesapce, v.VolumeSource.Secret.Name, "")
			if err != nil {
				return err
			}
		}
	}

	for _, r := range app.Registries {
		_, err := api.secretService.Get(namesapce, r.Name, "")
		if err != nil {
			return err
		}
	}
	return nil
}

func (api *API) isAppCanDelete(namesapce, name string) (bool, error) {
	for _, sysAppPrefix := range SystemApps {
		if strings.Contains(name, string(sysAppPrefix)) {
			nodeNames, err := api.indexService.ListNodesByApp(namesapce, name)
			if err != nil {
				return false, err
			}

			if len(nodeNames) > 0 {
				return false, nil
			}
		}
	}
	return true, nil
}

func (api *API) updateGeneratedConfigsOfFunctionApp(namespace string, configs []specV1.Configuration) error {
	for _, config := range configs {
		_, err := api.configService.Upsert(namespace, &config)
		if err != nil {
			return err
		}
	}
	return nil
}

func (api *API) cleanGeneratedConfigsOfFunctionApp(configs []specV1.Configuration, oldApp *specV1.Application) {
	m := map[string]bool{}
	for _, config := range configs {
		m[config.Name] = true
	}

	for _, v := range oldApp.Volumes {
		if v.VolumeSource.Config != nil {
			if _, ok := m[v.VolumeSource.Config.Name]; !ok &&
				strings.HasPrefix(v.VolumeSource.Config.Name, FunctionConfigPrefix) {
				err := api.configService.Delete(oldApp.Namespace, v.VolumeSource.Config.Name)
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

func getGeneratedConfigNameOfFunctionService(app *specV1.Application, serviceName string) (string, error) {
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

func generateConfigOfFunctionService(service *specV1.Service, app *specV1.Application) (*specV1.Configuration, error) {
	serviceFunctions := models.ServiceFunction{
		Functions: service.Functions,
	}

	data, err := json.Marshal(serviceFunctions)
	if err != nil {
		return nil, err
	}

	generatedConfigName, err := getGeneratedConfigNameOfFunctionService(app, service.Name)
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

func generateVolumeAndVolumeMount(serviceName, configName string) (specV1.VolumeMount, specV1.Volume) {
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

func populateFunctionVolumeMount(service *specV1.Service) {
	codeVm := getNameOfFunctionCodeVolumeMount(service.Name)
	confVm := getNameOfFunctionConfigVolumeMount(service.Name)

	for i := range service.VolumeMounts {
		mount := &service.VolumeMounts[i]
		if mount.Name == codeVm || mount.Name == confVm {
			mount.Immutable = true
		}
	}
}

func getNameOfFunctionConfigVolumeMount(serviceName string) string {
	return fmt.Sprintf("%s-%s", FunctionConfigPrefix, serviceName)
}

func getNameOfFunctionCodeVolumeMount(serviceName string) string {
	return fmt.Sprintf("%s-%s", FunctionCodePrefix, serviceName)
}
