package service

import (
	"fmt"
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigService_Create(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	namespace := "default"
	name := "ConfigService-Create"
	mConf := &specV1.Configuration{Name: name}
	mockObject.configuration.EXPECT().CreateConfig(namespace, mConf).Return(mConf, nil)
	cs, err := NewConfigService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.Create(namespace, mConf)
	assert.NoError(t, err)
	assert.Equal(t, name, res.Name)
}

func TestDefaultConfigService_Get(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	namespace := "default"
	name := "ConfigService-Get"
	mConf := &specV1.Configuration{
		Name:    name,
		Version: "get.0.0.1",
	}

	mockObject.configuration.EXPECT().GetConfig(namespace, name, "").Return(mConf, nil)

	cs, err := NewConfigService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.Get(namespace, name, "")
	assert.NoError(t, err)
	assert.Equal(t, mConf.Version, res.Version)
}

func TestDefaultConfigService_List(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	namespace := "default"
	name := "ConfigService-List"

	selector := &models.ListOptions{
		LabelSelector: "name=testlist",
	}

	mConf := specV1.Configuration{
		Name:    name,
		Version: "get.0.0.1",
	}

	configList := &models.ConfigurationList{
		Total:       1,
		ListOptions: selector,
		Items:       []specV1.Configuration{mConf},
	}

	mockObject.configuration.EXPECT().ListConfig(namespace, selector).Return(configList, nil)

	cs, err := NewConfigService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.List(namespace, selector)
	assert.NoError(t, err)
	assert.Equal(t, selector, res.ListOptions)
}

func TestDefaultConfigService_Update(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	cs := configService{
		config: mockObject.configuration,
	}

	namespace := "default"
	name := "config"
	mConf := &specV1.Configuration{
		Name:    name,
		Version: "1243",
	}

	mockObject.configuration.EXPECT().UpdateConfig(namespace, mConf).Return(nil, fmt.Errorf("error"))
	_, err := cs.Update(namespace, mConf)
	assert.NotNil(t, err)

	mockObject.configuration.EXPECT().UpdateConfig(namespace, mConf).Return(mConf, nil).AnyTimes()
	_, err = cs.Update(namespace, mConf)
	assert.NoError(t, err)
}

func TestDefaultConfigService_Upsert(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	cs := configService{
		config: mockObject.configuration,
	}

	namespace := "default"
	name := "config"
	mConf := &specV1.Configuration{
		Name:    name,
		Version: "1243",
	}

	mockObject.configuration.EXPECT().GetConfig(namespace, mConf.Name, "").Return(nil, fmt.Errorf("error"))
	mockObject.configuration.EXPECT().CreateConfig(namespace, mConf).Return(mConf, nil)
	_, err := cs.Upsert(namespace, mConf)
	assert.NoError(t, err)

	mockObject.configuration.EXPECT().GetConfig(namespace, mConf.Name, "").Return(mConf, nil)
	_, err = cs.Upsert(namespace, mConf)
	assert.NoError(t, err)
}

func TestDefaultConfigService_Delete(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	namespace := "default"
	name := "ConfigService-update"

	mockObject.configuration.EXPECT().DeleteConfig(namespace, name).Return(nil)
	mockObject.index.EXPECT().ListIndex(namespace, common.Application, common.Config, name).Return([]string{}, nil).AnyTimes()

	cs, err := NewConfigService(mockObject.conf)
	assert.NoError(t, err)
	err = cs.Delete(namespace, name)
	assert.NoError(t, err)
}
