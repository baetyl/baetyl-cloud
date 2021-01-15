package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

func TestPropertyService_GetProperty(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	property := mockObject.property
	is := propertyService{}
	is.property = mockObject.property

	name := "a"
	m := &models.Property{
		Name:  "a",
		Value: "a-value",
	}
	property.EXPECT().GetProperty(name).Return(m, nil).Times(1)
	res, err := is.GetProperty(name)
	assert.NoError(t, err)
	assert.Equal(t, res.Name, m.Name)
	assert.Equal(t, res.Value, m.Value)
}
