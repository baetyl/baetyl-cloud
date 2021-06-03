package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	_ "github.com/baetyl/baetyl-cloud/v2/plugin/default/auth"
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
