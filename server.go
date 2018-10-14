package main

import (
	"fmt"
	"net"
	"net/rpc"
)

/*
	Kademlia protocol
		- PING 		 : probe a node to see if it is online
		- STORE		 : store a (key, value) pair in one node
		- FIND_NODE	 : recipient returns k closest nodes to requested key
		- FIND_VALUE : like FIND_NODE, but return value if found in node
*/

// TODO - from xlattice:
// "Whenever a node receives a communication from another,
// it updates the corresponding bucket."
// Should probably do that inside the RecvX...

func (r *RPC) RecvPing(req *MessageRequestPing, res *MessageResponsePing) error {
	fmt.Println("RecvPing")
	res.Sender = r.peer.ID
	res.Nonce = req.Nonce
	return nil
}

func (r *RPC) RecvStore(req *MessageRequestStore, res *MessageResponseStore) error {
	fmt.Println("RecvStore")
	res.Sender = r.peer.ID
	res.Nonce = req.Nonce
	key, err := r.peer.store.Put(req.Data) // TODO should also store req.Sender
	if err != nil {
		return err
	}
	fmt.Printf("stored data at %s", key)
	return nil
}

func (r *RPC) RecvFindNode(req *MessageRequestFindNode, res *MessageResponseFindNode) error {
	fmt.Println("RecvFindNode")
	res.Sender = r.peer.ID
	res.Nonce = req.Nonce
	res.Contacts = r.peer.FindClosest(req.Target, k)
	return nil
}

func (r *RPC) RecvFindValue(req *MessageRequestFindValue, res *MessageResponseFindValue) error {
	fmt.Println("RecvFindValue")
	res.Sender = r.peer.ID
	res.Nonce = req.Nonce
	data, err := r.peer.store.Get(encodeID(req.Target))
	if err != nil {
		fmt.Println("data not found")
		return err
	}
	if data != nil {
		res.Data = data
		return nil
	}
	res.Contacts = r.peer.FindClosest(req.Target, k)
	return nil
}

// RPC is the receiver required by net/rpc.
type RPC struct {
	peer *Peer
}

// RunServer starts an RPC server that listens for RPC calls from other peers.
func RunServer(peer *Peer) error {
	defer wg.Done()
	rpc.Register(&RPC{peer})
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+peer.Port)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	fmt.Printf("Started RPC server on port %s.\n", peer.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(conn)
	}
}
