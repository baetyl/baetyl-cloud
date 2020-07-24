package service

import (
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func genPropertyTestCase() *models.Property {
	property := &models.Property{
		Key:   "baetyl_0.1.0",
		Value: "http://test/0.1.0",
	}
	return property
}

func TestCreateProperty(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mConf := genPropertyTestCase()

	mockObject.property.EXPECT().CreateProperty(mConf).Return(nil).Times(1)

	cs, err := NewPropertyService(mockObject.conf)
	assert.NoError(t, err)
	err = cs.CreateProperty(mConf)
	assert.NoError(t, err)
}

func TestDeleteProperty(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	mConf := genPropertyTestCase()
	mockObject.property.EXPECT().DeleteProperty(mConf.Key).Return(nil).Times(1)

	cs, err := NewPropertyService(mockObject.conf)
	assert.NoError(t, err)
	err = cs.DeleteProperty(mConf.Key)
	assert.NoError(t, err)
}

func TestGetProperty(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mConf := genPropertyTestCase()

	mockObject.property.EXPECT().GetProperty(mConf.Key).Return(mConf, nil).Times(1)
	cs, err := NewPropertyService(mockObject.conf)
	assert.NoError(t, err)
	res, err := cs.GetProperty(mConf.Key)
	assert.NoError(t, err)
	checkProperty(t, res, mConf)
}

func TestListProperty(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mConf := []models.Property{
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
	mockObject.property.EXPECT().ListProperty(page).Return(mConf, 1, nil).Times(1)

	cs, err := NewPropertyService(mockObject.conf)
	assert.NoError(t, err)
	res, count, err := cs.ListProperty(page)
	assert.NoError(t, err)
	assert.Equal(t, count, 1)
	checkProperty(t, &mConf[0], &res[0])
}

func TestUpdateProperty(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	mConf := genPropertyTestCase()

	mockObject.property.EXPECT().UpdateProperty(mConf).Return(nil).Times(1)

	cs, err := NewPropertyService(mockObject.conf)
	assert.NoError(t, err)
	err = cs.UpdateProperty(mConf)
	assert.NoError(t, err)
}

func checkProperty(t *testing.T, expect, actual *models.Property) {
	assert.Equal(t, expect.Key, actual.Key)
	assert.Equal(t, expect.Value, actual.Value)
}
