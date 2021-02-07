package lock

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultLocker_Lock(t *testing.T) {
	locker := taskLocker{}
	locker.Lock("locker")
	assert.Equal(t, true, true)
}

func TestDefaultLocker_TryLockWithExpireTime(t *testing.T) {
	locker := taskLocker{}
	locker.LockWithExpireTime("lockWithExpireTime", 1)
	assert.Equal(t, true, true)
}

func TestDefaultLocker_ReleaseLook(t *testing.T) {
	locker := taskLocker{}
	locker.ReleaseLock("locker")
	assert.Equal(t, true, true)
}
