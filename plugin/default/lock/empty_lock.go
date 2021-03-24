package lock

import (
	"context"

	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

func init() {
	plugin.RegisterFactory("defaultlocker", New)
}

type emptyLocker struct{}

func New() (plugin.Plugin, error) {
	return &emptyLocker{}, nil
}

func (l *emptyLocker) Lock(ctx context.Context, name string, ttl int64) (string, error) {
	return "", nil
}

func (l *emptyLocker) Unlock(ctx context.Context, name, version string) {
	return
}

func (l *emptyLocker) Close() error {
	return nil
}
