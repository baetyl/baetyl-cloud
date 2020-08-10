package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalDefault(t *testing.T) {
	assert.Equal(t, "etc/baetyl/cloud.yml", GetConfFile())
	assert.Equal(t, "requestId", GetTraceKey())
	assert.Equal(t, "x-bce-request-id", GetTraceHeader())

	SetConfFile("a.log")
	SetTraceKey("b")
	SetTraceHeader("c")

	assert.Equal(t, "a.log", GetConfFile())
	assert.Equal(t, "b", GetTraceKey())
	assert.Equal(t, "c", GetTraceHeader())
}
