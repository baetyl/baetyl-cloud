package api

import (
	"fmt"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

// TODO: optimize this layer, general abstraction

// GetSecret get a secret
func (api *API) GetSecret(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	res, err := wrapSecret(api.Secret.Get(ns, n, ""))
	if res == nil {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "secret"), common.Field("name", n))
	}
	return res, err
}

// ListSecret list secret
func (api *API) ListSecret(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	res, err := wrapSecretList(api.Secret.List(ns, api.parseListOptionsAppendSystemLabel(c)))
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	return res, err
}

// CreateSecret create one secret
func (api *API) CreateSecret(c *common.Context) (interface{}, error) {
	cfg, err := api.parseAndCheckSecretModel(c)
	if err != nil {
		return nil, err
	}
	ns, name := c.GetNamespace(), cfg.Name
	sd, err := api.Secret.Get(ns, name, "")
	if err != nil {
		if e, ok := err.(errors.Coder); !ok || e.Code() != common.ErrResourceNotFound {
			return nil, err
		}
	}
	if sd != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "this name is already in use"))
	}
	return wrapSecret(api.Secret.Create(ns, cfg.ToSecret()))
}

// UpdateSecret update the secret
func (api *API) UpdateSecret(c *common.Context) (interface{}, error) {
	cfg, err := api.parseAndCheckSecretModel(c)
	if err != nil {
		return nil, err
	}
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	sd, err := wrapSecret(api.Secret.Get(ns, n, ""))
	if err != nil {
		return nil, err
	}
	if sd.Equal(cfg) {
		return sd, nil
	}
	cfg.Version = sd.Version
	cfg.UpdateTimestamp = time.Now()
	res, err := wrapSecret(api.Secret.Update(ns, cfg.ToSecret()))
	if err != nil {
		return nil, err
	}
	err = api.updateAppSecret(ns, res.ToSecret())
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteSecret delete the secret
func (api *API) DeleteSecret(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	return api.deleteSecret(ns, n, "secret")
}

// GetAppBySecret list app
func (api *API) GetAppBySecret(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	res, err := wrapSecret(api.Secret.Get(ns, n, ""))
	if err != nil {
		return nil, err
	}
	return api.listAppBySecret(ns, res.Name)
}

// parseAndCheckSecretModel parse and check the config model
func (api *API) parseAndCheckSecretModel(c *common.Context) (*models.SecretView, error) {
	secret := new(models.SecretView)
	secret.Name = c.GetNameFromParam()
	err := c.LoadBody(secret)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	if name := c.GetNameFromParam(); name != "" {
		secret.Name = name
	}
	if secret.Name == "" {
		err = common.Error(common.ErrRequestParamInvalid, common.Field("error", "name is required"))
	}

	return secret, err
}

func wrapSecret(s *specV1.Secret, e error) (*models.SecretView, error) {
	if s != nil {
		return models.FromSecretToView(s), e
	}
	return nil, e
}

func wrapSecretList(s *models.SecretList, e error) (*models.SecretViewList, error) {
	if s != nil {
		return models.FromSecretListToView(s), e
	}
	return nil, e
}

func wrapSecretListOption(lo *models.ListOptions) *models.ListOptions {
	lo.LabelSelector = fmt.Sprintf("%s=%s", specV1.SecretLabel, specV1.SecretConfig)
	return lo
}

func (api *API) deleteSecret(namespace, secret, secretType string) (interface{}, error) {
	appNames, err := api.indexService.ListAppIndexBySecret(namespace, secret)
	if err != nil {
		return nil, err
	}
	if len(appNames) > 0 {
		return nil, common.Error(common.ErrResourceHasBeenUsed, common.Field("type", secretType), common.Field("name", secret))
	}
	return nil, api.Secret.Delete(namespace, secret)
}

func (api *API) updateAppSecret(namespace string, secret *specV1.Secret) error {
	appNames, err := api.indexService.ListAppIndexBySecret(namespace, secret.Name)
	if err != nil {
		return err
	}
	for _, appName := range appNames {
		app, err := api.App.Get(namespace, appName, "")
		if err != nil {
			if e, ok := err.(errors.Coder); ok && e.Code() == common.ErrResourceNotFound {
				continue
			}
			return err
		}
		if !needUpdateAppSecret(secret, app) {
			continue
		}
		app, err = api.App.Update(namespace, app)
		if err != nil {
			return err
		}
		_, err = api.nodeService.UpdateNodeAppVersion(namespace, app)
		if err != nil {
			return err
		}
	}
	return nil
}

func needUpdateAppSecret(secret *specV1.Secret, app *specV1.Application) bool {
	appNeedUpdate := false
	for _, volume := range app.Volumes {
		if volume.Secret != nil &&
			volume.Secret.Name == secret.Name &&
			// secret's version must increment
			common.CompareNumericalString(secret.Version, volume.Secret.Version) > 0 {
			volume.Secret.Version = secret.Version
			appNeedUpdate = true
		}
	}
	return appNeedUpdate
}

func (api *API) listAppBySecret(namespace, secret string) (*models.ApplicationList, error) {
	appNames, err := api.indexService.ListAppIndexBySecret(namespace, secret)
	if err != nil {
		return nil, err
	}
	return api.listAppByNames(namespace, appNames)
}

func (api *API) listAppByConfig(namespace, config string) (*models.ApplicationList, error) {
	appNames, err := api.indexService.ListAppIndexByConfig(namespace, config)
	if err != nil {
		return nil, err
	}
	return api.listAppByNames(namespace, appNames)
}

func (api *API) listAppByNames(namespace string, appNames []string) (*models.ApplicationList, error) {
	result := &models.ApplicationList{
		Total: 0,
		Items: []models.AppItem{},
	}
	for _, appName := range appNames {
		app, err := api.App.Get(namespace, appName, "")
		if err != nil {
			if e, ok := err.(errors.Coder); ok && e.Code() == common.ErrResourceNotFound {
				continue
			}
			return nil, err
		}
		result.Total++
		result.Items = append(result.Items, models.AppItem{
			Name:              app.Name,
			Labels:            app.Labels,
			Version:           app.Version,
			Namespace:         app.Namespace,
			Selector:          app.Selector,
			CreationTimestamp: app.CreationTimestamp,
			Description:       app.Description,
		})
	}
	return result, nil
}
