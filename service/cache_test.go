package service

import (
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func genCacheTestCase() *models.Cache{
	cache := &models.Cache{
		Key:   "baetyl_0.1.0",
		Value: "http://test/0.1.0",
	}
	return cache
}
func TestGet(t *testing.T){
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mConf := genCacheTestCase()

	mockObject.cacheStorage.EXPECT().GetCache(mConf.Key).Return(mConf.Value, nil).Times(1)
	cs, err := NewCacheService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.Get(mConf.Key)
	assert.NoError(t, err)
	assert.Equal(t, res, mConf.Value)
}
func TestSet(t *testing.T){
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mConf := genCacheTestCase()

	mockObject.cacheStorage.EXPECT().SetCache(mConf.Key, mConf.Value).Return( nil).Times(1)
	mockObject.cacheStorage.EXPECT().GetCache(mConf.Key).Return( mConf.Value, nil).Times(1)

	cs, err := NewCacheService(mockObject.conf)
	assert.NoError(t, err)
	err = cs.Set(mConf.Key, mConf.Value)
	assert.NoError(t, err)
	value, err := cs.Get(mConf.Key)
	assert.NoError(t, err)
	assert.Equal(t, value, mConf.Value)

}

func TestList(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mConf := []models.Cache{
		{
			Key:   "baetyl_0.1.0",
			Value: "http://test/0.1.0",
		},
	}
	page := &models.Filter{
		PageNo:   1,
		PageSize: 10,
		Name:     "%",
	}
	mockObject.cacheStorage.EXPECT().ListCache(page).Return(
		&models.ListView{
			Total: 1,
			PageNo: page.PageNo,
			PageSize: page.PageSize,
			Items: mConf,
		}, nil).Times(1)

	cs, err := NewCacheService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.List(page)
	assert.NoError(t, err)
	checkCache(t, &mConf[0], &res.Items.([]models.Cache)[0])
}

func TestDelete(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	mConf := genCacheTestCase()
	mockObject.cacheStorage.EXPECT().DeleteCache(mConf.Key).Return(nil).Times(1)

	cs, err := NewCacheService(mockObject.conf)
	assert.NoError(t, err)
	err = cs.Delete(mConf.Key)
	assert.NoError(t, err)
}

func checkCache(t *testing.T, expect, actual *models.Cache) {
	assert.Equal(t, expect.Key, actual.Key)
	assert.Equal(t, expect.Value, actual.Value)
}
