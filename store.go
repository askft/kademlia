package main

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"sync"
)

// Store is the interface for a peer's DHT data storage mechanism.
type Store interface {
	Put(value []byte) (string, error)
	Get(key string) ([]byte, error)
	Delete(key string) error
}

// LocalStore is an in-memory (volatile) store for DHT data.
type LocalStore struct {
	sync.Mutex
	data map[string][]byte
}

// NewLocalStore creates and returns a new LocalStore handle.
func NewLocalStore() *LocalStore {
	return &LocalStore{data: make(map[string][]byte)}
}

// Put stores `value` in volatile memory.
func (s *LocalStore) Put(value []byte) (string, error) {
	s.Lock()
	defer s.Unlock()
	key := encodeValue(value)
	s.data[key] = value[:]
	return key, nil
}

// Get returns the value at `key` if it exists, where
// `key` is a base64-encoded SHA-1 hash of some data.
func (s *LocalStore) Get(key string) ([]byte, error) {
	if value, ok := s.data[key]; ok {
		return value[:], nil
	}
	return nil, errors.New("invalid key")
}

// Delete removes the value at `key` if it exists, where
// `key` is a base64-encoded SHA-1 hash of some data.
func (s *LocalStore) Delete(key string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.data, key)
	return nil
}

// OptionalTODO use base58 instead of base64

func encodeID(id NodeID) string {
	return base64.StdEncoding.EncodeToString(id[:])
}

func encodeValue(value []byte) string {
	hash := sha1.Sum(value)
	return base64.StdEncoding.EncodeToString(hash[:])
}

// TODO_UNUSED
func decode(key string) ([]byte, error) {
	hash, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
