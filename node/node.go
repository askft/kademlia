package node

import (
	"fmt"
	"math/rand"
	"net"
	"sort"
	"time"

	"github.com/askft/kademlia/encoding"
)

/*
	Node operations:
		- generate random key
		- encode/decode to/from hexstring/bytes
		- prefix length (number of leading zeros)
		- xor distance
		- less?
		- equal?
		- sort list by distance
*/

const (
	// KeyByteLen is the length in bytes of a peer's key.
	KeyByteLen = encoding.Size

	// KeyBitLen is the length in bits of a peer's key.
	KeyBitLen = KeyByteLen * 8
)

// Key uniquely identifies a peer in the network
// and is used to address data in the DHT.
type Key [KeyByteLen]byte

// Contact is primarily used to group node key, host and port,
// but also contains some extra optional useful data.
type Contact struct {
	Key  Key
	Host net.IP
	Port string
	RTT  int
}

func (contact Contact) String() string {
	return fmt.Sprintf("%s, [ %s ]", contact.Address(), contact.Key)
}

// Address formats `contact` as a `host:port` string.
func (contact Contact) Address() string {
	return contact.Host.String() + ":" + contact.Port
}

// GenerateRandomKey creates a randomized node key.
// Uses time.Now().Unix() for the random seed.
func GenerateRandomKey() Key {
	rand.Seed(time.Now().Unix())
	key := Key{}
	for i := 0; i < KeyByteLen; i++ {
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
	for i := 0; i < KeyByteLen; i++ {
		if key[i] == 0 {
			continue
		}
		return i*8 + BytePrefixLength(key[i])
	}
	return KeyByteLen*8 - 1
}

// Distance returns the XOR distance between two `Key`s.
func Distance(lhs, rhs Key) Key {
	key := Key{}
	for i := 0; i < KeyByteLen; i++ {
		key[i] = lhs[i] ^ rhs[i]
	}
	return key
}

// Less returns `true` if `lhs` < `rhs` when both are interpreted as a number.
func Less(lhs, rhs Key) bool {
	for i := 0; i < KeyByteLen; i++ {
		if lhs[i] != rhs[i] {
			return lhs[i] < rhs[i]
		}
	}
	return false
}

// Equal returns `true` if all bytes of `lhs` and `rhs` are equal.
// 	TODO_UNUSED
func Equal(lhs, rhs Key) bool {
	for i := 0; i < KeyByteLen; i++ {
		if lhs[i] != rhs[i] {
			return false
		}
	}
	return true
}

// SortByDistance sorts `list` by distance to `key`.
func SortByDistance(list []Contact, key Key) {
	sort.SliceStable(list, func(i, j int) bool {
		d1 := Distance(key, list[i].Key)
		d2 := Distance(key, list[j].Key)
		return Less(d1, d2)
	})
}

// -------------------------------------------------------------

// func (key Key) String() string {
// 	return hex.EncodeToString(key[0:KeyByteLen])
// }

// decodeHexKey converts a hexadecimal string representation
// of a node ID into a Key representation (a list of bytes).
//  TODO_UNUSED
// func decodeHexKey(data string) (Key, error) {
// 	key := Key{}
// 	dec, err := hex.DecodeString(data)
// 	if err != nil {
// 		return key, err
// 	}
// 	for i := 0; i < KeyByteLen; i++ {
// 		key[i] = dec[i]
// 	}
// 	return key, nil
// }
