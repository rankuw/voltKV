package store

import (
	"fmt"
	"sync"
)

type Store struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

func (store *Store) Set(key, value string) {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.data[key] = value

	fmt.Println(store.data)
}

func (store *Store) Get(key string) (val string, ok bool) {
	store.mu.RLock()
	defer store.mu.RLocker().Unlock()

	val, ok = store.data[key]
	return
}
