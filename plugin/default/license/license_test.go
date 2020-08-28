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

func TestLicense_CheckQuota(t *testing.T) {
	l, err := New()
	assert.NoError(t, err)
	err = l.(*license).CheckQuota("ns", func(namespace string) (map[string]int, error) {
		return nil, nil
	})
	assert.NoError(t, err)
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
