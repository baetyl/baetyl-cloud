package service

import (
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/common"
	_ "github.com/baetyl/baetyl-cloud/v2/plugin/default/auth"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_Authenticate(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	c := &common.Context{}
	mockObject.auth.EXPECT().Authenticate(c).Return(nil).Times(1)

	as, err := NewAuthService(mockObject.conf)
	assert.Nil(t, err)
	err = as.Authenticate(c)
	assert.Nil(t, err)
}

func TestAuthService_Sign(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	data := []byte("test")
	mockObject.auth.EXPECT().SignToken(data).Return(data, nil).Times(1)

	as, err := NewAuthService(mockObject.conf)
	assert.Nil(t, err)
	res, err := as.SignToken(data)
	assert.Nil(t, err)
	assert.Equal(t, data, res)
}

func TestAuthService_Verify(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	data := []byte("test")
	mockObject.auth.EXPECT().VerifyToken(data, data).Return(true).Times(1)

	as, err := NewAuthService(mockObject.conf)
	assert.Nil(t, err)
	res := as.VerifyToken(data, data)
	assert.Equal(t, true, res)
}

func TestAuthService_GenToken(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	data := map[string]interface{}{
		"123": "123",
	}
	mockObject.auth.EXPECT().SignToken([]byte("{\"123\":\"123\"}")).Return([]byte("test"), nil).Times(1)

	as, err := NewAuthService(mockObject.conf)
	assert.Nil(t, err)
	res, err := as.GenToken(data)
	assert.Nil(t, err)
	assert.Equal(t, "098f6bcd467b22313233223a22313233227d", res)
}
