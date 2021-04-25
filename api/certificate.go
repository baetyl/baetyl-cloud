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

	secret, err := api.Secret.Get(ns, n, "")
	if err != nil {
		return nil, wrapSecretLikedResourceNotFoundError(n, common.Certificate, err)
	}
	return hideCertKey(api.ToCertificateView(secret)), nil
}

// ListCertificate list Certificate
func (api *API) ListCertificate(c *common.Context) (interface{}, error) {
	ns := c.GetNamespace()
	params, err := api.parseListOptionsAppendSystemLabel(c)
	if err != nil {
		return nil, err
	}
	params.LabelSelector += "," + fmt.Sprintf("%s=%s", specV1.SecretLabel, specV1.SecretCertificate)
	secrets, err := api.Secret.List(ns, params)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}

	return api.ToFilteredCertificateViewList(secrets), nil
}

// CreateCertificate create one Certificate
func (api *API) CreateCertificate(c *common.Context) (interface{}, error) {
	cfg, err := parseAndCheckCertificateModelWhenGet(c)
	if err != nil {
		return nil, err
	}

	ns, name := c.GetNamespace(), cfg.Name
	cfg.Namespace = ns

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
	res, err := api.Secret.Create(nil, ns, cfg.ToSecret())
	if err != nil {
		return nil, err
	}

	return hideCertKey(api.ToFilteredCertificateView(res)), err
}

// UpdateCertificate update the Certificate
func (api *API) UpdateCertificate(c *common.Context) (interface{}, error) {
	cfg, err := parseAndCheckCertificateModelWhenUpdate(c)
	if err != nil {
		return nil, err
	}

	ns, n := c.GetNamespace(), c.GetNameFromParam()
	secret, err := api.Secret.Get(ns, n, "")
	if err != nil {
		return nil, wrapSecretLikedResourceNotFoundError(n, common.Certificate, err)
	}

	sd := api.ToCertificateView(secret)
	if cfg.Data.Key == "" {
		cfg.Data = sd.Data
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
	secret, err = api.Secret.Update(ns, sd.ToSecret())
	if err != nil {
		return nil, err
	}
	return hideCertKey(api.ToCertificateView(secret)), nil
}

// DeleteCertificate delete the Certificate
func (api *API) DeleteCertificate(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	return api.deleteSecret(ns, n, "certificate")
}

// GetAppByCertificate list app
func (api *API) GetAppByCertificate(c *common.Context) (interface{}, error) {
	ns, n := c.GetNamespace(), c.GetNameFromParam()
	res, err := api.Secret.Get(ns, n, "")
	if err != nil {
		return nil, wrapSecretLikedResourceNotFoundError(n, common.Certificate, err)
	}
	return api.listAppBySecret(ns, res.Name)
}

func parseAndCheckCertificateModelWhenGet(c *common.Context) (*models.Certificate, error) {
	cert, err := parseAndCheckCertificateModel(c)
	if err != nil {
		return nil, err
	}
	if cert.Data.Key == "" || cert.Data.Certificate == "" {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "private key and certificate can't be empty"))
	}
	return cert, nil
}

func parseAndCheckCertificateModelWhenUpdate(c *common.Context) (*models.Certificate, error) {
	cert, err := parseAndCheckCertificateModel(c)
	if err != nil {
		return nil, err
	}
	if cert.Data.Key != "" && cert.Data.Certificate == "" {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "private key and certificate can't be empty"))
	}
	return cert, nil
}

func parseAndCheckCertificateModel(c *common.Context) (*models.Certificate, error) {
	certificate := new(models.Certificate)
	certificate.Name = c.GetNameFromParam()
	certificate.Namespace = c.GetNamespace()
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

func (api *API) ToFilteredCertificateView(s *specV1.Secret) *models.Certificate {
	return models.FromSecretToCertificate(s, true)
}

func (api *API) ToCertificateView(s *specV1.Secret) *models.Certificate {
	return models.FromSecretToCertificate(s, false)
}

func (api *API) ToFilteredCertificateViewList(s *models.SecretList) *models.CertificateList {
	res := models.FromSecretListToCertificateList(s, true)
	for i := range res.Items {
		hideCertKey(&res.Items[i])
	}
	return res
}

func (api *API) ToCertificateViewList(s *models.SecretList) *models.CertificateList {
	res := models.FromSecretListToCertificateList(s, false)
	for i := range res.Items {
		hideCertKey(&res.Items[i])
	}
	return res
}
