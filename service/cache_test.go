package service

import (
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/stretchr/testify/assert"
)

func TestCacheService(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	// good case
	mConf := &models.Property{
		Key:   "baetyl_0.1.0",
		Value: "http://test.baetyl/0.1.0",
	}
	mockObject.property.EXPECT().GetPropertyValue(mConf.Key).Return(mConf.Value, nil).AnyTimes()

	cache, err := NewCacheService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cache.Get(mConf.Key, mockObject.property.GetPropertyValue)
	assert.NoError(t, err)
	assert.Equal(t, res, mConf.Value)

	// bad case
	key := "bad key"
	mockObject.property.EXPECT().GetPropertyValue(key).Return("", common.Error(
		common.ErrResourceNotFound,
		common.Field("key", key)))
	assert.NoError(t, err)
	res, err = cache.Get(key, mockObject.property.GetPropertyValue)
	assert.Error(t, err)

}
