package mutex

import (
	"path"
	"sync"
)

type Unlock func()

type KeyMutex interface {
	Lock(key string) (Unlock, error)
}

func GenerateKey(subKeys ...string) string {
	return path.Join(subKeys...)
}

var mutexPool = sync.Pool{
	New: func() interface{} {
		return &sync.Mutex{}
	},
}

type lockerPerKey struct {
	lock  *sync.Mutex
	count int
}

type memKeyMutex struct {
	lock     *sync.RWMutex
	keyLocks map[string]*lockerPerKey
}

func NewKeyMutex() KeyMutex {
	return &memKeyMutex{
		lock:     new(sync.RWMutex),
		keyLocks: map[string]*lockerPerKey{},
	}
}

func (m *memKeyMutex) Lock(key string) (Unlock, error) {
	m.lock.Lock()
	locker := m.keyLocks[key]
	if locker == nil {
		lock, _ := mutexPool.Get().(*sync.Mutex)
		locker = &lockerPerKey{
			lock:  lock,
			count: 0,
		}
		m.keyLocks[key] = locker
	}
	locker.count++
	m.lock.Unlock()

	locker.lock.Lock()
	return func() {
		m.unlock(key)
	}, nil
}

func (m *memKeyMutex) unlock(key string) {
	m.lock.Lock()
	locker := m.keyLocks[key]
	if locker == nil {
		m.lock.Unlock()
		return
	}
	locker.lock.Unlock()
	locker.count--
	if locker.count == 0 {
		delete(m.keyLocks, key)
		mutexPool.Put(locker.lock)
	}
	m.lock.Unlock()
}
