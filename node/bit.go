package node

import (
	"fmt"
	"strings"
)

/*
	Useful operations on the bit level.
*/

// FindSetBits returns a list of indices for the bits of `bs` that equal 1.
// Example: `00100101 --> [5, 2, 0]`.
func FindSetBits(bs []byte) []int {
	setBits := []int{}
	for i := 0; i < len(bs); i++ {
		byteSet := []int{}
		for j := 7; j >= 0; j-- {
			if isSet(bs[i], j) {
				byteSet = append(byteSet, (len(bs)-i-1)*8+j)
			}
		}
		setBits = append(setBits, byteSet...)
	}
	return setBits
}

// BytePrefixLength returns the number of bits preceeding the most significant.
// Exampls: 00011111 --> 3. 10101010 --> 0. 00000000 --> 7 [sic!].
func BytePrefixLength(x byte) int {
	if x == 0 {
		return 7
	}
	return 7 - log2(x)
}

func isSet(b byte, i int) bool {
	return (b>>uint(i))&1 != 0
	// return b & (1 << uint(i)) != 0
}

func log2(x byte) int {
	r := -1
	for ; x != 0; x >>= 1 {
		r++
	}
	return r
}

// StringBytes returns a binary string representation of a byte sequence.
func stringBytes(bs []byte) string {
	s := []string{}
	for i := 0; i < len(bs); i++ {
		s = append(s, stringByte(bs[i]))
	}
	return strings.Join(s, " ")
}

func stringByte(b byte) string {
	return fmt.Sprintf("%08b", b)
}
