package service

import (
	"testing"

	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/stretchr/testify/assert"
)

func TestCacheService(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	// good case1: get from db
	mConf := &models.Property{
		Key:   "baetyl_0.1.0",
		Value: "http://test.baetyl/0.1.0",
	}
	mockObject.property.EXPECT().GetProperty(mConf.Key).Return(mConf, nil).AnyTimes()

	cache, err := NewCacheService(mockObject.property)
	assert.NoError(t, err)
	res, err := cache.Get(mConf.Key)
	assert.NoError(t, err)
	assert.Equal(t, res, mConf.Value)
	// good case2: get from cache
	res, err = cache.Get(mConf.Key)
	assert.NoError(t, err)
	assert.Equal(t, res, mConf.Value)

	// bad case
	key := "bad key"
	mockObject.property.EXPECT().GetProperty(key).Return(nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("key", key)))
	assert.NoError(t, err)
	res, err = cache.Get(key)
	assert.Error(t, err)

}
