package api

import (
	"bytes"
	"io"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v2"
	appv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubectl/pkg/scheme"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

const (
	TypeSecret    = "Secret"
	TypeConfig    = "ConfigMap"
	TypeDeploy    = "Deployment"
	TypeDaemonset = "DaemonSet"
	TypeJob       = "Job"
	TypeService   = "Service"
)

// yaml resources api
func (api *API) CreateYamlResource(c *common.Context) (interface{}, error) {
	var res models.YamlResourceList
	ns := c.GetNamespace()
	resources, err := api.parseYamlFileAndCheck(c)
	if err != nil {
		return nil, err
	}

	for _, r := range resources {
		switch r.GetObjectKind().GroupVersionKind().Kind {
		case TypeSecret:
			se, err := api.generateSecret(ns, r)
			if err != nil {
				return nil, err
			}
			res.Items = append(res.Items, se)
			res.Total++
		case TypeConfig:
			cfg, err := api.generateConfig(ns, c.GetUser().ID, r)
			if err != nil {
				return nil, err
			}
			res.Items = append(res.Items, cfg)
			res.Total++
		case TypeDeploy, TypeDaemonset, TypeJob:
			app, err := api.generateApplication(ns, r)
			if err != nil {
				return nil, err
			}
			res.Items = append(res.Items, app)
			res.Total++
		case TypeService:
			err = api.generateService(ns, r)
			if err != nil {
				return nil, err
			}
		}
	}
	return res, nil
}

func (api *API) UpdateYamlResource(c *common.Context) (interface{}, error) {
	var res models.YamlResourceList
	ns := c.GetNamespace()
	resources, err := api.parseYamlFileAndCheck(c)
	if err != nil {
		return nil, err
	}

	for _, r := range resources {
		switch r.GetObjectKind().GroupVersionKind().Kind {
		case TypeSecret:
			se, err := api.updateSecret(ns, r)
			if err != nil {
				return nil, err
			}
			res.Items = append(res.Items, se)
			res.Total++
		case TypeConfig:
			cfg, err := api.updateConfig(ns, c.GetUser().ID, r)
			if err != nil {
				return nil, err
			}
			res.Items = append(res.Items, cfg)
			res.Total++
		case TypeDeploy, TypeDaemonset, TypeJob:
			app, err := api.updateApplication(ns, r)
			if err != nil {
				return nil, err
			}
			res.Items = append(res.Items, app)
			res.Total++
		case TypeService:
			err = api.updateService(ns, r)
			if err != nil {
				return nil, err
			}
		}
	}
	return res, nil
}

func (api *API) DeleteYamlResource(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	resources, err := api.parseYamlFileAndCheck(c)
	if err != nil {
		return nil, err
	}
	// 逆序删除，避免依赖
	for i := len(resources) - 1; i >= 0; i-- {
		switch resources[i].GetObjectKind().GroupVersionKind().Kind {
		case TypeSecret:
			_, err := api.deleteSecret(ns, resources[i])
			if err != nil {
				return nil, err
			}
		case TypeConfig:
			_, err := api.deleteConfig(ns, resources[i])
			if err != nil {
				return nil, err
			}
		case TypeDeploy, TypeDaemonset, TypeJob:
			_, err := api.deleteApplication(ns, resources[i])
			if err != nil {
				return nil, err
			}
		case TypeService:
			_, err := api.deleteService(ns, resources[i])
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, err
}

func (api *API) parseYamlFileAndCheck(c *common.Context) ([]runtime.Object, error) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		err = common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
		return nil, err
	}
	err = fileCheck(header.Filename)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	if err != nil {
		return nil, err
	}

	res := api.parseK8SYaml(buf.Bytes())
	return res, nil
}

func (api *API) parseK8SYaml(fileR []byte) []runtime.Object {
	acceptedK8sTypes := regexp.MustCompile(`(Secret|ConfigMap|Deployment|DaemonSet|Job|Service)`)
	fileAsString := string(fileR[:])
	sepYamlfiles := strings.Split(fileAsString, "---")

	res := make([]runtime.Object, 0, len(sepYamlfiles))
	deploys := make([]runtime.Object, 0)
	services := make([]runtime.Object, 0)

	for _, f := range sepYamlfiles {
		if f == "\n" || f == "" {
			// ignore empty cases
			continue
		}

		decode := scheme.Codecs.UniversalDeserializer().Decode
		obj, groupVersionKind, err := decode([]byte(f), nil, nil)
		if err != nil {
			api.log.Warn("Error while decoding YAML object", log.Error(err))
			continue
		}

		if kind := groupVersionKind.Kind; !acceptedK8sTypes.MatchString(kind) {
			api.log.Error("K8s object types not supported!", log.Any("Skipping object with type: %s", kind))
		} else if kind == TypeSecret || kind == TypeConfig {
			res = append(res, obj)
		} else if kind == TypeDeploy || kind == TypeDaemonset || kind == TypeJob {
			deploys = append(deploys, obj)
		} else {
			services = append(services, obj)
		}
	}

	res = append(res, deploys...)
	res = append(res, services...)
	return res
}

// secret resource
func (api *API) generateSecret(ns string, r runtime.Object) (interface{}, error) {
	var err error
	sec, ok := r.(*corev1.Secret)
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "k8s secret typecasting failed"))
	}

	secret, err := api.generateSecretResource(ns, sec)
	if err != nil {
		return nil, err
	}

	sd, err := api.Secret.Get(ns, secret.Name, "")
	if err != nil {
		if e, ok := err.(errors.Coder); !ok || e.Code() != common.ErrResourceNotFound {
			return nil, err
		}
	}
	if sd != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "secret name is already in use"))
	}

	res, err := api.Facade.CreateSecret(ns, secret)
	if err != nil {
		return nil, err
	}
	return api.ToFilteredSecretView(res), nil
}

func (api *API) updateSecret(ns string, r runtime.Object) (interface{}, error) {
	var err error
	sec, ok := r.(*corev1.Secret)
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "k8s secret typecasting failed"))
	}

	secret, err := api.generateSecretResource(ns, sec)
	if err != nil {
		return nil, err
	}

	oldSecret, err := api.Secret.Get(ns, secret.Name, "")
	if err != nil {
		return nil, err
	}
	sv := api.ToSecretView(oldSecret)

	if sv.Equal(api.ToSecretView(secret)) {
		return oldSecret, nil
	}

	secret.Version = oldSecret.Version
	secret.UpdateTimestamp = time.Now()

	res, err := api.Facade.UpdateSecret(ns, secret)
	if err != nil {
		return nil, err
	}
	return api.ToSecretView(res), nil
}

func (api *API) deleteSecret(ns string, r runtime.Object) (string, error) {
	sec, ok := r.(*corev1.Secret)
	if !ok {
		return "", common.Error(common.ErrRequestParamInvalid, common.Field("error", "k8s secret typecasting failed"))
	}

	err := common.ValidateResourceName(sec.Name)
	if err != nil {
		return "", err
	}

	switch sec.Type {
	case corev1.SecretTypeDockerConfigJson:
		_, err = api.DeleteSecretResource(ns, sec.Name, "registry")
		if err != nil {
			return "", err
		}
	case corev1.SecretTypeTLS:
		_, err = api.DeleteSecretResource(ns, sec.Name, "certificate")
		if err != nil {
			return "", err
		}
	default:
		_, err = api.DeleteSecretResource(ns, sec.Name, "secret")
		if err != nil {
			return "", err
		}
	}
	return sec.Name, err
}

func (api *API) generateSecretResource(ns string, sec *corev1.Secret) (*specV1.Secret, error) {
	var err error
	secret := &specV1.Secret{}
	switch sec.Type {
	case corev1.SecretTypeDockerConfigJson:
		registry, err := api.generateRegistry(ns, sec)
		if err != nil {
			return nil, err
		}
		secret = registry.ToSecret()
	case corev1.SecretTypeTLS:
		certificate, err := api.generateCertificate(ns, sec)
		if err != nil {
			return nil, err
		}
		secret = certificate.ToSecret()
	default:
		secret, err = api.generateCommonSecret(ns, sec)
		if err != nil {
			return nil, err
		}
	}

	err = validateSecret(secret)
	if err != nil {
		return nil, err
	}

	return secret, err
}

func (api *API) generateRegistry(ns string, s *corev1.Secret) (*models.Registry, error) {
	registry := new(models.Registry)
	dockercfgjson, ok := s.Data[corev1.DockerConfigJsonKey]
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "no dockerconfigjson found in registry"))
	}
	auths := map[string]interface{}{}
	err := json.Unmarshal(dockercfgjson, &auths)
	if err != nil {
		return nil, err
	}
	if data, ok := auths["auths"]; ok {
		if dd, ok := data.(map[string]interface{}); ok {
			for key, value := range dd {
				registry.Address = key
				registry.Name = s.Name
				registry.Namespace = ns
				if vmap, ok := value.(map[string]interface{}); ok {
					for k, v := range vmap {
						if strings.Contains(strings.ToLower(k), "username") {
							registry.Username = v.(string)
						}
						if strings.Contains(strings.ToLower(k), "password") {
							registry.Password = v.(string)
						}
					}
				}
				if err = api.ValidateRegistryModel(registry); err != nil {
					return nil, err
				}
				return registry, nil
			}
		}
	}

	return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "failed to generete registry"))
}

func (api *API) generateCertificate(ns string, s *corev1.Secret) (*models.Certificate, error) {
	certificate := new(models.Certificate)
	if key, ok := s.Data[corev1.TLSPrivateKeyKey]; ok {
		certificate.Data.Key = string(key)
	}
	if cert, ok := s.Data[corev1.TLSCertKey]; ok {
		certificate.Data.Certificate = string(cert)
	}
	if certificate.Data.Key == "" || certificate.Data.Certificate == "" {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "private key and certificate can't be empty"))
	}
	certificate.Name = s.Name
	certificate.Namespace = ns

	if err := certificate.ParseCertInfo(); err != nil {
		return nil, err
	}

	return certificate, nil
}

func (api *API) generateCommonSecret(ns string, s *corev1.Secret) (*specV1.Secret, error) {
	secret := &specV1.Secret{}
	err := copier.Copy(secret, s)
	if err != nil {
		return nil, err
	}
	secret.Namespace = ns
	secret.Labels = common.AddSystemLabel(secret.Labels, map[string]string{
		specV1.SecretLabel: specV1.SecretConfig,
	})
	return secret, nil
}

func fileCheck(filename string) error {
	fn := strings.Split(filename, ".")
	if fType := fn[len(fn)-1]; fType != "yaml" && fType != "yml" {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", "file type invalid"))
	}
	return nil
}

func validateSecret(s *specV1.Secret) error {
	if s.Name == "" {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", "secret name is required"))
	}
	return common.ValidateResourceName(s.Name)
}

// config resource
func (api *API) generateConfig(ns, userId string, r runtime.Object) (interface{}, error) {
	cfg, ok := r.(*corev1.ConfigMap)
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "k8s config typecasting failed"))
	}

	config := &specV1.Configuration{
		Namespace: ns,
		Name:      cfg.Name,
		Labels:    cfg.Labels,
		Data:      map[string]string{},
	}

	err := generateConfigData(userId, cfg.Data, config.Data)
	if err != nil {
		return nil, err
	}

	err = validateConfig(config)
	if err != nil {
		return nil, err
	}

	oldConfig, err := api.Config.Get(ns, cfg.Name, "")
	if err != nil {
		if e, ok := err.(errors.Coder); !ok || e.Code() != common.ErrResourceNotFound {
			return nil, err
		}
	}

	if oldConfig != nil {
		return nil, common.Error(common.ErrRequestParamInvalid,
			common.Field("error", "this name is already in use"))
	}

	config, err = api.Facade.CreateConfig(ns, config)
	if err != nil {
		return nil, err
	}

	return api.ToConfigurationView(config)
}

func (api *API) updateConfig(ns, userId string, r runtime.Object) (interface{}, error) {
	cfg, ok := r.(*corev1.ConfigMap)
	if !ok {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "k8s config typecasting failed"))
	}

	config := &specV1.Configuration{
		Namespace: ns,
		Name:      cfg.Name,
		Labels:    cfg.Labels,
		Data:      map[string]string{},
	}

	err := generateConfigData(userId, cfg.Data, config.Data)
	if err != nil {
		return nil, err
	}

	err = validateConfig(config)
	if err != nil {
		return nil, err
	}

	res, err := api.Config.Get(ns, cfg.Name, "")
	if err != nil {
		return nil, err
	}

	// labels can't be modified of sys apps
	if CheckIsSysResources(res.Labels) && !reflect.DeepEqual(res.Labels, config.Labels) {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "labels can't be modified of sys apps"))
	}

	if models.EqualConfig(res, config) {
		return res, nil
	}

	config.Version = res.Version
	config.UpdateTimestamp = time.Now()
	config.CreationTimestamp = res.CreationTimestamp

	res, err = api.Facade.UpdateConfig(ns, config)
	if err != nil {
		return nil, err
	}

	return api.ToConfigurationView(res)
}

func (api *API) deleteConfig(ns string, r runtime.Object) (string, error) {
	cfg, ok := r.(*corev1.ConfigMap)
	if !ok {
		return "", common.Error(common.ErrRequestParamInvalid, common.Field("error", "k8s config typecasting failed"))
	}

	err := common.ValidateResourceName(cfg.Name)
	if err != nil {
		return "", err
	}

	res, err := api.Config.Get(ns, cfg.Name, "")
	if err != nil {
		if e, ok := err.(errors.Coder); ok && e.Code() == common.ErrResourceNotFound {
			return "", nil
		}
		api.log.Error("get config failed", log.Error(err), log.Any("name", cfg.Name), log.Any("namespace", ns))
		return "", err
	}

	appNames, err := api.Index.ListAppIndexByConfig(ns, res.Name)
	if err != nil {
		api.log.Error("list app index by config failed", log.Error(err), log.Any("name", cfg.Name), log.Any("namespace", ns))
		return "", err
	}

	if len(appNames) > 0 {
		return "", common.Error(common.ErrResourceHasBeenUsed,
			common.Field("type", "config"),
			common.Field("name", cfg.Name))
	}

	return cfg.Name, api.Facade.DeleteConfig(ns, cfg.Name)
}

func generateConfigData(userId string, cfgData, configData map[string]string) error {
	for name, content := range cfgData {
		vmap := map[string]string{}
		err := yaml.Unmarshal([]byte(content), &vmap)
		if err != nil {
			if strings.HasPrefix(name, common.ConfigObjectPrefix) {
				return common.Error(common.ErrRequestParamInvalid,
					common.Field("error", "key of kv data can't start with "+common.ConfigObjectPrefix))
			}
			configData[name] = content
			continue
		}
		ctype, ok := vmap["type"]
		if !ok {
			if strings.HasPrefix(name, common.ConfigObjectPrefix) {
				return common.Error(common.ErrRequestParamInvalid,
					common.Field("error", "key of kv data can't start with "+common.ConfigObjectPrefix))
			}
			configData[name] = content
			continue
		}

		switch ctype {
		case ConfigTypeObject:
			ok = checkElementsExist(vmap, "source")
			if !ok {
				return common.Error(common.ErrRequestParamInvalid,
					common.Field("error", "failed to validate object data of config"))
			}
			if vmap["source"] == ConfigObjectTypeHttp {
				ok = checkElementsExist(vmap, "url")
				if !ok {
					return common.Error(common.ErrRequestParamInvalid,
						common.Field("error", "failed to validate object data of config"))
				}
			}
			object := &specV1.ConfigurationObject{
				URL:      vmap["url"],
				MD5:      vmap["md5"],
				Unpack:   vmap["unpack"],
				Metadata: map[string]string{},
			}
			object.Metadata = vmap
			object.Metadata["userID"] = userId
			bytes, err := json.Marshal(object)
			if err != nil {
				return err
			}
			configData[common.ConfigObjectPrefix+name] = string(bytes)
		case ConfigTypeFunction:
			ok = checkElementsExist(vmap, "function", "version", "runtime",
				"handler", "bucket", "object")
			if !ok {
				return common.Error(common.ErrRequestParamInvalid,
					common.Field("error", "failed to validate function data of config"))
			}
			object := &specV1.ConfigurationObject{
				URL:      vmap["url"],
				MD5:      vmap["md5"],
				Unpack:   vmap["unpack"],
				Metadata: map[string]string{},
			}
			object.Metadata = vmap
			object.Metadata["userID"] = userId
			bytes, err := json.Marshal(object)
			if err != nil {
				return err
			}
			configData[common.ConfigObjectPrefix+name] = string(bytes)
		default:
			continue
		}
	}
	return nil
}

func validateConfig(c *specV1.Configuration) error {
	if c.Name == "" {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", "config name is required"))
	}
	err := common.ValidateResourceName(c.Name)
	if err != nil {
		return err
	}
	for k, _ := range c.Data {
		err = common.ValidateKeyValue(k)
		if err != nil {
			return err
		}
	}
	return nil
}

// app resource
func (api *API) generateApplication(ns string, r runtime.Object) (interface{}, error) {
	app, err := api.generateAppData(ns, r)
	if err != nil {
		return nil, err
	}

	err = validateApp(app)
	if err != nil {
		return nil, err
	}

	oldApp, err := api.App.Get(ns, app.Name, "")
	if err != nil {
		if e, ok := err.(errors.Coder); !ok || e.Code() != common.ErrResourceNotFound {
			return nil, err
		}
	}
	if oldApp != nil {
		return nil, common.Error(common.ErrResourceHasBeenUsed,
			common.Field("error", "this name is already in use"))
	}

	app, err = api.Facade.CreateApp(ns, nil, app, nil)
	if err != nil {
		return nil, err
	}

	return api.ToApplicationView(app)
}

func (api *API) updateApplication(ns string, r runtime.Object) (interface{}, error) {
	app, err := api.generateAppData(ns, r)
	if err != nil {
		return nil, err
	}

	err = validateApp(app)
	if err != nil {
		return nil, err
	}

	oldApp, err := api.App.Get(ns, app.Name, "")
	if err != nil {
		if e, ok := err.(errors.Coder); !ok || e.Code() != common.ErrResourceNotFound {
			return nil, err
		}
	}
	// sys app: core、init、function is not visible
	if common.ValidIsInvisible(oldApp.Labels) {
		return nil, common.Error(common.ErrResourceInvisible, common.Field("type", common.APP), common.Field("name", oldApp.Name))
	}

	// labels and Selector can't be modified of sys apps
	if CheckIsSysResources(oldApp.Labels) &&
		(oldApp.Selector != app.Selector || !reflect.DeepEqual(oldApp.Labels, app.Labels) || !app.System) {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "selector，labels or system field can't be modified of sys apps"))
	}

	app.Version = oldApp.Version
	app.CreationTimestamp = oldApp.CreationTimestamp
	app.Selector = oldApp.Selector
	app.NodeSelector = oldApp.NodeSelector

	app, err = api.Facade.UpdateApp(ns, oldApp, app, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return api.ToApplicationView(app)
}

func (api *API) deleteApplication(ns string, r runtime.Object) (string, error) {
	var err error
	var name string
	kind := r.GetObjectKind().GroupVersionKind().Kind
	switch kind {
	case TypeDeploy:
		deploy, ok := r.(*appv1.Deployment)
		if !ok {
			return "", common.Error(common.ErrRequestParamInvalid, common.Field("error", "k8s deployment typecasting failed"))
		}
		name = deploy.Name
	case TypeDaemonset:
		ds, ok := r.(*appv1.DaemonSet)
		if !ok {
			return "", common.Error(common.ErrRequestParamInvalid, common.Field("error", "k8s daemonset typecasting failed"))
		}
		name = ds.Name
	case TypeJob:
		job, ok := r.(*batchv1.Job)
		if !ok {
			return "", common.Error(common.ErrRequestParamInvalid, common.Field("error", "k8s job typecasting failed"))
		}
		name = job.Name
	}

	err = common.ValidateResourceName(name)
	if err != nil {
		return "", err
	}

	app, err := api.App.Get(ns, name, "")
	if err != nil {
		if e, ok := err.(errors.Coder); ok && e.Code() == common.ErrResourceNotFound {
			return "", nil
		}
		return "", err
	}

	if canDelete, err := api.IsAppCanDelete(ns, app.Name); err != nil {
		return "", err
	} else if !canDelete {
		return "", common.Error(common.ErrAppReferencedByNode, common.Field("name", app.Name))
	}

	err = api.Facade.DeleteApp(ns, app.Name, app)
	return app.Name, err
}

func (api *API) generateAppData(ns string, r runtime.Object) (*specV1.Application, error) {
	var err error
	app := new(specV1.Application)
	kind := r.GetObjectKind().GroupVersionKind().Kind
	switch kind {
	case TypeDeploy:
		deploy, ok := r.(*appv1.Deployment)
		if !ok {
			return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "k8s deployment typecasting failed"))
		}
		app, err = api.generateDeployApp(ns, deploy)
		if err != nil {
			return nil, err
		}
	case TypeDaemonset:
		ds, ok := r.(*appv1.DaemonSet)
		if !ok {
			return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "k8s daemonset typecasting failed"))
		}
		app, err = api.generateDaemonSetApp(ns, ds)
		if err != nil {
			return nil, err
		}
	case TypeJob:
		job, ok := r.(*batchv1.Job)
		if !ok {
			return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "k8s job typecasting failed"))
		}
		app, err = api.generateJobApp(ns, job)
		if err != nil {
			return nil, err
		}
	}
	return app, nil
}

func (api *API) generateDeployApp(ns string, deploy *appv1.Deployment) (*specV1.Application, error) {
	app, err := api.generateCommonAppInfo(ns, &deploy.Spec.Template.Spec)
	if err != nil {
		return nil, err
	}
	app.Name = deploy.Name
	app.Namespace = ns
	app.Replica = int(*deploy.Spec.Replicas)

	labels := map[string]string{}
	for k, v := range deploy.Labels {
		labels[k] = v
	}
	app.Labels = labels

	app.Type = common.ContainerApp
	app.Mode = context.RunModeKube
	app.Workload = "deployment"

	app.Labels = common.AddSystemLabel(app.Labels, map[string]string{
		common.LabelAppMode: app.Mode,
	})

	return app, nil
}

func (api *API) generateDaemonSetApp(ns string, ds *appv1.DaemonSet) (*specV1.Application, error) {
	app, err := api.generateCommonAppInfo(ns, &ds.Spec.Template.Spec)
	if err != nil {
		return nil, err
	}
	app.Name = ds.Name
	app.Namespace = ns

	labels := map[string]string{}
	for k, v := range ds.Labels {
		labels[k] = v
	}
	app.Labels = labels

	app.Type = common.ContainerApp
	app.Mode = context.RunModeKube
	app.Workload = "daemonset"

	app.Labels = common.AddSystemLabel(app.Labels, map[string]string{
		common.LabelAppMode: app.Mode,
	})

	return app, nil
}

func (api *API) generateJobApp(ns string, job *batchv1.Job) (*specV1.Application, error) {
	app, err := api.generateCommonAppInfo(ns, &job.Spec.Template.Spec)
	if err != nil {
		return nil, err
	}
	app.Name = job.Name
	app.Namespace = ns

	labels := map[string]string{}
	for k, v := range job.Labels {
		labels[k] = v
	}
	app.Labels = labels

	var jobConfig specV1.AppJobConfig
	jobSepc := job.Spec
	if jobSepc.Template.Spec.RestartPolicy != "" {
		jobConfig.RestartPolicy = string(jobSepc.Template.Spec.RestartPolicy)
	}
	if jobSepc.Parallelism != nil {
		jobConfig.Parallelism = int(*jobSepc.Parallelism)
	}
	if jobSepc.Completions != nil {
		jobConfig.Completions = int(*jobSepc.Completions)
	}
	if jobSepc.BackoffLimit != nil {
		jobConfig.BackoffLimit = int(*jobSepc.BackoffLimit)
	}
	app.JobConfig = &jobConfig

	app.Type = common.ContainerApp
	app.Mode = context.RunModeKube
	app.Workload = "job"

	app.Labels = common.AddSystemLabel(app.Labels, map[string]string{
		common.LabelAppMode: app.Mode,
	})

	return app, nil
}

func (api *API) generateCommonAppInfo(ns string, podSpec *corev1.PodSpec) (*specV1.Application, error) {
	var volumes []specV1.Volume
	for _, v := range podSpec.Volumes {
		var vol specV1.Volume
		if v.ConfigMap != nil {
			_, err := api.Config.Get(ns, v.ConfigMap.Name, "")
			if err != nil {
				return nil, err
			}
			vol.Name = v.Name
			vol.VolumeSource.Config = &specV1.ObjectReference{
				Name: v.ConfigMap.Name,
			}
		} else if v.Secret != nil {
			_, err := api.Secret.Get(ns, v.Secret.SecretName, "")
			if err != nil {
				return nil, err
			}
			vol.Name = v.Name
			vol.VolumeSource.Secret = &specV1.ObjectReference{
				Name: v.Secret.SecretName,
			}
		} else if v.HostPath != nil {
			vol.Name = v.Name
			vol.VolumeSource.HostPath = &specV1.HostPathVolumeSource{
				Path: v.HostPath.Path,
			}
			if v.HostPath.Type != nil {
				vol.VolumeSource.HostPath.Type = string(*v.HostPath.Type)
			}
		} else if v.EmptyDir != nil {
			vol.Name = v.Name
			vol.VolumeSource.EmptyDir = &specV1.EmptyDirVolumeSource{
				Medium: string(v.EmptyDir.Medium),
			}
			if v.EmptyDir.SizeLimit != nil {
				vol.VolumeSource.EmptyDir.SizeLimit = v.EmptyDir.SizeLimit.String()
			}
		} else {
			api.log.Warn("volume type unsupported")
			continue
		}
		volumes = append(volumes, vol)
	}

	for _, v := range podSpec.ImagePullSecrets {
		registry, err := api.Secret.Get(ns, v.Name, "")
		if err != nil {
			return nil, err
		}
		if !isRegistrySecret(*registry) {
			return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "imagePullSecrets is not registry type"))
		}
		secretVolume := specV1.Volume{
			Name: v.Name,
			VolumeSource: specV1.VolumeSource{
				Secret: &specV1.ObjectReference{
					Name: v.Name,
				},
			},
		}
		volumes = append(volumes, secretVolume)
	}

	var svcs []specV1.Service
	for _, c := range podSpec.Containers {
		var svc specV1.Service
		err := TransContainerToSvc(&svc, &c)
		if err != nil {
			return nil, err
		}
		svcs = append(svcs, svc)
	}

	var initSvcs []specV1.Service
	for _, c := range podSpec.InitContainers {
		var svc specV1.Service
		err := TransContainerToSvc(&svc, &c)
		if err != nil {
			return nil, err
		}
		initSvcs = append(initSvcs, svc)
	}
	return &specV1.Application{
		InitServices: initSvcs,
		Services:     svcs,
		Volumes:      volumes,
		HostNetwork:  podSpec.HostNetwork,
	}, nil
}

func TransContainerToSvc(svc *specV1.Service, c *corev1.Container) error {
	if err := copier.Copy(&svc, &c); err != nil {
		return errors.Trace(err)
	}
	if c.Resources.Limits != nil {
		svc.Resources.Limits = map[string]string{}
		for k, v := range c.Resources.Limits {
			svc.Resources.Limits[string(k)] = v.String()
		}
	}
	if c.Resources.Requests != nil {
		svc.Resources.Requests = map[string]string{}
		for k, v := range c.Resources.Requests {
			svc.Resources.Requests[string(k)] = v.String()
		}
	}
	if c.SecurityContext != nil {
		svc.SecurityContext.Privileged = *c.SecurityContext.Privileged
	}

	for i, _ := range svc.Ports {
		svc.Ports[i].ServiceType = string(corev1.ServiceTypeClusterIP)
	}

	return nil
}

func isRegistrySecret(secret specV1.Secret) bool {
	registry, ok := secret.Labels[specV1.SecretLabel]
	return ok && registry == specV1.SecretRegistry
}

func validateApp(app *specV1.Application) error {
	if app.Name == "" {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", "app name is required"))
	}
	err := common.ValidateResourceName(app.Name)
	if err != nil {
		return err
	}

	ports := make(map[int32]bool)
	for _, svc := range app.Services {
		err = common.ValidateResourceName(svc.Name)
		if err != nil {
			return nil
		}
		if svc.Image == "" {
			common.Error(common.ErrRequestParamInvalid, common.Field("error", "image is required"))
		}
		if app.Replica > 1 && len(svc.Ports) > 0 {
			return common.Error(common.ErrRequestParamInvalid,
				common.Field("error", "port mapping is only supported under single replica"))
		}
		err = isValidSvcPort(&svc, ports)
		if err != nil {
			return err
		}
	}

	for _, svc := range app.InitServices {
		err = common.ValidateResourceName(svc.Name)
		if err != nil {
			return nil
		}
		if svc.Image == "" {
			common.Error(common.ErrRequestParamInvalid, common.Field("error", "image is required"))
		}
		if app.Replica > 1 && len(svc.Ports) > 0 {
			return common.Error(common.ErrRequestParamInvalid,
				common.Field("error", "port mapping is only supported under single replica"))
		}
		err = isValidSvcPort(&svc, ports)
		if err != nil {
			return err
		}
	}

	for _, v := range app.Volumes {
		err = common.ValidateResourceName(v.Name)
		if err != nil {
			return nil
		}
	}

	return nil
}

func isValidSvcPort(service *specV1.Service, ports map[int32]bool) error {
	for _, port := range service.Ports {
		if port.ServiceType == string(corev1.ServiceTypeNodePort) {
			if port.NodePort <= 0 {
				return common.Error(common.ErrRequestParamInvalid, common.Field("error", "invalid NodePort"))
			}
			if _, ok := ports[port.NodePort]; ok {
				return common.Error(common.ErrRequestParamInvalid, common.Field("error", "duplicate host ports"))
			} else {
				ports[port.NodePort] = true
			}
		} else {
			if port.HostPort == 0 {
				continue
			}
			if _, ok := ports[port.HostPort]; ok {
				return common.Error(common.ErrRequestParamInvalid, common.Field("error", "duplicate host ports"))
			} else {
				ports[port.HostPort] = true
			}
		}
	}
	return nil
}

// service resource
func (api *API) generateService(ns string, r runtime.Object) error {
	var err error
	svc, ok := r.(*corev1.Service)
	if !ok {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", "k8s service typecasting failed"))
	}
	seletor := labels.SelectorFromSet(svc.Spec.Selector).String()

	apps, err := api.App.List(ns, &models.ListOptions{LabelSelector: seletor})
	if err != nil {
		return err
	}

	for _, item := range apps.Items {
		app, err := api.App.Get(ns, item.Name, "")
		if err != nil {
			return err
		}

		updateAppPort(app, svc)

		app, err = api.App.Update(nil, ns, app)
		if err != nil {
			return err
		}
	}
	return err
}

func (api *API) updateService(ns string, r runtime.Object) error {
	var err error
	svc, ok := r.(*corev1.Service)
	if !ok {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", "k8s service typecasting failed"))
	}
	seletor := labels.SelectorFromSet(svc.Spec.Selector).String()

	apps, err := api.App.List(ns, &models.ListOptions{LabelSelector: seletor})
	if err != nil {
		return err
	}
	for _, item := range apps.Items {
		app, err := api.App.Get(ns, item.Name, "")
		if err != nil {
			return err
		}

		updateAppPort(app, svc)

		app, err = api.App.Update(nil, ns, app)
		if err != nil {
			return err
		}
	}
	return err
}

func (api *API) deleteService(ns string, r runtime.Object) (string, error) {
	var err error
	svc, ok := r.(*corev1.Service)
	if !ok {
		return "", common.Error(common.ErrRequestParamInvalid, common.Field("error", "k8s service typecasting failed"))
	}
	seletor := labels.SelectorFromSet(svc.Spec.Selector).String()

	apps, err := api.App.List(ns, &models.ListOptions{LabelSelector: seletor})
	if err != nil {
		return "", err
	}
	for _, item := range apps.Items {
		app, err := api.App.Get(ns, item.Name, "")
		if err != nil {
			return "", err
		}

		resetAppPort(app, svc)

		app, err = api.App.Update(nil, ns, app)
		if err != nil {
			return "", err
		}
	}
	return svc.Name, err
}

func updateAppPort(app *specV1.Application, svc *corev1.Service) {
	for _, sp := range svc.Spec.Ports {
		for i, _ := range app.Services {
			for j, _ := range app.Services[i].Ports {
				if sp.TargetPort.IntVal == app.Services[i].Ports[j].ContainerPort {
					if svc.Spec.Type == corev1.ServiceTypeNodePort {
						app.Services[i].Ports[j].ServiceType = string(corev1.ServiceTypeNodePort)
						app.Services[i].Ports[j].NodePort = sp.NodePort
					}
				}
			}
		}
	}
}

func resetAppPort(app *specV1.Application, svc *corev1.Service) {
	for _, sp := range svc.Spec.Ports {
		for i, _ := range app.Services {
			for j, _ := range app.Services[i].Ports {
				if sp.TargetPort.IntVal == app.Services[i].Ports[j].ContainerPort {
					if svc.Spec.Type == corev1.ServiceTypeNodePort {
						app.Services[i].Ports[j].ServiceType = string(corev1.ServiceTypeClusterIP)
						app.Services[i].Ports[j].NodePort = 0
					}
				}
			}
		}
	}
}
