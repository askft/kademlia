package peer

import (
	"fmt"

	"p2p/encoding"
	"p2p/node"
)

/*
	This file contains iterative node lookup functions:
		- IterativeStore
		- IterativeFindNode
		- IterativeFindValue

	NOTE:
		- IterativeFindNode and IterativeFindValue are almost exactly the same.

	TODO
		- Look up if we need WaitGroups
		- xlattice: When an IterativeFindValue succeeds, the initiator
		  must store the key/value pair at the closest node seen which
		  did not return the value.
*/

// IterativeStore finds the <=k closest nodes to `target`
// and sends `data` in a STORE RPC to each.
func (peer *Peer) IterativeStore(target node.Key, data []byte) {
	done := make(chan MessageResponseStore) // TODO empty means error, check all <-done
	contacts := peer.IterativeFindNode(target)
	for _, contact := range contacts {
		go peer.SendStore(contact, data, done)
	}
}

// IterativeFindNode finds the <=k closest nodes to `target`.
func (peer *Peer) IterativeFindNode(target node.Key) []node.Contact {
	var (
		results = []node.Contact{}
		todo    = []node.Contact{}
		seen    = make(map[string]bool)              // the string is Key.String()
		done    = make(chan MessageResponseFindNode) // TODO empty means error, check all <-done
	)
	for _, contact := range peer.FindClosest(target, α) {
		results = append(results, contact)
		todo = append(todo, contact)
		seen[contact.Key.String()] = true
	}

	// Number of pending nodes
	pending := 0

	// Send async FIND_NODE RPCs to α nodes
	for i := 0; i < α && len(todo) > 0; i++ {
		contact := todo[0]
		todo = todo[1:]
		go peer.SendFindNode(contact, target, done) // the recieves node does FindClosest
		pending++
	}

	// While there are still nodes to query
	for pending > 0 {
		res := <-done // Get the RPC result from a node
		pending--

		for _, contact := range res.Contacts {
			// Node hasn't been queried before
			if _, ok := seen[contact.Key.String()]; !ok {
				results = append(results, contact)
				todo = append(todo, contact)
				seen[contact.Key.String()] = true
			}
		}

		// Again, send async FIND_NODE RPCs to α nodes
		for pending < α && len(todo) > 0 {
			contact := todo[0]
			todo = todo[1:]
			go peer.SendFindNode(contact, target, done)
			pending++
		}
	}
	node.SortByDistance(results, target)
	if len(results) > k {
		results = results[:k]
	}
	return results
}

// IterativeFindValue attemps to find the value at `target`. If the value
// can't be found, the <=k closest nodes to `target` are returned.
func (peer *Peer) IterativeFindValue(target node.Key) ([]byte, []node.Contact) {
	var (
		results = []node.Contact{}
		todo    = []node.Contact{}
		seen    = make(map[string]bool)               // the string is Key.String()
		done    = make(chan MessageResponseFindValue) // TODO empty means error, check all <-done
	)
	for _, contact := range peer.FindClosest(target, α) {
		results = append(results, contact)
		todo = append(todo, contact)
		seen[contact.Key.String()] = true
	}

	// Number of pending nodes
	pending := 0

	// Send async FIND_VALUE RPCs to α nodes
	for i := 0; i < α && len(todo) > 0; i++ {
		contact := todo[0]
		todo = todo[1:]
		go peer.SendFindValue(contact, target, done) // the recieves node does FindClosest
		pending++
	}

	// While there are still nodes to query
	for pending > 0 {
		res := <-done // Get the RPC result from a node
		pending--

		// If a value was found, return it immediately
		if res.Data != nil || len(res.Data) > 0 {
			// TODO this condition will always be true if we get here
			if encoding.EncodeHash(target) == encoding.EncodeData(res.Data) {
				fmt.Println("found value, returning")
				// TODO store in cache, see top of this file
				return res.Data, nil
			}
			fmt.Println("this should not print. value was found, but not the correct one. search continues...")
		}

		for _, contact := range res.Contacts {
			// Contact hasn't been queried before
			if _, ok := seen[contact.Key.String()]; !ok {
				results = append(results, contact)
				todo = append(todo, contact)
				seen[contact.Key.String()] = true
			}
		}

		// Again, send async FIND_VALUE RPCs to α nodes
		for pending < α && len(todo) > 0 {
			contact := todo[0]
			todo = todo[1:]
			go peer.SendFindValue(contact, target, done)
			pending++
		}
	}
	node.SortByDistance(results, target)
	if len(results) > k {
		results = results[:k]
	}
	return nil, results
}
