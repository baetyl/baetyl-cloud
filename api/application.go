package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/jinzhu/copier"
	v1 "k8s.io/api/core/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

const (
	ConfigDir                   = "/etc/baetyl"
	ProgramConfigDir            = "/var/lib/baetyl/bin"
	FunctionConfigPrefix        = "baetyl-function-config"
	ProgramConfigPrefix         = "baetyl-program-config"
	FunctionProgramConfigPrefix = "baetyl-function-program-config"
	FunctionCodePrefix          = "baetyl-function-code"
	FunctionDefaultConfigFile   = "conf.yml"

	HookCreateApplicationOta = "hookCreateApplicationOta"
	HookUpdateApplicationOta = "hookUpdateApplicationOta"
	HookDeleteApplicationOta = "hookDeleteApplicationOta"
)

type CreateApplicationOta = func(c *common.Context, app *specV1.Application) (*specV1.Application, error)
type UpdateApplicationOta = func(c *common.Context, app *specV1.Application) (*specV1.Application, error)
type DeleteApplicationOta = func(c *common.Context, app *specV1.Application) error

// GetApplication get a application
func (api *API) GetApplication(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	app, err := api.Facade.GetApp(ns, n, "")
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
	params, err := api.ParseListOptionsAppendSystemLabel(c)
	if err != nil {
		return nil, err
	}
	apps, err := api.App.List(ns, params)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	api.ToApplicationListView(apps)
	return apps, err
}

// CreateApplication create one application
func (api *API) CreateApplication(c *common.Context) (interface{}, error) {
	appView, err := api.ParseApplication(c)
	if err != nil {
		return nil, err
	}
	ns, name := c.GetNamespace(), appView.Name
	appView.Namespace = ns

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

	if f, exist := api.Hooks[HookCreateApplicationOta]; exist {
		if hk, ok := f.(CreateApplicationOta); ok {
			app, err = hk(c, app)
			if err != nil {
				return nil, err
			}
		}
	}

	log.L().Info("", log.Any("app2", app))
	app, err = api.Facade.CreateApp(ns, baseApp, app, configs)
	if err != nil {
		return nil, err
	}

	return api.ToApplicationView(app)
}

// UpdateApplication update the application
func (api *API) UpdateApplication(c *common.Context) (interface{}, error) {
	appView, err := api.ParseApplication(c)
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
	if CheckIsSysResources(oldApp.Labels) &&
		(oldApp.Selector != appView.Selector || !reflect.DeepEqual(oldApp.Labels, appView.Labels) || !appView.System) {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "selector，labels or system field can't be modified of sys apps"))
	}

	if appView.CronStatus != oldApp.CronStatus {
		if oldApp.CronStatus != specV1.CronWait || appView.CronStatus != specV1.CronNotSet {
			return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "failed to add cron, can't set up a cron job which has been deployed"))
		}
	} else {
		appView.CronTime = oldApp.CronTime
	}

	appView.Version = oldApp.Version
	appView.CreationTimestamp = oldApp.CreationTimestamp
	app, configs, err := api.ToApplication(appView, oldApp)
	if err != nil {
		return nil, err
	}

	// ota can not modify
	app.Ota = oldApp.Ota

	if f, exist := api.Hooks[HookUpdateApplicationOta]; exist {
		if hk, ok := f.(UpdateApplicationOta); ok {
			app, err = hk(c, app)
			if err != nil {
				return nil, err
			}
		}
	}

	app, err = api.Facade.UpdateApp(ns, oldApp, app, configs)
	if err != nil {
		return nil, errors.Trace(err)
	}

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

	if canDelete, err := api.IsAppCanDelete(ns, name); err != nil {
		return nil, err
	} else if !canDelete {
		return nil, common.Error(common.ErrAppReferencedByNode, common.Field("name", name))
	}

	if f, exist := api.Hooks[HookDeleteApplicationOta]; exist {
		if hk, ok := f.(DeleteApplicationOta); ok {
			err = hk(c, app)
			if err != nil {
				return nil, err
			}
		}
	}

	err = api.Facade.DeleteApp(ns, name, app)
	return nil, err
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

func (api *API) ParseApplication(c *common.Context) (*models.ApplicationView, error) {
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
			if app.Mode == context.RunModeKube && v.Image == "" {
				return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "image is required in kube mode"))
			}
			if app.Mode == context.RunModeNative && v.ProgramConfig == "" {
				return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "program config is required in native mode"))
			}
		}
	} else if app.Type == common.FunctionApp {
		for _, v := range app.Services {
			if v.FunctionConfig == nil {
				return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "function config can't be empty in function app"))
			}
		}
		if len(app.Registries) != 0 {
			return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "registries should be empty in function app"))
		}
	} else {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "type is invalid"))
	}

	// multi-container compatibility
	api.compatibleAppDeprecatedFiled(app)

	if app.Workload != specV1.WorkloadDeployment &&
		app.Workload != specV1.WorkloadDaemonSet &&
		app.Workload != specV1.WorkloadStatefulSet &&
		app.Workload != specV1.WorkloadJob {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error",
			"failed to parse service type, service type should be deployment / daemonset / statefulset / job"))
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

func (api *API) ParseListOptions(c *common.Context) (*models.ListOptions, error) {
	params := &models.ListOptions{}
	if err := c.Bind(params); err != nil {
		return nil, err
	}
	return params, nil
}

func (api *API) ParseListOptionsAppendSystemLabel(c *common.Context) (*models.ListOptions, error) {
	opt, err := api.ParseListOptions(c)
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
	nodes, err := api.Node.UpdateNodeAppVersion(nil, namespace, app)
	if err != nil {
		return err
	}
	return api.Index.RefreshNodesIndexByApp(nil, namespace, app.Name, nodes)
}

func (api *API) ToApplicationView(app *specV1.Application) (*models.ApplicationView, error) {
	appView := &models.ApplicationView{}
	copier.Copy(appView, app)

	err := api.translateSecretsToSecretLikedResources(appView)
	if err != nil {
		return nil, err
	}

	err = api.translateToNativeAppView(appView)
	if err != nil {
		return nil, err
	}

	api.compatibleAppDeprecatedFiled(appView)
	populateAppDefaultField(appView)

	if app.Type != common.FunctionApp {
		delete(appView.Labels, common.LabelAppMode)
		return appView, nil
	}
	for index := range appView.Services {
		service := &appView.Services[index]
		generatedConfigName, err := getGenConfigNameOfFunctionService(app, service.Name)
		if err != nil {
			return nil, err
		}

		var generatedProgramConfigName string
		if app.Mode == context.RunModeNative {
			generatedProgramConfigName, err = getGenProgramNameOfFunctionService(app, service.Name)
			if err != nil {
				return nil, err
			}
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

		if app.Mode == context.RunModeNative {
			_, err = api.Config.Get(appView.Namespace, generatedProgramConfigName, "")
			if err != nil {
				return nil, err
			}
		}

		populateFunctionVolumeMount(service)
	}

	// kube function set program nil
	if appView.Type == common.FunctionApp && appView.Mode == context.RunModeKube {
		for i := 0; i < len(appView.Volumes); i++ {
			if strings.HasPrefix(appView.Volumes[i].Name, "baetyl-function-program-config") {
				appView.Volumes = append(appView.Volumes[:i], appView.Volumes[i+1:]...)
				i--
			}
		}
		for i, svc := range appView.Services {
			for j := 0; j < len(svc.VolumeMounts); j++ {
				if strings.HasPrefix(svc.VolumeMounts[j].Name, "baetyl-function-program-config") {
					svc.VolumeMounts = append(svc.VolumeMounts[:j], svc.VolumeMounts[j+1:]...)
					j--
				}
			}
			appView.Services[i] = svc
		}
	}
	delete(appView.Labels, common.LabelAppMode)
	return appView, nil
}

func (api *API) ToApplication(appView *models.ApplicationView, oldApp *specV1.Application) (*specV1.Application, []specV1.Configuration, error) {
	app := new(specV1.Application)
	copier.Copy(app, appView)

	translateSecretLikedModelsToSecrets(appView, app)
	translateNativeApp(appView, app, oldApp)

	if app.Type != common.FunctionApp {
		return app, nil, nil
	}
	oldServices := map[string]bool{}
	if oldApp != nil {
		for _, service := range oldApp.Services {
			oldServices[service.Name] = true
		}
	}
	volMap := map[string]bool{}
	for _, vol := range app.Volumes {
		volMap[vol.Name] = true
	}

	var configs []specV1.Configuration
	for index := range app.Services {
		service := &app.Services[index]
		config, err := generateConfigOfFunctionService(service, app)
		if err != nil {
			return nil, nil, err
		}
		configs = append(configs, *config)

		var programConfig *specV1.Configuration
		if app.Mode == context.RunModeNative {
			programConfig, err = generateProgramOfFunctionService(api, service, app)
			if err != nil {
				return nil, nil, err
			}
			configs = append(configs, *programConfig)
		}

		if _, ok := oldServices[service.Name]; !ok {
			vmName := getNameOfFunctionConfigVolumeMount(service.Name)
			if !volMap[vmName] {
				volumeMount, volume := generateVmAndMount(config.Name, vmName, ConfigDir)
				service.VolumeMounts = append(service.VolumeMounts, volumeMount)
				app.Volumes = append(app.Volumes, volume)
			}

			if app.Mode == context.RunModeNative {
				vmName = getNameOfFunctionProgramVmMount(service.Name)
				if !volMap[vmName] {
					volumeMountPrpgram, volumeProgram := generateVmAndMount(programConfig.Name, vmName, ProgramConfigDir)
					service.VolumeMounts = append(service.VolumeMounts, volumeMountPrpgram)
					app.Volumes = append(app.Volumes, volumeProgram)
				}
			}
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

func (api *API) ToApplicationListView(apps *models.ApplicationList) {
	for i := range apps.Items {
		if apps.Items[i].Mode == "" {
			apps.Items[i].Mode = context.RunModeKube
		}
		delete(apps.Items[i].Labels, common.LabelAppMode)
	}
}

func translateNativeApp(appView *models.ApplicationView,
	app *specV1.Application, oldApp *specV1.Application) {
	if appView.Mode != context.RunModeNative || appView.Type == common.FunctionApp {
		return
	}

	for index := range appView.Services {
		service := &app.Services[index]
		serviceView := &appView.Services[index]
		vmName := getNameOfNativeProgramVolumeMount(service.Name)
		volumeMount, volume := generateVmAndMount(serviceView.ProgramConfig, vmName, ProgramConfigDir)

		var exist bool
		for _, v := range service.VolumeMounts {
			if v.Name == vmName {
				exist = true
				break
			}
		}
		if !exist {
			service.VolumeMounts = append(service.VolumeMounts, volumeMount)
		}

		exist = false
		for i, v := range app.Volumes {
			if v.Name == vmName {
				app.Volumes[i] = volume
				exist = true
				break
			}
		}
		if !exist {
			app.Volumes = append(app.Volumes, volume)
		}
	}
}

func (api *API) translateToNativeAppView(appView *models.ApplicationView) error {
	if appView.Mode != context.RunModeNative || appView.Type == common.FunctionApp {
		return nil
	}
	for index := range appView.Services {
		service := &appView.Services[index]
		vmName := getNameOfNativeProgramVolumeMount(service.Name)
		configName, err := getNameOfNativeProgramConfig(appView, vmName)
		if err != nil {
			return err
		}
		service.ProgramConfig = configName

		_, err = api.Config.Get(appView.Namespace, configName, "")
		if err != nil {
			return err
		}

		for i := range service.VolumeMounts {
			mount := &service.VolumeMounts[i]
			if mount.Name == vmName {
				mount.Immutable = true
				break
			}
		}
	}
	return nil
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

func (api *API) validApplication(namespace string, app *models.ApplicationView) error {
	for _, v := range app.Volumes {
		if v.Config != nil {
			// native program config will be validated by service.ProgramConfig
			if isProgramConfig(v.Name) {
				continue
			}
			_, err := api.Config.Get(namespace, v.Config.Name, "")
			if err != nil {
				return err
			}
		}
		if v.Secret != nil {
			_, err := api.Secret.Get(namespace, v.Secret.Name, "")
			if err != nil {
				return err
			}
		}
		if v.Certificate != nil {
			_, err := api.Secret.Get(namespace, v.Certificate.Name, "")
			if err != nil {
				return err
			}
		}
	}

	for _, r := range app.Registries {
		_, err := api.Secret.Get(namespace, r.Name, "")
		if err != nil {
			return err
		}
	}

	ports := make(map[int32]bool)
	for _, service := range app.Services {
		if service.ProgramConfig != "" {
			_, err := api.Config.Get(namespace, service.ProgramConfig, "")
			if err != nil {
				return err
			}
		}
		hostPortNum, err := isValidPort(&service, ports)
		if err != nil {
			return err
		}
		if app.Replica > 1 && hostPortNum > 0 {
			return common.Error(common.ErrRequestParamInvalid,
				common.Field("error", "port mapping is only supported under single replica"))
		}
	}
	for _, service := range app.InitServices {
		hostPortNum, err := isValidPort(&service, ports)
		if err != nil {
			return err
		}
		if app.Replica > 1 && hostPortNum > 0 {
			return common.Error(common.ErrRequestParamInvalid,
				common.Field("error", "port mapping is only supported under single replica"))
		}
	}

	if app.CronStatus == specV1.CronWait && app.CronTime.Before(time.Now()) {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", "failed to add cron job, time should be set after now"))
	}

	if app.AutoScaleCfg != nil {
		if app.AutoScaleCfg.MinReplicas <= 0 {
			return common.Error(common.ErrRequestParamInvalid,
				common.Field("error", "minReplicas must be greater than 0"))
		}
		if app.AutoScaleCfg.MaxReplicas < app.AutoScaleCfg.MinReplicas {
			return common.Error(common.ErrRequestParamInvalid,
				common.Field("error", "minReplicas must be less than maxReplicas"))
		}
		for _, metrics := range app.AutoScaleCfg.Metrics {
			if metrics.Resource.Name != "cpu" && metrics.Resource.Name != "memory" {
				return common.Error(common.ErrRequestParamInvalid,
					common.Field("error", "only cpu and memory is supported"))
			}
		}
	}

	app.Labels = common.AddSystemLabel(app.Labels, map[string]string{
		common.LabelAppMode: app.Mode,
	})
	return nil
}

func isValidPort(service *models.ServiceView, ports map[int32]bool) (int, error) {
	hostPortNum := 0
	for _, port := range service.Ports {
		if port.ServiceType == string(v1.ServiceTypeNodePort) {
			if port.NodePort <= 0 {
				return hostPortNum, common.Error(common.ErrRequestParamInvalid, common.Field("error", "invalid NodePort"))
			}
			if _, ok := ports[port.NodePort]; ok {
				return hostPortNum, common.Error(common.ErrRequestParamInvalid, common.Field("error", "duplicate node ports"))
			} else {
				ports[port.NodePort] = true
			}
		} else {
			if port.HostPort == 0 {
				continue
			}
			if _, ok := ports[port.HostPort]; ok {
				return hostPortNum, common.Error(common.ErrRequestParamInvalid, common.Field("error", "duplicate host ports"))
			} else {
				hostPortNum++
				ports[port.HostPort] = true
			}
		}
	}
	return hostPortNum, nil
}

func isProgramConfig(volume string) bool {
	if strings.HasPrefix(volume, ProgramConfigPrefix) {
		return true
	}
	return false
}

func (api *API) IsAppCanDelete(namesapce, name string) (bool, error) {
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

func (api *API) compatibleAppDeprecatedFiled(app *models.ApplicationView) {
	// Workload
	if app.Workload == "" {
		// compatible with the original one service corresponding to one workload
		if len(app.Services) > 0 && app.Services[0].Type != "" {
			api.log.Debug("workload is empty, use the services[0].Type ")
			app.Workload = app.Services[0].Type
			for i := 1; i < len(app.Services); i++ {
				if app.Services[i].Type != app.Workload {
					api.log.Warn("app service type is inconsistent", log.Any("index", i), log.Any("name", app.Services[i].Name))
				}
			}
		} else {
			app.Workload = specV1.WorkloadDeployment
		}
	} else {
		for i, svc := range app.Services {
			if svc.Type == "" {
				app.Services[i].Type = app.Workload
			}
		}
	}

	// HostNetwork
	if !app.HostNetwork && len(app.Services) > 0 && app.Services[0].HostNetwork {
		api.log.Debug("app.HostNetwork is false, use the services[0].HostNetwork true ")
		app.HostNetwork = true
	} else {
		if app.HostNetwork {
			for i, _ := range app.Services {
				app.Services[i].HostNetwork = true
			}
		}
	}

	// Replica
	if app.Replica == 0 {
		// compatible with the original one service corresponding to one workload
		if len(app.Services) > 0 && app.Services[0].Replica != 0 {
			api.log.Debug("app.Replica is 0, use the services[0].Replica ")
			app.Replica = app.Services[0].Replica
			for i := 1; i < len(app.Services); i++ {
				if app.Services[i].Replica != app.Replica {
					api.log.Warn("app service replica is inconsistent", log.Any("index", i), log.Any("name", app.Services[i].Name))
				}
			}
		} else {
			app.Replica = 1
		}
	} else {
		for i, svc := range app.Services {
			if svc.Replica == 0 {
				app.Services[i].Replica = app.Replica
			}
		}
	}

	// JobConfig
	if app.JobConfig == nil || app.JobConfig.RestartPolicy == "" {
		// compatible with the original one service corresponding to one workload
		if len(app.Services) > 0 && app.Services[0].JobConfig != nil {
			api.log.Debug("app.JobConfig is 0, use the services[0].JobConfig ")
			app.JobConfig = &specV1.AppJobConfig{
				Completions:   app.Services[0].JobConfig.Completions,
				Parallelism:   app.Services[0].JobConfig.Parallelism,
				BackoffLimit:  app.Services[0].JobConfig.BackoffLimit,
				RestartPolicy: app.Services[0].JobConfig.RestartPolicy,
			}
		} else {
			app.JobConfig = &specV1.AppJobConfig{RestartPolicy: "Never", Completions: 1}
			for i, svc := range app.Services {
				if svc.JobConfig == nil || svc.JobConfig.RestartPolicy == "" {
					app.Services[i].JobConfig = &specV1.ServiceJobConfig{
						Completions:   app.JobConfig.Completions,
						Parallelism:   app.JobConfig.Parallelism,
						BackoffLimit:  app.JobConfig.BackoffLimit,
						RestartPolicy: app.JobConfig.RestartPolicy,
					}
				}
			}
		}
	} else {
		for i, svc := range app.Services {
			if svc.JobConfig == nil || svc.JobConfig.RestartPolicy == "" {
				app.Services[i].JobConfig = &specV1.ServiceJobConfig{
					Completions:   app.JobConfig.Completions,
					Parallelism:   app.JobConfig.Parallelism,
					BackoffLimit:  app.JobConfig.BackoffLimit,
					RestartPolicy: app.JobConfig.RestartPolicy,
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

func generateVmAndMount(configName string, vmName, mountPath string) (specV1.VolumeMount, specV1.Volume) {
	volumeMount := specV1.VolumeMount{
		Name:      vmName,
		MountPath: mountPath,
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

func populateFunctionVolumeMount(service *models.ServiceView) {
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

func populateAppDefaultField(appView *models.ApplicationView) {
	if appView.Mode == "" {
		appView.Mode = context.RunModeKube
	}
	if appView.Workload == "" {
		appView.Workload = specV1.WorkloadDeployment
	}
}

func getNameOfNativeProgramConfig(appView *models.ApplicationView, vmName string) (string, error) {
	for _, v := range appView.Volumes {
		if v.Name == vmName {
			if v.Config == nil {
				return "", common.Error(common.ErrVolumeType, common.Field("name", v.Name), common.Field("type", common.Config))
			}
			return v.Config.Name, nil
		}
	}
	return "", common.Error(common.ErrResourceNotFound,
		common.Field("type", "volume"),
		common.Field("name", vmName))
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

func getNameOfNativeProgramVolumeMount(serviceName string) string {
	return fmt.Sprintf("%s-%s", ProgramConfigPrefix, serviceName)
}
