package store

import (
	"sync"
	"time"
)

type Store struct {
	data    map[string]string
	expires map[string]time.Time
	mu      sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		data:    make(map[string]string),
		expires: make(map[string]time.Time),
	}
}

func (store *Store) Set(key, value string, ttl time.Duration) {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.data[key] = value

	if ttl > 0 {
		store.expires[key] = time.Now().Add(ttl)
	} else {
		delete(store.expires, key)
	}

}

func (store *Store) Get(key string) (val string, ok bool) {
	store.mu.RLock()
	defer store.mu.RLocker().Unlock()

	if expTime, ok := store.expires[key]; ok {
		if time.Now().After(expTime) {
			delete(store.expires, key)
			delete(store.data, key)

			return "", false
		}
	}

	val, ok = store.data[key]
	return
}
