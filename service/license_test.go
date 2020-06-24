package service

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLicenseService_CheckLicense(t *testing.T) {
	services := InitMockEnvironment(t)
	ls, err := NewLicenseService(services.conf)
	assert.NoError(t, err)

	services.license.EXPECT().CheckLicense().Return(nil)
	err = ls.CheckLicense()
	assert.NoError(t, err)
}

func TestLicenseService_CheckQuota(t *testing.T) {

	namespace := "test"
	services := InitMockEnvironment(t)
	ls, err := NewLicenseService(services.conf)
	assert.NoError(t, err)

	services.license.EXPECT().CheckQuota(namespace, gomock.Any()).Return(nil)
	err = ls.CheckQuota(namespace, func(namespace string) (map[string]int, error) {
		return nil, nil
	})
	assert.NoError(t, err)
}

func TestLicenseService_ProtectCode(t *testing.T) {
	services := InitMockEnvironment(t)
	ls, err := NewLicenseService(services.conf)
	assert.NoError(t, err)
	services.license.EXPECT().ProtectCode().Return(nil)
	err = ls.ProtectCode()
	assert.NoError(t, err)
}
