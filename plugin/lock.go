package plugin

// Locker - the lock manager for baetyl cloud
type Locker interface {

	// Lock lock the resource
	// PARAMS:
	//   -name: the lock's name
	// RETURNs:
	//   true, if locked success
	//   error, if has error
	Lock(name string) (bool, error)

	// LockWithExpireTime lock the resouce with expire time
	// PARAMS:
	//   - name: the lock's name
	//   - expireTime(seconds): the expire time of the lock, if the lock hasn't been released
	// RETURNs:
	//   true, if locked success
	//   error, if has error
	LockWithExpireTime(name string, expireTime int64) (bool, error)

	// ReleaseLock release the lock by name
	// PARAMS:
	//	 - name: the lock's name
	// RETURNs:
	//   true, if lock is released
	//   error, if has error
	ReleaseLock(name string) (bool, error)
}
