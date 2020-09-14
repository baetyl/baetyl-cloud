package license

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLicense_CheckLicense(t *testing.T) {
	l, err := New()
	assert.NoError(t, err)
	err = l.(*license).CheckLicense()
	assert.NoError(t, err)
}

func TestQuota_GetQuota(t *testing.T) {
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
