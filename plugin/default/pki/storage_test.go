package pki

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewStorage(t *testing.T) {
	cfg := Persistent{
		Kind: "error",
	}
	_, err := NewStorage(cfg)
	assert.Error(t, err, os.ErrInvalid)
}
