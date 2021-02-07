package lock

import (
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type taskLocker struct {
	task plugin.Task
}

var _ plugin.Locker = &taskLocker{}

const DefaultExpireTime = 100

// Lock lock the resource DefaultExpireTime seconds
func (l taskLocker) Lock(name string) (bool, error) {
	return l.LockWithExpireTime(name, DefaultExpireTime)
}

// LockWithExpireTime lock the resource expireTime seconds
func (l taskLocker) LockWithExpireTime(name string, expireTime int64) (bool, error) {
	t, err := l.task.GetTask(name)
	if err != nil {
		return false, err
	}

	t.ExpireTime = expireTime
	return l.task.AcquireTaskLock(t)
}

// ReleaseLook releaseLock
func (l taskLocker) ReleaseLock(name string) (bool, error) {
	return l.LockWithExpireTime(name, 0)
}
