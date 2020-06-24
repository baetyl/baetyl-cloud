package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandString(t *testing.T) {
	test1, test2 := RandString(10), RandString(10)
	assert.NotEqual(t, test1, test2)
}
