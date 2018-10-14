package main

import (
	"fmt"
	"strings"
)

/*
	Useful operations on the bit level.
*/

// FindSetBits returns a list of indices for the bits of `id` that equal 1.
// Example: `00100101 --> [5, 2, 0]`.
func FindSetBits(id NodeID) []int {
	idSet := []int{}
	for i := 0; i < len(id); i++ {
		byteSet := []int{}
		for j := 7; j >= 0; j-- {
			if isSet(id[i], j) {
				byteSet = append(byteSet, (len(id)-i-1)*8+j)
			}
		}
		idSet = append(idSet, byteSet...)
	}
	return idSet
}

// BytePrefixLength returns the index of the MSB starting from index(LSB) == 1.
// Example: 00011111 --> 5.
func BytePrefixLength(x byte) int {
	return 8 - log2(x)
}

func isSet(b byte, i int) bool {
	return (b>>uint(i))&1 != 0
	// return b & (1 << uint(i)) != 0
}

func log2(x byte) int {
	r := 0
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
