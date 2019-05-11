package node

import (
	"math/rand"
	"time"

	"github.com/askft/kademlia/encoding"
)

// Key uniquely identifies a peer in the network
// and is used to address data in the DHT.
type Key [KeySizeBytes]byte

const (
	// KeySizeBytes is the length in bytes of a peer's key.
	KeySizeBytes = encoding.Size

	// KeySizeBits is the length in bits of a peer's key.
	KeySizeBits = KeySizeBytes * 8
)

// GenerateRandomKey creates a randomized node key.
// Uses time.Now().Unix() for the random seed.
func GenerateRandomKey() Key {
	rand.Seed(time.Now().Unix()) // TODO Dangerous?
	key := Key{}
	for i := 0; i < KeySizeBytes; i++ {
		key[i] = uint8(rand.Intn(256))
	}
	return key
}

// String returns a string representation of `key`.
func (key Key) String() string {
	return encoding.EncodeHash(key)
}

// PrefixLength returns the number of leading zeros in a Key.
func (key Key) PrefixLength() int {
	for i := 0; i < KeySizeBytes; i++ {
		if key[i] == 0 {
			continue
		}
		return i*8 + BytePrefixLength(key[i])
	}
	return KeySizeBytes*8 - 1
}

// Distance returns the XOR distance between two keys.
func (lhs Key) Distance(rhs Key) Key {
	key := Key{}
	for i := 0; i < KeySizeBytes; i++ {
		key[i] = lhs[i] ^ rhs[i]
	}
	return key
}

// Less returns true if lhs < rhs when both are interpreted as a number.
func (lhs Key) Less(rhs Key) bool {
	for i := 0; i < KeySizeBytes; i++ {
		if lhs[i] != rhs[i] {
			return lhs[i] < rhs[i]
		}
	}
	return false
}

// Equal returns true if all bytes of lhs and rhs are equal.
func (lhs Key) Equal(rhs Key) bool {
	for i := 0; i < KeySizeBytes; i++ {
		if lhs[i] != rhs[i] {
			return false
		}
	}
	return true
}

// decodeHexKey converts a hexadecimal string representation
// of a node ID into a Key representation (a list of bytes).
// TODO_UNUSED
// func decodeHexKey(data string) (Key, error) {
// 	key := Key{}
// 	dec, err := hex.DecodeString(data)
// 	if err != nil {
// 		return key, err
// 	}
// 	for i := 0; i < KeySizeBytes; i++ {
// 		key[i] = dec[i]
// 	}
// 	return key, nil
// }
