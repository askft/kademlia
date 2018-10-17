package peer

import (
	"fmt"
	"sync"
	"time"

	"p2p/node"
	"p2p/store"
)

// Bucket is a list of contacts. Note that a bucket
// should maximally hold `k` elements.
type Bucket []node.Contact

// Peer keeps track of relevant state for the Kademlia network.
// TODO embed Options instead of having all this shit
type Peer struct {
	node.Contact
	store        store.Store
	networkID    string                    // Prevents networks merging together.
	routingTable [node.KeyBitLen]Bucket    // Every bucket corresponds to a specific distance.
	refreshMap   [node.KeyBitLen]time.Time // TODO Look closer into when/where to refresh.
	mutex        sync.Mutex                // TODO Use RWMutex instead? And check carefully where this might be needed.
}

// NewPeer initializes a peer and returns a handle to it.
func NewPeer(options *Options) (*Peer, error) {
	peer := &Peer{}
	peer.Key = options.Key
	peer.Host = options.Host
	peer.Port = options.Port
	peer.store = options.Store
	peer.networkID = options.NetworkID
	peer.routingTable = [node.KeyBitLen]Bucket{}
	peer.refreshMap = [node.KeyBitLen]time.Time{}
	return peer, nil
}

// Bootstrap lets `peer` join a network using a predefined set of nodes.
//  See http://xlattice.sourceforge.net/components/protocol/kademlia/specs.html#join
func (peer *Peer) Bootstrap(bootstrapContact node.Contact) {

	// Add the bootstrap node into this peer's appropriate bucket.
	peer.BucketUpdate(bootstrapContact)

	// Perform a self-lookup against the bootstrap node. This populates other
	// peers' k-buckets with this peer, and populates this peer's k-buckets
	// with peers known by the bootstrap node.
	done := make(chan MessageResponseFindNode)
	peer.SendFindNode(bootstrapContact, peer.Key, done)
	res := <-done

	fmt.Println("contacts from IterativeFN")
	printContacts(res.Contacts)

	// Populate this peer's buckets with found contacts.
	for _, contact := range res.Contacts {
		q := peer.bucketIndex(contact.Key)
		peer.RefreshBucket(q)
		peer.BucketUpdate(contact) // TODO Recently added this - correct?
	}
}

func printContacts(contacts []node.Contact) {
	for _, contact := range contacts {
		fmt.Println(" -", contact)
	}
}

func (peer *Peer) PrintAllContacts() {
	for _, bucket := range peer.routingTable {
		for _, contact := range bucket {
			fmt.Println(" -", contact)
		}
	}
}

// RefreshBucket resets the last refresh time for bucket number `q`.
func (peer *Peer) RefreshBucket(q int) {
	peer.refreshMap[q] = time.Now()
}

// FindClosest finds the `n` closest contacts to `target` in
// the peer's routing table.
func (peer *Peer) FindClosest(target node.Key, n int) []node.Contact {
	d := node.Distance(peer.Key, target)
	closest := []node.Contact{}
	seq := NewIntSet()
	seq.AddMany(node.FindSetBits(d[:]))

	// Descend through 1-bits in `d` toward 0 and try to fill `closest`.
	for _, q := range seq.SortedReverse() {
		bucket := peer.routingTable[q]
		if tryFill(&closest, bucket, n) {
			fmt.Println("Filled up `closest` at bucket", q)
			break
		}
	}

	// If `closest` is still not filled, search unvisisted buckets [0, 160).
	for q := 0; q < node.KeyBitLen; q++ {
		if !seq.Has(q) {
			bucket := peer.routingTable[q]
			if tryFill(&closest, bucket, n) {
				fmt.Println("Filled up `closest` at bucket", q)
				break
			}
		}
	}
	return closest
}

func tryFill(closest *[]node.Contact, bucket Bucket, n int) bool {
	for _, contact := range bucket {
		*closest = append(*closest, contact)
		if len(*closest) == n {
			return true
		}
	}
	return false
}

// BucketUpdate adds `contact` into `peer`'s appropriate bucket if necessary.
func (peer *Peer) BucketUpdate(contact node.Contact) {
	bucket := peer.bucketFor(contact.Key)

	// If the contact already exists, move it to the end of the bucket.
	for i, c := range *bucket {
		if c.Key == contact.Key {
			bucket.moveToTail(i)
			fmt.Printf("Updated %s, bucket %d with contact %s (tail move).\n",
				peer.Address(), peer.bucketIndex(contact.Key), contact.Address())
			return
		}
	}

	// If the bucket has space, add the new contact to the bucket.
	if len(*bucket) < k {
		bucket.addToTail(contact)
		fmt.Printf("Updated %s, bucket %d with contact %s (tail add).\n",
			peer.Address(), peer.bucketIndex(contact.Key), contact.Address())
		return
	}

	// If the bucket is full, ping its head and replace it iff
	// it did not respond within a reasonable time.
	pingChan := make(chan MessageResponsePing)
	peer.SendPing(contact, (*bucket)[0].Key, pingChan) // TODO just send Contact, see client.go
	select {
	case res := <-pingChan:
		fmt.Println("got ping back:", res)
		bucket.moveToTail(0)
	case <-time.After(updateTimeout * time.Millisecond):
		fmt.Println("ping timed out")
		(*bucket)[0] = contact // Replace first item...
		bucket.moveToTail(0)   // ... and move it to the tail.
	}
	fmt.Printf("Updated %s, bucket %d with contact %s.\n",
		peer.Address(), peer.bucketIndex(contact.Key), contact.Address())
}

func (peer *Peer) bucketFor(key node.Key) *Bucket {
	fmt.Printf("%s's bucket index for %s: %d\n",
		peer.Key, key, peer.bucketIndex(key))
	return &peer.routingTable[peer.bucketIndex(key)]
}

func (peer *Peer) bucketIndex(key node.Key) int {
	return node.KeyBitLen - 1 - node.Distance(peer.Key, key).PrefixLength()
}

func (bucket *Bucket) moveToTail(i int) {
	*bucket = append((*bucket)[:i], append((*bucket)[i+1:], (*bucket)[i])...)
}

func (bucket *Bucket) addToTail(contact node.Contact) {
	*bucket = append(*bucket, contact)
}

// Store operations ----------------------------------------------------------

// Put stores `value` in `peer`'s storage.
func (peer *Peer) Put(value []byte) (string, error) {
	return peer.store.Put(value)
}

// Get returns the value at `key` in `peer`'s storage if it exists.
func (peer *Peer) Get(key string) ([]byte, error) {
	return peer.store.Get(key)
}

// Delete removes the value at `key` from `peer`'s storage.
func (peer *Peer) Delete(key string) error {
	return peer.store.Delete(key)
}
