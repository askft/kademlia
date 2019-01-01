package peer

import (
	"net/rpc"

	"github.com/askft/kademlia/node"
)

/*
	RPC client for the Kademlia protocol (PING, STORE, FIND_NODE, FIND_VALUE).

	TODO
		- Uninitialized MessageResponse array values are `nil`. BE CAREFUL!
*/

// SendPing sends a PING RPC.
func (peer *Peer) SendPing(contact node.Contact, target node.Key, done chan MessageResponsePing) {
	req := &MessageRequestPing{
		MessageCommon: createCommonWithNonce(peer.Contact),
	}
	res := &MessageResponsePing{}
	err := peer.call(contact, "RPC.RecvPing", req, res)
	if err != nil {
		panic(err)
	}
	done <- *res
	peer.UpdateTable(res.Sender)
}

// SendStore sends a STORE RPC.
// 	TODO send two RPCs - first one to check if it exists already,
//  and if not then send the data.
func (peer *Peer) SendStore(contact node.Contact, data []byte, done chan MessageResponseStore) {
	req := &MessageRequestStore{
		MessageCommon: createCommonWithNonce(peer.Contact),
		Data:          data,
	}
	res := &MessageResponseStore{}
	err := peer.call(contact, "RPC.RecvStore", req, res)
	if err != nil {
		panic(err)
	}
	done <- *res
	peer.UpdateTable(res.Sender)
}

// SendFindNode sends a FIND_NODE RPC.
func (peer *Peer) SendFindNode(contact node.Contact, target node.Key, done chan MessageResponseFindNode) {
	req := &MessageRequestFindNode{
		MessageCommon: createCommonWithNonce(peer.Contact),
		Target:        target,
	}
	res := &MessageResponseFindNode{}
	err := peer.call(contact, "RPC.RecvFindNode", req, res)
	if err != nil {
		panic(err)
	}
	done <- *res
	peer.UpdateTable(res.Sender)
}

// SendFindValue sends a FIND_VALUE_RPC.
func (peer *Peer) SendFindValue(contact node.Contact, target node.Key, done chan MessageResponseFindValue) {
	req := &MessageRequestFindValue{
		MessageCommon: createCommonWithNonce(peer.Contact),
		Target:        target,
	}
	res := &MessageResponseFindValue{}
	err := peer.call(contact, "RPC.RecvFindValue", req, res)
	if err != nil {
		panic(err)
	}
	done <- *res
	peer.UpdateTable(res.Sender)
}

func (peer *Peer) call(contact node.Contact, method string, args, reply interface{}) error {
	client, err := rpc.Dial("tcp", contact.Address())
	if err != nil {
		return err
	}
	return client.Call(method, args, reply)
}
