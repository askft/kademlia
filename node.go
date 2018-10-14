package main

import (
	"encoding/hex"
	"math/rand"
	"net"
)

/*
	Node operations:
		- generate random id
		- encode/decode to/from hexstring/bytes
		- prefix length (number of leading zeros)
		- xor distance
		- less?
		- equal?

	TODO
		- Use len(id) instead of idLength?
		- Have config for the bootstrap node.
*/

type (
	// NodeID uniquely identifies a node in the network.
	NodeID [idLength]byte

	// Key is used to address data in the DHT.
	Key NodeID // TODO unused
)

// Contact is primarily used to group node ID, host and port,
// but also contains some extra optional useful data.
type Contact struct {
	ID   NodeID
	Host net.IP
	Port string
	RTT  int
}

// Address formats `contact` as a `host:port` string.
func (contact Contact) Address() string {
	return contact.Host.String() + ":" + contact.Port
}

// GenerateRandomNodeID creates a randomized node ID.
func GenerateRandomNodeID() NodeID {
	id := NodeID{}
	for i := 0; i < idLength; i++ {
		id[i] = uint8(rand.Intn(256))
	}
	return id
}

// decodeHexNodeID converts a hexadecimal string representation
// of a node ID into a NodeID representation (a list of bytes).
func decodeHexNodeID(data string) (NodeID, error) {
	id := NodeID{}
	dec, err := hex.DecodeString(data)
	if err != nil {
		return id, err
	}
	for i := 0; i < idLength; i++ {
		id[i] = dec[i]
	}
	return id, nil
}

// String on NodeID encodes the bytes as a hexadecimal string.
func (id NodeID) String() string {
	return hex.EncodeToString(id[0:idLength])
}

// PrefixLength returns the number of leading zeros in a NodeID.
func (id NodeID) PrefixLength() int {
	for i := 0; i < idLength; i++ {
		if id[i] == 0 {
			continue
		}
		return i*8 + BytePrefixLength(id[i])
	}
	return idLength * 8 // TODO idLength * 8 - 1?
}

// Distance returns the XOR distance between two `NodeID`s.
func Distance(lhs, rhs NodeID) NodeID {
	id := NodeID{}
	for i := 0; i < idLength; i++ {
		id[i] = lhs[i] ^ rhs[i]
	}
	return id
}

// Less returns `true` if `lhs` < `rhs` when both are interpreted as a number.
func Less(lhs, rhs NodeID) bool {
	for i := 0; i < idLength; i++ {
		if lhs[i] != rhs[i] {
			return lhs[i] < rhs[i]
		}
	}
	return false
}

// Equal returns `true` if all bytes of `lhs` and `rhs` are equal.
// 	TODO_UNUSED
func Equal(lhs, rhs NodeID) bool {
	for i := 0; i < idLength; i++ {
		if lhs[i] != rhs[i] {
			return false
		}
	}
	return true
}
