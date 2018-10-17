package encoding

import (
	"crypto/sha1"
	"encoding/base64"
)

// Size is the length in bytes of a key.
const Size = sha1.Size

// EncodeData returns the base64-encoded SHA1-hash of `data`.
// (Data -> SHA-1 -> Base64).
func EncodeData(data []byte) string {
	return EncodeHash(HashData(data))
}

// HashData returns the SHA-1 hash of `data`.
// (Data -> SHA-1).
func HashData(data []byte) [sha1.Size]byte {
	return sha1.Sum(data)
}

// EncodeHash encodes a SHA-1 hash into a base64 string.
// (SHA-1 -> Base64).
func EncodeHash(hash [sha1.Size]byte) string {
	return base64.StdEncoding.EncodeToString(hash[:])
}

// DecodeKeyStr decodes a byte array encoded as a base64 string.
// (Base64 -> SHA-1).
func DecodeKeyStr(key string) ([sha1.Size]byte, error) {
	hash := [sha1.Size]byte{}
	dec, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return hash, err
	}
	for i := 0; i < sha1.Size; i++ {
		hash[i] = dec[i]
	}
	return hash, nil
}
