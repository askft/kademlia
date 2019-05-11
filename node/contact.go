package node

import (
	"fmt"
	"net"
	"sort"
)

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

// SortByDistance sorts the list of contacts by distance to key.
func SortByDistance(contacts []Contact, key Key) {
	sort.SliceStable(contacts, func(i, j int) bool {
		d1 := key.Distance(contacts[i].Key)
		d2 := key.Distance(contacts[j].Key)
		return d1.Less(d2)
	})
}
