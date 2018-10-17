package node

import (
	"testing"
)

func init() {

}

func TestFindSetBits(t *testing.T) {
	bs := []byte{5, 2}
	res := FindSetBits(bs)
	t.Log(res)
	// t.Fail()
}

func TestKeyPrefixLength(t *testing.T) {
	key := Key{}
	assertEqual(t, key.PrefixLength(), 159)

	key[19] = 1
	assertEqual(t, key.PrefixLength(), 159)

	key[1] = 1
	assertEqual(t, key.PrefixLength(), 15)

	key[1] = 128
	assertEqual(t, key.PrefixLength(), 8)

	key[0] = 64
	assertEqual(t, key.PrefixLength(), 1)

	key[0] = 128
	assertEqual(t, key.PrefixLength(), 0)
}

func TestBytePrefixLength(t *testing.T) {
	assertEqual(t, BytePrefixLength(0), 7)
	assertEqual(t, BytePrefixLength(1), 7)
	assertEqual(t, BytePrefixLength(2), 6)
	assertEqual(t, BytePrefixLength(4), 5)
	assertEqual(t, BytePrefixLength(8), 4)
	assertEqual(t, BytePrefixLength(16), 3)
	assertEqual(t, BytePrefixLength(32), 2)
	assertEqual(t, BytePrefixLength(64), 1)
	assertEqual(t, BytePrefixLength(128), 0)
	assertEqual(t, BytePrefixLength(255), 0)
}

func TestLog2(t *testing.T) {
	assertEqual(t, log2(255), 7)
	assertEqual(t, log2(128), 7)
	assertEqual(t, log2(127), 6)
	assertEqual(t, log2(35), 5)
	assertEqual(t, log2(35), 5)
	assertEqual(t, log2(10), 3)
	assertEqual(t, log2(4), 2)
	assertEqual(t, log2(3), 1)
	assertEqual(t, log2(2), 1)
	assertEqual(t, log2(1), 0)
	assertEqual(t, log2(0), -1)
}

func assertEqual(t *testing.T, value, expected interface{}) {
	if value != expected {
		t.Errorf("Expected %v, got %v.\n", expected, value)
	}
}
