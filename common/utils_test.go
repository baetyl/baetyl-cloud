package common

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUUIDPrune(t *testing.T) {
	uuid := UUIDPrune()
	assert.Equal(t, 32, len(uuid))
	assert.NotEqual(t, uuid, UUIDPrune())
	assert.False(t, strings.Contains(uuid, "-"))
}

func TestUUID(t *testing.T) {
	uuid := UUID()
	assert.Equal(t, 36, len(uuid))
	assert.NotEqual(t, uuid, UUID())
}

func TestCompareNumericalString(t *testing.T) {
	a := "123"
	b := "12"
	assert.Equal(t, 1, CompareNumericalString(a, b))

	b = "121"
	assert.Equal(t, 1, CompareNumericalString(a, b))

	b = "123"
	assert.Equal(t, 0, CompareNumericalString(a, b))

	b = "1231"
	assert.Equal(t, -1, CompareNumericalString(a, b))
}

func TestAddSystemLabel(t *testing.T) {
	var labels map[string]string
	infos := map[string]string{
		"a": "b",
		"c": "d",
	}

	labels = AddSystemLabel(labels, infos)
	for k, v := range infos {
		l, ok := labels[k]
		assert.True(t, ok)
		assert.Equal(t, v, l)
	}
}
