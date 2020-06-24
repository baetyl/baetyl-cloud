package pki

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestToCertModel(t *testing.T) {
	c := &Cert{
		CertId:      "12345678",
		ParentId:    "",
		Type:        "ROOT",
		CommonName:  "bie",
		Description: "desc",
		NotBefore:   time.Now(),
		NotAfter:    time.Now(),
	}

	// bad case 0
	c.Csr = "error"
	res := ToCertModel(c)
	assert.Nil(t, res)

	// bad case 1
	c.Csr = "MQ==" // base64(1) -> MQ==
	c.Content = "error"
	res = ToCertModel(c)
	assert.Nil(t, res)

	// bad case 2
	c.Csr = "MQ=="
	c.Content = "MQ=="
	c.PrivateKey = "error"
	res = ToCertModel(c)
	assert.Nil(t, res)

	// good case
	c.Csr = "MQ=="
	c.Content = "MQ=="
	c.PrivateKey = "MQ=="
	res = ToCertModel(c)
	assert.NotNil(t, res)
}
