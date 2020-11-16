package service

import (
	"fmt"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
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
	namespace := "default"
	services := InitMockEnvironment(t)
	ls, err := NewLicenseService(services.conf)
	assert.NoError(t, err)
	quotas := map[string]int{
		plugin.QuotaNode: 10,
	}

	services.license.EXPECT().GetQuota(namespace).Return(quotas, nil)
	err = ls.CheckQuota(namespace, func(namespace string) (map[string]int, error) {
		return map[string]int{plugin.QuotaNode: 1}, nil
	})
	assert.NoError(t, err)

	services.license.EXPECT().GetQuota(namespace).Return(nil, nil)
	err = ls.CheckQuota(namespace, func(namespace string) (map[string]int, error) {
		return map[string]int{plugin.QuotaNode: 1}, nil
	})
	assert.NoError(t, err)

	services.license.EXPECT().GetQuota(namespace).Return(quotas, nil)
	err = ls.CheckQuota(namespace, func(namespace string) (map[string]int, error) {
		return nil, nil
	})
	assert.NoError(t, err)

	errGetQuota := fmt.Errorf("get quota error")
	services.license.EXPECT().GetQuota(namespace).Return(nil, errGetQuota)
	err = ls.CheckQuota(namespace, func(namespace string) (map[string]int, error) {
		return map[string]int{plugin.QuotaNode: 1}, nil
	})
	assert.Equal(t, err, errGetQuota)

	errGetNodeCount := fmt.Errorf("get node count error")
	services.license.EXPECT().GetQuota(namespace).Return(quotas, nil)
	err = ls.CheckQuota(namespace, func(namespace string) (map[string]int, error) {
		return nil, errGetNodeCount
	})
	assert.Equal(t, err, errGetNodeCount)

	services.license.EXPECT().GetQuota(namespace).Return(quotas, nil)
	err = ls.CheckQuota(namespace, func(namespace string) (map[string]int, error) {
		return map[string]int{plugin.QuotaNode: 11}, nil
	})
	assert.Error(t, err)

}

func TestLicenseService_ProtectCode(t *testing.T) {
	services := InitMockEnvironment(t)
	ls, err := NewLicenseService(services.conf)
	assert.NoError(t, err)
	services.license.EXPECT().ProtectCode().Return(nil)
	err = ls.ProtectCode()
	assert.NoError(t, err)
}

func TestLicenseService_GetQuota(t *testing.T) {
	namespace := "default"
	services := InitMockEnvironment(t)
	ls, err := NewLicenseService(services.conf)
	assert.NoError(t, err)
	quotas := map[string]int{
		plugin.QuotaNode: 10,
	}
	services.license.EXPECT().GetQuota(namespace).Return(quotas, nil)
	result, err := ls.GetQuota(namespace)

	assert.NoError(t, err)
	assert.Equal(t, quotas, result)
}

func TestLicenseService_AcquireQuota(t *testing.T) {
	namespace := "default"
	number := 1
	services := InitMockEnvironment(t)
	ls, err := NewLicenseService(services.conf)
	assert.NoError(t, err)

	services.license.EXPECT().AcquireQuota(namespace, plugin.QuotaNode, number).Return(nil)
	err = ls.AcquireQuota(namespace, plugin.QuotaNode, number)

	assert.NoError(t, err)

}

func TestLicenseService_CreateQuota(t *testing.T) {
	namespace := "default"
	number := 10
	services := InitMockEnvironment(t)
	ls, err := NewLicenseService(services.conf)
	assert.NoError(t, err)
	quotas := map[string]int{
		plugin.QuotaNode: number,
	}
	services.license.EXPECT().CreateQuota(namespace, quotas).Return(nil)
	err = ls.CreateQuota(namespace, quotas)

	assert.NoError(t, err)
}

func TestLicenseService_DeleteQuota(t *testing.T) {
	namespace := "default"
	services := InitMockEnvironment(t)
	ls, err := NewLicenseService(services.conf)
	assert.NoError(t, err)

	services.license.EXPECT().DeleteQuota(namespace, plugin.QuotaNode).Return(nil)
	err = ls.DeleteQuota(namespace, plugin.QuotaNode)

	assert.NoError(t, err)
}

func TestLicenseService_DeleteQuotaByNamespace(t *testing.T) {
	namespace := "default"
	services := InitMockEnvironment(t)
	ls, err := NewLicenseService(services.conf)
	assert.NoError(t, err)

	services.license.EXPECT().DeleteQuotaByNamespace(namespace).Return(nil)
	err = ls.DeleteQuotaByNamespace(namespace)

	assert.NoError(t, err)
}

func TestLicenseService_GetDefaultQuotas(t *testing.T) {
	namespace := "default"
	services := InitMockEnvironment(t)
	ls, err := NewLicenseService(services.conf)
	assert.NoError(t, err)
	quotas := map[string]int{
		plugin.QuotaNode: 10,
	}
	services.license.EXPECT().GetDefaultQuotas(namespace).Return(quotas, nil)
	result, err := ls.GetDefaultQuotas(namespace)

	assert.NoError(t, err)
	assert.Equal(t, quotas, result)
}

func TestLicenseService_ReleaseQuota(t *testing.T) {
	namespace := "default"
	number := 1
	services := InitMockEnvironment(t)
	ls, err := NewLicenseService(services.conf)
	assert.NoError(t, err)

	services.license.EXPECT().ReleaseQuota(namespace, plugin.QuotaNode, number).Return(nil)
	err = ls.ReleaseQuota(namespace, plugin.QuotaNode, number)

	assert.NoError(t, err)
}

func TestLicenseService_UpdateQuota(t *testing.T) {
	namespace := "default"
	quota := 10
	services := InitMockEnvironment(t)
	ls, err := NewLicenseService(services.conf)
	assert.NoError(t, err)

	services.license.EXPECT().UpdateQuota(namespace, plugin.QuotaNode, quota).Return(nil)
	err = ls.UpdateQuota(namespace, plugin.QuotaNode, quota)

	assert.NoError(t, err)
}
