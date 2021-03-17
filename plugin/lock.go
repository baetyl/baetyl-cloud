package plugin

//go:generate mockgen -destination=../mock/plugin/lock.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Locker

// Locker - the lock manager for baetyl cloud
type Locker interface {

	// Lock lock the resource, Lock should be paired with Unlock.
	// PARAMS:
	//   - name: the lock's name
	// RETURNS:
	//   error: if has error else nil
	Lock(name, value string) error

	// LockWithExpireTime lock the resource with expire time
	// PARAMS:
	//   - name: the lock's name
	//   - expireTime(seconds): the expire time of the lock, if acquired the lock
	// RETURNS:
	//   error: if has error else nil
	LockWithExpireTime(name, value string, expireTime int64) error

	// Unlock release the lock by name
	// PARAMS:
	//	 - name: the lock's name
	// RETURNS:
	//   error: if has error else nil
	Unlock(name, value string) error
}
