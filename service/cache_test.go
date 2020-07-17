package service

import (
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func genSystemConfigTestCase() *models.SystemConfig{
	systemConfig := &models.SystemConfig{
		Key:   "baetyl_0.1.0",
		Value: "http://test/0.1.0",
	}
	return systemConfig
}
func TestSystemConfigService_GetSystemConfig(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	mConf := genSystemConfigTestCase()

	mockObject.cacheStorage.EXPECT().GetSystemConfig(mConf.Key).Return(mConf, nil)

	cs, err := NewCacheService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.GetSystemConfig(mConf.Key)
	assert.NoError(t, err)
	assert.Equal(t, mConf.Value, res.Value)
}

func TestSystemConfigService_ListSystemConfig(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mConf := genSystemConfigTestCase()
	page := &models.Filter{
		PageNo:   1,
		PageSize: 10,
		Name:     "%",
	}
	mockObject.cacheStorage.EXPECT().CountSystemConfig(gomock.Any()).Return(1, nil)
	mockObject.cacheStorage.EXPECT().ListSystemConfig(page.Name, page.PageNo, page.PageSize).Return([]models.SystemConfig{*mConf}, nil)

	cs, err := NewCacheService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.ListSystemConfig(page)
	assert.NoError(t, err)
	assert.EqualValues(t, *mConf, res.Items.([]models.SystemConfig)[0])

}

func TestSystemConfigService_CreateSystemConfig(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	mConf := genSystemConfigTestCase()

	mockObject.cacheStorage.EXPECT().CreateSystemConfig(mConf).Return(nil, nil).AnyTimes()
	mockObject.cacheStorage.EXPECT().GetSystemConfig(mConf.Key).Return(mConf, nil).AnyTimes()

	cs, err := NewCacheService(mockObject.conf)
	assert.NoError(t, err)
	_, err = cs.CreateSystemConfig(mConf)
	assert.NoError(t, err)
}

func TestSystemConfigService_UpdateSystemConfig(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	mConf := genSystemConfigTestCase()

	mockObject.cacheStorage.EXPECT().UpdateSystemConfig(mConf).Return(nil, nil)
	mockObject.cacheStorage.EXPECT().GetSystemConfig(mConf.Key).Return(mConf, nil).AnyTimes()

	cs, err := NewCacheService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.UpdateSystemConfig(mConf)
	assert.NoError(t, err)
	assert.Equal(t, mConf.Value, res.Value)
}

func TestSystemConfigService_DeleteSystemConfig(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	mConf := genSystemConfigTestCase()

	mockObject.cacheStorage.EXPECT().DeleteSystemConfig(mConf.Key).Return(nil, nil)

	cs, err := NewCacheService(mockObject.conf)
	assert.NoError(t, err)
	err = cs.DeleteSystemConfig(mConf.Key)
	assert.NoError(t, err)
}
