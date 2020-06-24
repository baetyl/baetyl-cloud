package pki

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewPVC(t *testing.T) {
	cfg := Persistent{
		Kind: "error",
	}
	_, err := NewPVC(cfg)
	assert.Error(t, err, os.ErrInvalid)
}
