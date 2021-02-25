package lock

import (
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

func init() {
	plugin.RegisterFactory("defaultlocker", New)
}

type taskLocker struct {
	task plugin.Task
}

var _ plugin.Locker = &taskLocker{}

const DefaultExpireTime = 100

func New() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, err
	}

	task, err := plugin.GetPlugin(cfg.DefaultLocker.Storage)

	if err != nil {
		return nil, err
	}

	return &taskLocker{
		task: task.(plugin.Task),
	}, nil

}

// Lock lock the resource DefaultExpireTime seconds
func (l *taskLocker) Lock(name string) (bool, error) {
	return l.LockWithExpireTime(name, DefaultExpireTime)
}

// LockWithExpireTime lock the resource expireTime seconds
func (l *taskLocker) LockWithExpireTime(name string, expireTime int64) (bool, error) {
	t, err := l.task.GetTask(name)
	if err != nil {
		return false, err
	}

	if t == nil {
		return false, nil
	}

	t.ExpireTime = expireTime
	return l.task.AcquireTaskLock(t)
}

// ReleaseLook releaseLock
func (l *taskLocker) ReleaseLock(name string) (bool, error) {
	return l.LockWithExpireTime(name, 0)
}

func (l *taskLocker) Close() error {
	return l.task.Close()
}
