package lock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyLocker(t *testing.T) {
	locker := &emptyLocker{}

	err := locker.Lock("", "")
	assert.NoError(t, err)

	err = locker.LockWithExpireTime("", "", 10)
	assert.NoError(t, err)

	err = locker.Unlock("", "")
	assert.NoError(t, err)

	err = locker.Close()
	assert.NoError(t, err)
}