package plugin

import (
	"context"
	"io"
)

//go:generate mockgen -destination=../mock/plugin/lock.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Locker

// Locker - the lock manager for baetyl cloud
type Locker interface {

	// Lock lock the resource, Lock should be paired with Unlock.
	// PARAMS:
	//   - name: the lock's name
	//   - ttl: expire time of lock, if 0, use default time.
	// RETURNS:
	//   error: if has error else nil
	Lock(ctx context.Context, name string, ttl int64) (string, error)

	// Unlock release the lock by name
	// PARAMS:
	//	 - name: the lock's name
	// RETURNS:
	//   error: if has error else nil
	Unlock(ctx context.Context, name, version string)
	io.Closer
}
