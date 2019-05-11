package store

import (
	"errors"
	"sync"

	"github.com/askft/kademlia/encoding"
)

// MemStore is an in-memory (volatile) store for DHT data.
type MemStore struct {
	sync.Mutex
	m map[string][]byte
}

// NewMemStore creates and returns a new MemStore handle.
func NewMemStore() *MemStore {
	return &MemStore{m: make(map[string][]byte)}
}

// Put stores `data` in volatile memory and returns its key.
func (s *MemStore) Put(data []byte) (string, error) {
	s.Lock()
	defer s.Unlock()
	key := encoding.EncodeData(data)
	s.m[key] = data[:]
	return key, nil
}

// Get returns the data at `key` if it exists, where
// `key` is a base64-encoded SHA-1 hash of some data.
func (s *MemStore) Get(key string) ([]byte, error) {
	if data, ok := s.m[key]; ok {
		return data[:], nil
	}
	return nil, errors.New("invalid key")
}

// Delete removes the data at `key` if it exists, where
// `key` is a base64-encoded SHA-1 hash of some data.
func (s *MemStore) Delete(key string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.m, key)
	return nil
}
