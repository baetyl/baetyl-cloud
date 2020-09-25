package api

import (
	"fmt"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

// GetCertificate get a Certificate
func (api *API) GetCertificate(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	res, err := wrapCertificate(api.Secret.Get(ns, n, ""))
	if err != nil {
		return nil, err
	}
	return hideCertKey(res), nil
}

// ListCertificate list Certificate
func (api *API) ListCertificate(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	res, err := wrapCertificateList(api.Secret.List(ns, wrapCertificateListOption(api.parseListOptions(c))))
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	for i := range res.Items {
		hideCertKey(&res.Items[i])
	}
	return res, err
}

// CreateCertificate create one Certificate
func (api *API) CreateCertificate(c *common.Context) (interface{}, error) {
	cfg, err := api.parseAndCheckCertificateModel(c)
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
	if err = cfg.ParseCertInfo(); err != nil {
		return nil, err
	}
	res, err := wrapCertificate(api.Secret.Create(ns, cfg.ToSecret()))
	return hideCertKey(res), err
}

// UpdateCertificate update the Certificate
func (api *API) UpdateCertificate(c *common.Context) (interface{}, error) {
	cfg, err := api.parseAndCheckCertificateModel(c)
	if err != nil {
		return nil, err
	}
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	sd, err := wrapCertificate(api.Secret.Get(ns, n, ""))
	if err != nil {
		return nil, err
	}
	if cfg.Data.Key == "" {
		cfg.Data.Key = sd.Data.Key
	}
	// only edit description by design
	if cfg.Equal(sd) {
		return hideCertKey(sd), nil
	}
	sd.Description = cfg.Description
	sd.UpdateTimestamp = time.Now()
	sd.Data = cfg.Data
	if err = sd.ParseCertInfo(); err != nil {
		return nil, err
	}
	res, err := wrapCertificate(api.Secret.Update(ns, sd.ToSecret()))
	if err != nil {
		return nil, err
	}
	return hideCertKey(res), nil
}

// DeleteCertificate delete the Certificate
func (api *API) DeleteCertificate(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	return api.deleteSecret(ns, n, "certificate")
}

// GetAppByCertificate list app
func (api *API) GetAppByCertificate(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	res, err := wrapCertificate(api.Secret.Get(ns, n, ""))
	if err != nil {
		return nil, err
	}
	return api.listAppBySecret(ns, res.Name)
}

// parseAndCheckCertificateModel parse and check the config model
func (api *API) parseAndCheckCertificateModel(c *common.Context) (*models.Certificate, error) {
	certificate := new(models.Certificate)
	certificate.Name = c.GetNameFromParam()
	err := c.LoadBody(certificate)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	if name := c.GetNameFromParam(); name != "" {
		certificate.Name = name
	}
	if certificate.Name == "" {
		err = common.Error(common.ErrRequestParamInvalid, common.Field("error", "name is required"))
	}
	return certificate, err
}

func hideCertKey(r *models.Certificate) *models.Certificate {
	if r != nil {
		r.Data.Key = ""
	}
	return r
}

func wrapCertificate(s *specV1.Secret, e error) (*models.Certificate, error) {
	if s != nil {
		return models.FromSecretToCertificate(s), e
	}
	return nil, e
}

func wrapCertificateList(s *models.SecretList, e error) (*models.CertificateList, error) {
	if s != nil {
		return models.FromSecretListToCertificateList(s), e
	}
	return nil, e
}

func wrapCertificateListOption(lo *models.ListOptions) *models.ListOptions {
	// TODO 增加type字段代替label标签
	lo.LabelSelector = fmt.Sprintf("%s=%s", specV1.SecretLabel, specV1.SecretCustomCertificate)
	return lo
}
