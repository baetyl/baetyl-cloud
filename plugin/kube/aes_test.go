package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAES(t *testing.T) {
	testmap := map[string][]byte{
		"key-a": []byte("hello"),
		"key-b": []byte("world"),
	}
	emap, _ := EncryptMap(testmap, []byte("0123456789abcdef"))
	dmap, _ := DecryptMap(emap, []byte("0123456789abcdef"))

	assert.Equal(t, testmap["key-a"], dmap["key-a"])
	assert.Equal(t, testmap["key-b"], dmap["key-b"])

	itext, _ := Encrypt([]byte("hello world"), []byte("0123456789abcdef"))
	otext, _ := Decrypt(itext, []byte("0123456789abcdef"))
	assert.Equal(t, "hello world", string(otext))
}
