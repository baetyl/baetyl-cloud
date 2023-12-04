// Package entities 数据库存储基本结构与方法
package entities

import (
	"testing"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

func TestFromConfigModel(t *testing.T) {
	cfg := &Configuration{
		Name:        "testCfg",
		Namespace:   "namespace",
		Labels:      `{"baetyl-cloud-system":"true"}`,
		Data:        `{"cfg1":"123","cfg2":"abc"}`,
		Description: "desc",
	}

	mCfg := &specV1.Configuration{
		Namespace: "namespace",
		Name:      "testCfg",
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Data: map[string]string{
			"cfg1": "123",
			"cfg2": "abc",
		},
		Description: "desc",
	}
	res, err := FromConfigModel("namespace", mCfg)

	assert.NoError(t, err)
	assert.Equal(t, cfg.Name, res.Name)
	assert.Equal(t, cfg.Namespace, res.Namespace)
	assert.Equal(t, cfg.Labels, res.Labels)
	assert.Equal(t, cfg.Description, res.Description)
	assert.Equal(t, cfg.Data, res.Data)

	//// size too large
	//largeCfg := &specV1.Configuration{
	//	Namespace: "namespace",
	//	Name:      "largeCfg",
	//	Labels: map[string]string{
	//		common.LabelSystem: "true",
	//	},
	//	Data:        map[string]string{},
	//	Description: "desc",
	//}
	//
	//data := ""
	//for i := 0; i < 1024*1024; i++ {
	//	data += "1"
	//}
	//largeCfg.Data["data"] = data
	//_, err = FromConfigModel("namespace", largeCfg)
	//assert.Error(t, err)
}

func TestToConfigModel(t *testing.T) {
	cfg := &Configuration{
		ID:          123,
		Name:        "testCfg",
		Namespace:   "namespace",
		Labels:      "{\"baetyl-cloud-system\":\"true\"}",
		Data:        "{\"cfg1\":\"123\",\"cfg2\":\"abc\"}",
		Description: "desc",
	}

	mCfg := &specV1.Configuration{
		Namespace: "namespace",
		Name:      "testCfg",
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Data: map[string]string{
			"cfg1": "123",
			"cfg2": "abc",
		},
		Description: "desc",
	}

	res, err := ToConfigModel(cfg)
	assert.NoError(t, err)
	assert.Equal(t, mCfg.Name, res.Name)
	assert.Equal(t, mCfg.Namespace, res.Namespace)
	assert.Equal(t, mCfg.Labels, res.Labels)
	assert.Equal(t, mCfg.Description, res.Description)
	assert.Equal(t, mCfg.Data, res.Data)

	cfg = &Configuration{
		ID:        123,
		Name:      "testCfg",
		Namespace: "namespace",
	}

	_, err = ToConfigModel(cfg)
	assert.NotNil(t, err)
}

func TestEqualConfig(t *testing.T) {
	cfg1 := &specV1.Configuration{
		Namespace: "namespace",
		Name:      "testCfg",
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Data: map[string]string{
			"cfg1": "123",
			"cfg2": "abc",
		},
		Description: "desc",
	}
	cfg2 := &specV1.Configuration{
		Namespace: "namespace",
		Name:      "testCfg",
		Labels: map[string]string{
			common.LabelSystem: "true",
		},
		Data: map[string]string{
			"cfg1": "123",
			"cfg2": "abc",
		},
		Description: "desc",
	}
	res := EqualConfig(cfg1, cfg2)
	assert.True(t, res)

	cfg2.Namespace = "namespace-2"
	res = EqualConfig(cfg1, cfg2)
	assert.False(t, res)
	cfg2.Namespace = cfg1.Namespace

	cfg2.Name = "testCfg-2"
	res = EqualConfig(cfg1, cfg2)
	assert.False(t, res)
	cfg2.Name = cfg1.Name

	cfg2.Description = "desc-2"
	res = EqualConfig(cfg1, cfg2)
	assert.False(t, res)
	cfg2.Description = cfg1.Description

	cfg2.Labels = map[string]string{}
	res = EqualConfig(cfg1, cfg2)
	assert.False(t, res)
	cfg2.Labels = cfg1.Labels

	cfg2.Data = map[string]string{}
	res = EqualConfig(cfg1, cfg2)
	assert.False(t, res)
	cfg2.Data = cfg1.Data
}
