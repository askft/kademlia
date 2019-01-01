package peer

import (
	"testing"

	"github.com/askft/kademlia/node"
)

func init() {
}

func TestXorDistancePrefixLength(t *testing.T) {
	a := node.Key{}
	b := node.Key{}

	b[19] = 0 // default
	assertEqual(t, node.Distance(a, b).PrefixLength(), 159)

	b[19] = 1
	assertEqual(t, node.Distance(a, b).PrefixLength(), 159)

	b[19] = 2
	assertEqual(t, node.Distance(a, b).PrefixLength(), 158)

	b[0] = 255
	assertEqual(t, node.Distance(a, b).PrefixLength(), 0)
}

func assertEqual(t *testing.T, value, expected interface{}) {
	if value != expected {
		t.Errorf("Expected %v, got %v.\n", expected, value)
	}
}

func assertNotEqual(t *testing.T, value, expected interface{}) {
	if value == expected {
		t.Errorf("Expected something else than %v, but got %v.\n",
			expected, value)
	}
}
