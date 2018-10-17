package store

import ()

// Store is the interface for a peer's DHT data storage mechanism.
type Store interface {
	Put(data []byte) (string, error)
	Get(key string) ([]byte, error)
	Delete(key string) error
}
