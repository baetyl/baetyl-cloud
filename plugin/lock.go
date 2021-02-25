package plugin

//go:generate mockgen -destination=../mock/plugin/lock.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Locker

// Locker - the lock manager for baetyl cloud
type Locker interface {

	// Lock lock the resource
	// PARAMS:
	//   - name: the lock's name
	// RETURNS:
	//   true: if locked success
	//   error: if has error else nil
	Lock(name string) (bool, error)

	// LockWithExpireTime lock the resouce with expire time
	// PARAMS:
	//   - name: the lock's name
	//   - expireTime(seconds): the expire time of the lock, if acquired the lock
	// RETURNS:
	//   true: if locked success
	//   error: if has error else nil
	LockWithExpireTime(name string, expireTime int64) (bool, error)

	// ReleaseLock release the lock by name
	// PARAMS:
	//	 - name: the lock's name
	// RETURNS:
	//   true: if lock is released
	//   error: if has error else nil
	ReleaseLock(name string) (bool, error)
}
