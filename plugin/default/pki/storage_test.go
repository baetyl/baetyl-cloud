package pki

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	cfg := Persistent{
		Kind: "error",
	}
	_, err := NewStorage(cfg)
	assert.Error(t, err, os.ErrInvalid)
}
