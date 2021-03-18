package lock

import "github.com/baetyl/baetyl-cloud/v2/plugin"

func init() {
	plugin.RegisterFactory("defaultlocker", New)
}

type emptyLocker struct{}

func New() (plugin.Plugin, error) {
	return &emptyLocker{}, nil
}

func (l *emptyLocker) Lock(name, value string) error {
	return nil
}

func (l *emptyLocker) LockWithExpireTime(name, value string, expireTime int64) error {
	return nil
}

func (l *emptyLocker) Unlock(name, value string) error {
	return nil
}

func (l *emptyLocker) Close() error {
	return nil
}
