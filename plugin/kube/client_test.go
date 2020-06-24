package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	p, err := New()
	assert.Error(t, err)
	assert.Nil(t, p)
	assert.EqualError(t, err, "open etc/baetyl/cloud.yml: no such file or directory")
}
