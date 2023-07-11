package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

func TestLicenseService_CheckQuota(t *testing.T) {
	namespace := "default"
	services := InitMockEnvironment(t)
	ls, err := NewQuotaService(services.conf)
	assert.NoError(t, err)
	quotas := map[string]int{
		plugin.QuotaNode: 10,
	}

	services.quota.EXPECT().GetQuota(namespace).Return(quotas, nil)
	err = ls.CheckQuota(namespace, func(namespace string) (map[string]int, error) {
		return map[string]int{plugin.QuotaNode: 1}, nil
	})
	assert.NoError(t, err)

	services.quota.EXPECT().GetQuota(namespace).Return(nil, nil)
	err = ls.CheckQuota(namespace, func(namespace string) (map[string]int, error) {
		return map[string]int{plugin.QuotaNode: 1}, nil
	})
	assert.NoError(t, err)

	services.quota.EXPECT().GetQuota(namespace).Return(quotas, nil)
	err = ls.CheckQuota(namespace, func(namespace string) (map[string]int, error) {
		return nil, nil
	})
	assert.NoError(t, err)

	errGetQuota := fmt.Errorf("get quota error")
	services.quota.EXPECT().GetQuota(namespace).Return(nil, errGetQuota)
	err = ls.CheckQuota(namespace, func(namespace string) (map[string]int, error) {
		return map[string]int{plugin.QuotaNode: 1}, nil
	})
	assert.Equal(t, err, errGetQuota)

	errGetNodeCount := fmt.Errorf("get node count error")
	services.quota.EXPECT().GetQuota(namespace).Return(quotas, nil)
	err = ls.CheckQuota(namespace, func(namespace string) (map[string]int, error) {
		return nil, errGetNodeCount
	})
	assert.Equal(t, err, errGetNodeCount)

	services.quota.EXPECT().GetQuota(namespace).Return(quotas, nil)
	err = ls.CheckQuota(namespace, func(namespace string) (map[string]int, error) {
		return map[string]int{plugin.QuotaNode: 11}, nil
	})
	assert.Error(t, err)
}
