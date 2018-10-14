package main

import (
	"fmt"
	"sync"
	"time"
)

// Bucket is a list of contacts. Note that a bucket
// should maximally hold `k` elements.
type Bucket []Contact

// Peer keeps track of relevant state for the Kademlia network.
// TODO embed Options instead of having all this shit
type Peer struct {
	Contact
	store        Store
	networkID    string       // Prevents networks merging together.
	routingTable [b]Bucket    // Every bucket corresponds to a specific distance.
	refreshMap   [b]time.Time // TODO Look closer into when/where to refresh.
	mutex        sync.Mutex   // TODO Use RWMutex instead? And check carefully where this might be needed.
}

// NewPeer initializes a peer and returns a handle to it.
func NewPeer(options *Options) (*Peer, error) {
	peer := &Peer{}
	peer.ID = options.id
	peer.Host = options.host
	peer.Port = options.port
	peer.store = options.store
	peer.networkID = options.networkID
	peer.routingTable = [b]Bucket{}
	peer.refreshMap = [b]time.Time{}
	return peer, nil
}

// Bootstrap lets `peer` join a network using a predefined set of nodes.
func (peer *Peer) Bootstrap(port string) {
	// See http://xlattice.sourceforge.net/components/protocol/kademlia/specs.html#join
	n := getBootstrapNode(port)
	peer.Update(n)
	contacts := peer.IterativeFindNode(peer.ID) // [sic! - self-lookup]
	for _, contact := range contacts {
		q := peer.bucketIndex(contact.ID)
		peer.RefreshBucket(q)
	}
}

// RefreshBucket resets the last refresh time for bucket number `q`.
func (peer *Peer) RefreshBucket(q int) {
	peer.refreshMap[q] = time.Now()
}

// FindClosest finds the `n` closest contacts to `targetID` in
// the peer's routing table.
func (peer *Peer) FindClosest(targetID NodeID, n int) []Contact {
	d := Distance(peer.ID, targetID)
	closest := []Contact{}
	seq := NewIntSet()
	seq.AddMany(FindSetBits(d))

	// Descend through 1-bits in `d` toward 0 and try to fill `closest`.
	for _, q := range seq.SortedReverse() {
		bucket := peer.routingTable[q]
		if tryFill(&closest, bucket, n) {
			fmt.Println("Filled up `closest` at bucket", q)
			break
		}
	}

	// If `closest` is still not filled, search unvisisted buckets [0, 160).
	for q := 0; q < b; q++ {
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

func tryFill(closest *[]Contact, bucket Bucket, n int) bool {
	for _, contact := range bucket {
		*closest = append(*closest, contact)
		if len(*closest) == n {
			return true
		}
	}
	return false
}

// Update adds `contact` into the `peer`'s appropriate bucket if necessary.
func (peer *Peer) Update(contact Contact) {
	bucket := peer.bucketFor(contact.ID)

	// If the contact already exists, move it to the end of the bucket.
	for i, c := range *bucket {
		if c.ID == contact.ID {
			bucket.moveToTail(i)
			fmt.Printf("Updated %s, bucket %d with contact %s (tail move).\n",
				peer.ID.String(), peer.bucketIndex(contact.ID), contact.ID.String())
			return
		}
	}

	// If the bucket has space, add the new contact to the bucket.
	if len(*bucket) < k {
		bucket.addToTail(contact)
		fmt.Printf("Updated %s, bucket %d with contact %s (tail add).\n",
			peer.ID.String(), peer.bucketIndex(contact.ID), contact.ID.String())
		return
	}

	// If the bucket is full, ping its head and replace it iff
	// it did not respond within a reasonable time.
	pingChan := make(chan MessageResponsePing)
	peer.SendPing(&contact, (*bucket)[0].ID, pingChan) // TODO just send Contact, see client.go
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
		peer.ID.String(), peer.bucketIndex(contact.ID), contact.ID.String())
}

func (peer *Peer) bucketFor(id NodeID) *Bucket {
	return &peer.routingTable[peer.bucketIndex(id)]
}

func (peer *Peer) bucketIndex(id NodeID) int {
	return b - 1 - Distance(peer.ID, id).PrefixLength()
}

/*
	Some bucket operations
*/

func (bucket *Bucket) moveToTail(i int) {
	*bucket = append((*bucket)[:i], append((*bucket)[i+1:], (*bucket)[i])...)
}

func (bucket *Bucket) addToTail(contact Contact) {
	*bucket = append(*bucket, contact)
}
