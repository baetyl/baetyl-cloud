package service

import (
	"fmt"
	"os"
	"testing"

	"github.com/baetyl/baetyl-cloud/common"
	"github.com/stretchr/testify/assert"
)

func TestDefaultCallbackService_Create(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	callback := genCallbackTestCase()
	mockObject.dbStorage.EXPECT().CreateCallback(callback).Return(nil, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().GetCallback(callback.Name, callback.Namespace).Return(callback, nil).AnyTimes()
	cs, err := NewCallbackService(mockObject.conf)
	assert.NoError(t, err)
	_, err = cs.Create(callback)
	assert.NoError(t, err)
}

func TestDefaultCallbackService_Delete(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	callback := genCallbackTestCase()
	res := mockSQLResult{
		lastId: 0,
		affect: 0,
	}
	mockObject.dbStorage.EXPECT().CountBatchByCallback(callback.Name, callback.Namespace).Return(0, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().DeleteCallback(callback.Name, callback.Namespace).Return(&res, nil).AnyTimes()
	cs, err := NewCallbackService(mockObject.conf)
	assert.NoError(t, err)
	err = cs.Delete(callback.Name, callback.Namespace)
	assert.NotNil(t, err)

	res = mockSQLResult{
		lastId: 0,
		affect: 1,
	}
	err = cs.Delete(callback.Name, callback.Namespace)
	assert.NoError(t, err)
}

func TestDefaultCallbackService_Delete_ErrCount(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	callback := genCallbackTestCase()
	mockObject.dbStorage.EXPECT().CountBatchByCallback(callback.Name, callback.Namespace).Return(1, nil).AnyTimes()
	cs, err := NewCallbackService(mockObject.conf)
	assert.NoError(t, err)
	err = cs.Delete(callback.Name, callback.Namespace)
	assert.Error(t, err, common.Error(common.ErrRegisterDeleteCallback, common.Field("name", callback.Name)))
}

func TestDefaultCallbackService_Update(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	callback := genCallbackTestCase()
	mockObject.dbStorage.EXPECT().UpdateCallback(callback).Return(nil, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().GetCallback(callback.Name, callback.Namespace).Return(callback, nil).AnyTimes()
	cs, err := NewCallbackService(mockObject.conf)
	assert.NoError(t, err)
	_, err = cs.Update(callback)
	assert.NoError(t, err)
}

func TestDefaultCallbackService_Update_ErrUpdate(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	callback := genCallbackTestCase()
	mockObject.dbStorage.EXPECT().UpdateCallback(callback).Return(nil, fmt.Errorf("db err update")).AnyTimes()
	cs, err := NewCallbackService(mockObject.conf)
	assert.NoError(t, err)
	_, err = cs.Update(callback)
	assert.NotNil(t, err)
}

func TestDefaultCallbackService_Update_ErrGet(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	callback := genCallbackTestCase()
	mockObject.dbStorage.EXPECT().UpdateCallback(callback).Return(nil, nil).AnyTimes()
	mockObject.dbStorage.EXPECT().GetCallback(callback.Name, callback.Namespace).Return(nil, fmt.Errorf("db err get")).AnyTimes()
	cs, err := NewCallbackService(mockObject.conf)
	assert.NoError(t, err)
	_, err = cs.Update(callback)
	assert.NotNil(t, err)
}

func TestDefaultCallbackService_Get(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	callback := genCallbackTestCase()
	cs, err := NewCallbackService(mockObject.conf)
	assert.NoError(t, err)

	// good case
	mockObject.dbStorage.EXPECT().GetCallback(callback.Name, callback.Namespace).Return(callback, nil).Times(1)
	_, err = cs.Get(callback.Name, callback.Namespace)
	assert.NoError(t, err)

	// bad case 0
	mockObject.dbStorage.EXPECT().GetCallback(callback.Name, callback.Namespace).Return(nil, nil).Times(1)
	_, err = cs.Get(callback.Name, callback.Namespace)
	assert.Error(t, err, common.ErrResourceNotFound, common.Field("type", "callback"), common.Field("name", callback.Name))

	// bad case 1
	mockObject.dbStorage.EXPECT().GetCallback(callback.Name, callback.Namespace).Return(nil, os.ErrInvalid).Times(1)
	_, err = cs.Get(callback.Name, callback.Namespace)
	assert.Error(t, err, common.ErrDatabase, common.Field("error", os.ErrInvalid))
}

func TestDefaultCallbackService_Get_Err(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	callback := genCallbackTestCase()
	mockObject.dbStorage.EXPECT().GetCallback(callback.Name, callback.Namespace).Return(nil, fmt.Errorf("db error")).AnyTimes()
	cs, err := NewCallbackService(mockObject.conf)
	assert.NoError(t, err)
	_, err = cs.Get(callback.Name, callback.Namespace)
	assert.NotNil(t, err)
}

func TestDefaultCallbackService_Callback(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	callback := genCallbackTestCase()
	mockObject.dbStorage.EXPECT().GetCallback(callback.Name, callback.Namespace).Return(callback, nil).AnyTimes()
	cs, err := NewCallbackService(mockObject.conf)
	assert.NoError(t, err)
	data := map[string]string{
		"z": "x",
	}
	_, err = cs.Callback(callback.Name, callback.Namespace, data)
	assert.NotNil(t, err)
}
