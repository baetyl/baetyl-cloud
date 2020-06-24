package service

import (
	"testing"

	"github.com/baetyl/baetyl-cloud/models"
	"github.com/stretchr/testify/assert"
)

func TestSysConfigService_Get(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	mConf := &models.SysConfig{
		Type:  "baetyl",
		Key:   "0.1.0",
		Value: "http://test/0.1.0",
	}

	mockObject.dbStorage.EXPECT().GetSysConfig(mConf.Type, mConf.Key).Return(mConf, nil)

	cs, err := NewSysConfigService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.GetSysConfig(mConf.Type, mConf.Key)
	assert.NoError(t, err)
	assert.Equal(t, mConf.Value, res.Value)
}

func TestSysConfigService_ListAll(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	mConf := []models.SysConfig{
		models.SysConfig{
			Type:  "baetyl",
			Key:   "0.1.0",
			Value: "http://test/0.1.0",
		},
	}

	mockObject.dbStorage.EXPECT().ListSysConfigAll("baetyl").Return(mConf, nil)

	cs, err := NewSysConfigService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.ListSysConfigAll("baetyl")
	assert.NoError(t, err)
	assert.Equal(t, len(mConf), len(res))
}
