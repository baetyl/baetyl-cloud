package service

import (
	"os"
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPkiService_NewPKIService(t *testing.T) {
	conf := &config.CloudConfig{}
	conf.Plugin.PKI = "dummy-" + t.Name()
	_, err := NewPKIService(conf)
	assert.Error(t, err)
}

func TestPkiService_GetCA(t *testing.T) {
	mc := InitMockEnvironment(t)
	defer mc.Close()

	ps, err := NewPKIService(mc.conf)
	assert.NoError(t, err)

	mc.pki.EXPECT().GetRootCertId().Return("123").Times(1)
	mc.pki.EXPECT().GetRootCert("123").Return([]byte("test"), nil).Times(1)
	res, err := ps.GetCA()
	assert.NoError(t, err)
	assert.Equal(t, "test", string(res))
}

func TestPkiService_SignClientCertificate(t *testing.T) {
	mc := InitMockEnvironment(t)
	defer mc.Close()

	cn := "test"
	altNames := models.AltNames{}
	certId := "132"
	certPem := []byte("pem")

	ps, err := NewPKIService(mc.conf)
	assert.NoError(t, err)
	rootId := "12345678"

	// good case
	mc.pki.EXPECT().GetRootCertId().Return(rootId).Times(3)
	mc.pki.EXPECT().CreateClientCert(gomock.Any(), rootId).Return(certId, nil).Times(1)
	mc.pki.EXPECT().GetClientCert(certId).Return(certPem, nil).Times(1)
	res, err := ps.SignClientCertificate(cn, altNames)
	assert.NoError(t, err)
	assert.Equal(t, certPem, res.CertPEM)

	//bad case 1
	mc.pki.EXPECT().CreateClientCert(gomock.Any(), rootId).Return("", os.ErrNotExist).Times(1)
	res, err = ps.SignClientCertificate(cn, altNames)
	assert.Error(t, err)

	//bad case 2
	mc.pki.EXPECT().CreateClientCert(gomock.Any(), rootId).Return(certId, nil).Times(1)
	mc.pki.EXPECT().GetClientCert(certId).Return(nil, os.ErrNotExist).Times(1)
	res, err = ps.SignClientCertificate(cn, altNames)
	assert.Error(t, err)
}

func TestPkiService_SignServerCertificate(t *testing.T) {
	mc := InitMockEnvironment(t)
	defer mc.Close()

	cn := "test"
	altNames := models.AltNames{}
	rootId := "12345678"
	certId := "132"
	certPem := []byte("pem")

	ps, err := NewPKIService(mc.conf)
	assert.NoError(t, err)

	// good case
	mc.pki.EXPECT().GetRootCertId().Return(rootId).Times(1)
	mc.pki.EXPECT().CreateServerCert(gomock.Any(), rootId).Return(certId, nil).Times(1)
	mc.pki.EXPECT().GetServerCert(certId).Return(certPem, nil).Times(1)
	res, err := ps.SignServerCertificate(cn, altNames)
	assert.NoError(t, err)
	assert.Equal(t, certPem, res.CertPEM)
}

func TestPkiService_DeleteCertificate(t *testing.T) {
	mc := InitMockEnvironment(t)
	defer mc.Close()

	certId := "12345678"

	ps, err := NewPKIService(mc.conf)
	assert.NoError(t, err)

	mc.pki.EXPECT().DeleteClientCert(certId).Return(nil).Times(1)
	err = ps.DeleteClientCertificate(certId)
	assert.NoError(t, err)

	mc.pki.EXPECT().DeleteServerCert(certId).Return(nil).Times(1)
	err = ps.DeleteServerCertificate(certId)
	assert.NoError(t, err)
}
