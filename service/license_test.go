package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLicenseService_CheckLicense(t *testing.T) {
	services := InitMockEnvironment(t)
	ls, err := NewLicenseService(services.conf)
	assert.NoError(t, err)

	services.license.EXPECT().CheckLicense().Return(nil)
	err = ls.CheckLicense()
	assert.NoError(t, err)
}
