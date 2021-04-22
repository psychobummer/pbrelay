package keystore

import (
	"fmt"
	"sync"
)

// MemoryStorage is an in-memory map-based implementation of the
// `Storage` interface. Safe for concurrent access.
type MemoryStorage struct {
	mutex sync.RWMutex
	store map[string]Artifact
}

// NewMemoryStorage returns a pointer to a new MemoryStorage backend
func NewMemoryStorage() *MemoryStorage {
	mockStore := MemoryStorage{
		store: make(map[string]Artifact),
	}
	return &mockStore
}

// Get returns the value stored at `key`
func (s *MemoryStorage) Get(key string) (Artifact, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if v, ok := s.store[key]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("not found")
}

// HasKey returns bool if `key` exists in storage
func (s *MemoryStorage) HasKey(key string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	_, ok := s.store[key]
	return ok
}

// Set sets `key` to `val` in the underlying storage.
func (s *MemoryStorage) Set(key string, val Artifact) error {
	if s.HasKey(key) {
		return fmt.Errorf("key already exists")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.store[key] = val
	return nil
}

// Delete deletes the entry at `key`
func (s *MemoryStorage) Delete(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.store, key)
	return nil
}
