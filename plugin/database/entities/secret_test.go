// Package entities 数据库存储基本结构与方法
package entities

import (
	"testing"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

func TestFromSecretModel(t *testing.T) {
	secret := &Secret{
		Name:        "testSecret",
		Namespace:   "namespace",
		Labels:      `{"baetyl-cloud-system":"true"}`,
		Data:        "{\"cfg1\":\"MTIz\",\"cfg2\":\"YWJj\"}",
		Description: "desc",
	}

	mSecret := &specV1.Secret{
		Namespace: "namespace",
		Name:      "testSecret",
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Data: map[string][]byte{
			"cfg1": []byte("123"),
			"cfg2": []byte("abc"),
		},
		Description: "desc",
	}
	res, err := FromSecretModel("namespace", mSecret)

	assert.NoError(t, err)
	assert.Equal(t, secret.Name, res.Name)
	assert.Equal(t, secret.Namespace, res.Namespace)
	assert.Equal(t, secret.Labels, res.Labels)
	assert.Equal(t, secret.Description, res.Description)
	assert.Equal(t, secret.Data, res.Data)

	// size too large
	largeSert := &specV1.Secret{
		Namespace: "namespace",
		Name:      "largeSert",
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Data:        map[string][]byte{},
		Description: "desc",
	}

	data := ""
	for i := 0; i < 64*1024; i++ {
		data += "1"
	}
	largeSert.Data["data"] = []byte(data)
	_, err = FromSecretModel("namespace", largeSert)
	assert.Error(t, err)
}

func TestToSecretModel(t *testing.T) {
	secret := &Secret{
		ID:          123,
		Name:        "testSecret",
		Namespace:   "namespace",
		Labels:      "{\"baetyl-cloud-system\":\"true\"}",
		Data:        "{\"cfg1\":\"MTIz\",\"cfg2\":\"YWJj\"}",
		Description: "desc",
	}

	mSecret := &specV1.Secret{
		Namespace: "namespace",
		Name:      "testSecret",
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Data: map[string][]byte{
			"cfg1": []byte("123"),
			"cfg2": []byte("abc"),
		},
		Description: "desc",
	}

	res, err := ToSecretModel(secret)
	assert.NoError(t, err)
	assert.Equal(t, mSecret.Name, res.Name)
	assert.Equal(t, mSecret.Namespace, res.Namespace)
	assert.Equal(t, mSecret.Labels, res.Labels)
	assert.Equal(t, mSecret.Description, res.Description)
	assert.Equal(t, mSecret.Data, res.Data)

	secret = &Secret{
		ID:        123,
		Name:      "testSecret",
		Namespace: "namespace",
	}

	_, err = ToSecretModel(secret)
	assert.NotNil(t, err)
}
