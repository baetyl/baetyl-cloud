package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

// TODO: optimize this layer, general abstraction

// GetRegistry get a Registry
func (api *API) GetRegistry(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	secret, err := api.Secret.Get(ns, n, "")
	if err != nil {
		return nil, wrapSecretLikedResourceNotFoundError(n, common.Registry, err)
	}

	return hidePwd(api.ToRegistryView(secret)), nil
}

// ListRegistry list Registry
func (api *API) ListRegistry(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	params, err := api.parseListOptionsAppendSystemLabel(c)
	if err != nil {
		return nil, err
	}
	params.LabelSelector += "," + fmt.Sprintf("%s=%s", specV1.SecretLabel, specV1.SecretRegistry)
	secrets, err := api.Secret.List(ns, params)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}

	return api.ToFilteredRegistryViewList(secrets), nil
}

// CreateRegistry create one Registry
func (api *API) CreateRegistry(c *common.Context) (interface{}, error) {
	cfg, err := api.parseAndCheckRegistryModel(c)
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
	if err = api.validateRegistryModel(cfg); err != nil {
		return nil, err
	}

	secret, err := api.Secret.Create(ns, cfg.ToSecret())
	if err != nil {
		return nil, err
	}
	return hidePwd(api.ToFilteredRegistryView(secret)), nil
}

// UpdateRegistry update the Registry
func (api *API) UpdateRegistry(c *common.Context) (interface{}, error) {
	cfg, err := api.parseAndCheckRegistryModel(c)
	if err != nil {
		return nil, err
	}
	ns, n := c.GetNamespace(), c.GetNameFromParam()

	secret, err := api.Secret.Get(ns, n, "")
	if err != nil {
		return nil, wrapSecretLikedResourceNotFoundError(n, common.Registry, err)
	}

	sd := api.ToRegistryView(secret)
	// only edit description by design
	if cfg.Description == sd.Description {
		return hidePwd(sd), nil
	}
	sd.Description = cfg.Description
	sd.UpdateTimestamp = time.Now()
	if err = api.validateRegistryModel(sd); err != nil {
		return nil, err
	}
	secret, err = api.Secret.Update(ns, sd.ToSecret())
	if err != nil {
		return nil, err
	}
	return hidePwd(api.ToRegistryView(secret)), nil
}

func (api *API) RefreshRegistryPassword(c *common.Context) (interface{}, error) {
	cfg, err := api.parseAndCheckRegistryModel(c)
	if err != nil {
		return nil, err
	}
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	secret, err := api.Secret.Get(ns, n, "")
	if err != nil {
		return nil, wrapSecretLikedResourceNotFoundError(n, common.Registry, err)
	}

	sd := api.ToRegistryView(secret)
	sd.UpdateTimestamp = time.Now()
	sd.Password = cfg.Password

	secret, err = api.Secret.Update(ns, sd.ToSecret())
	if err != nil {
		return nil, err
	}
	res := api.ToRegistryView(secret)
	if err = api.validateRegistryModel(res); err != nil {
		return nil, err
	}
	err = api.updateAppSecret(ns, secret)
	if err != nil {
		return nil, err
	}
	return hidePwd(res), nil
}

// DeleteRegistry delete the Registry
func (api *API) DeleteRegistry(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	return api.deleteSecret(ns, n, "registry")
}

// GetAppByRegistry list app
func (api *API) GetAppByRegistry(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	secret, err := api.Secret.Get(ns, n, "")
	if err != nil {
		return nil, wrapSecretLikedResourceNotFoundError(n, common.Registry, err)
	}
	return api.listAppBySecret(ns, secret.Name)
}

// parseAndCheckRegistryModel parse and check the config model
func (api *API) parseAndCheckRegistryModel(c *common.Context) (*models.Registry, error) {
	registry := new(models.Registry)
	registry.Name = c.GetNameFromParam()
	registry.Namespace = c.GetNamespace()
	err := c.LoadBody(registry)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	if name := c.GetNameFromParam(); name != "" {
		registry.Name = name
	}
	if registry.Name == "" {
		err = common.Error(common.ErrRequestParamInvalid, common.Field("error", "name is required"))
	}
	return registry, err
}

func hidePwd(r *models.Registry) *models.Registry {
	if r != nil {
		r.Password = ""
	}
	return r
}

func (api *API) ToFilteredRegistryView(s *specV1.Secret) *models.Registry {
	return models.FromSecretToRegistry(s, true)
}

func (api *API) ToRegistryView(s *specV1.Secret) *models.Registry {
	return models.FromSecretToRegistry(s, false)
}

func (api *API) ToFilteredRegistryViewList(s *models.SecretList) *models.RegistryList {
	res := models.FromSecretListToRegistryList(s, true)
	for i := range res.Items {
		hidePwd(&res.Items[i])
	}
	return res
}

func (api *API) ToRegistryViewList(s *models.SecretList) *models.RegistryList {
	res := models.FromSecretListToRegistryList(s, false)
	for i := range res.Items {
		hidePwd(&res.Items[i])
	}
	return res
}

// validateRegistryModel validate the registry model
func (api *API) validateRegistryModel(r *models.Registry) error {
	if r.Address == "" || r.Username == "" || r.Password == "" {
		return common.Error(common.ErrRequestParamInvalid, common.Field("error", "address/username/password is required"))
	}
	return nil
}

func wrapSecretLikedResourceNotFoundError(name string, secretType common.Resource, err error) error {
	if err != nil {
		e, ok := err.(errors.Coder)
		if (ok && e.Code() == common.ErrResourceNotFound) || (!ok && strings.Contains(err.Error(), "not found")) {
			return common.Error(common.ErrResourceNotFound, common.Field("type", secretType),
				common.Field("name", name))
		}
	}
	return err
}
