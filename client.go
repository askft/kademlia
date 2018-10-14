package main

import (
	"net/rpc"
)

/*
	RPC client for the Kademlia protocol (PING, STORE, FIND_NODE, FIND_VALUE).

	TODO
		- Should probably have `peer` as a parameter, not a receiver,
		  in case we want to put this file in a separate package.
		- Uninitialized MessageResponse array values are `nil`. BE CAREFUL!
		- Take Contact, not *Contact
*/

// SendPing sends a PING RPC.
func (peer *Peer) SendPing(contact *Contact, target NodeID, done chan MessageResponsePing) {
	req := &MessageRequestPing{
		MessageCommon: createCommon(peer.ID, GenerateRandomNodeID()),
	}
	res := &MessageResponsePing{}
	err := call(peer, contact, "RPC.RecvPing", req, res)
	if err != nil {
		panic(err)
	}
	done <- *res
}

// SendStore sends a STORE RPC.
// 	TODO send two RPCs - first one to check if it exists already,
//  and if not then send the data.
func (peer *Peer) SendStore(contact *Contact, data []byte, done chan MessageResponseStore) {
	req := &MessageRequestStore{
		MessageCommon: createCommon(peer.ID, GenerateRandomNodeID()),
		Data:          data,
	}
	res := &MessageResponseStore{}
	err := call(peer, contact, "RPC.RecvStore", req, res)
	if err != nil {
		panic(err)
	}
	done <- *res
}

// SendFindNode sends a FIND_NODE RPC.
func (peer *Peer) SendFindNode(contact *Contact, target NodeID, done chan MessageResponseFindNode) {
	req := &MessageRequestFindNode{
		MessageCommon: createCommon(peer.ID, GenerateRandomNodeID()),
		Target:        target,
	}
	res := &MessageResponseFindNode{}
	err := call(peer, contact, "RPC.RecvFindNode", req, res)
	if err != nil {
		panic(err)
	}
	done <- *res
}

// SendFindValue sends a FIND_VALUE_RPC.
func (peer *Peer) SendFindValue(contact *Contact, target NodeID, done chan MessageResponseFindValue) {
	req := &MessageRequestFindValue{
		MessageCommon: createCommon(peer.ID, GenerateRandomNodeID()),
		Target:        target,
	}
	res := &MessageResponseFindValue{}
	err := call(peer, contact, "RPC.RecvFindValue", req, res)
	if err != nil {
		panic(err)
	}
	done <- *res
}

func call(peer *Peer, contact *Contact, method string, args, reply interface{}) error {
	client, err := rpc.Dial("tcp", contact.Address())
	if err != nil {
		return err
	}
	err = client.Call(method, args, reply)
	if err != nil {
		return err
	}
	// peer.AddNode(*contact) // TODO should we really update here?
	return nil
}
