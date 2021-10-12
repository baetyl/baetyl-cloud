package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskService(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	_, err := NewTaskService(mockObject.conf)
	assert.NoError(t, err)
}
