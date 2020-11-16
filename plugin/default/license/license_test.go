package license

import (
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLicense_CheckLicense(t *testing.T) {
	l, err := New()
	assert.NoError(t, err)
	err = l.(*license).CheckLicense()
	assert.NoError(t, err)
}

func TestLicense_GetQuota(t *testing.T) {
	l, err := New()
	assert.NoError(t, err)
	quotas, err := l.(*license).GetQuota("ns")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(quotas))
}

func TestLicense_ProtectCode(t *testing.T) {
	l, err := New()
	assert.NoError(t, err)
	err = l.(*license).ProtectCode()
	assert.NoError(t, err)
}

func TestLicense_Close(t *testing.T) {
	l, err := New()
	assert.NoError(t, err)
	err = l.(*license).Close()
	assert.NoError(t, err)
}

func TestLicense_AcquireQuota(t *testing.T) {
	l, err := New()
	assert.NoError(t, err)
	namespace := "default"
	err = l.(*license).AcquireQuota(namespace, plugin.QuotaNode, 1)
	assert.NoError(t, err)
}

func TestLicense_CreateQuota(t *testing.T) {
	l, err := New()
	assert.NoError(t, err)
	namespace := "default"
	quotas := map[string]int{plugin.QuotaNode: 10}
	err = l.(*license).CreateQuota(namespace, quotas)
	assert.NoError(t, err)
}

func TestLicense_DeleteQuota(t *testing.T) {
	l, err := New()
	assert.NoError(t, err)
	namespace := "default"
	err = l.(*license).DeleteQuota(namespace, plugin.QuotaNode)
	assert.NoError(t, err)
}

func TestLicense_DeleteQuotaByNamespace(t *testing.T) {
	l, err := New()
	assert.NoError(t, err)
	namespace := "default"
	err = l.(*license).DeleteQuotaByNamespace(namespace)
	assert.NoError(t, err)
}

func TestLicense_GetDefaultQuotas(t *testing.T) {
	l, err := New()
	assert.NoError(t, err)
	namespace := "default"
	res, err := l.(*license).GetDefaultQuotas(namespace)
	assert.NoError(t, err)
	assert.Equal(t, map[string]int{}, res)
}

func TestLicense_ReleaseQuota(t *testing.T) {
	l, err := New()
	assert.NoError(t, err)
	namespace := "default"
	err = l.(*license).ReleaseQuota(namespace, plugin.QuotaNode, 1)
	assert.NoError(t, err)
}

func TestLicense_UpdateQuota(t *testing.T) {
	l, err := New()
	assert.NoError(t, err)
	namespace := "default"
	err = l.(*license).UpdateQuota(namespace, plugin.QuotaNode, 1)
	assert.NoError(t, err)
}
