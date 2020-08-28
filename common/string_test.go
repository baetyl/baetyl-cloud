package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandString(t *testing.T) {
	test1, test2 := RandString(10), RandString(10)
	assert.NotEqual(t, test1, test2)
}
