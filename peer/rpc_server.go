package peer

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"

	"github.com/pkg/errors"

	"github.com/askft/kademlia/encoding"
)

/*
	Kademlia protocol
		- PING 		 : probe a node to see if it is online
		- STORE		 : store a (key, value) pair in one node
		- FIND_NODE	 : recipient returns k closest nodes to requested key
		- FIND_VALUE : like FIND_NODE, but return value if found in node
*/

// RPC is the receiver required by net/rpc.
type RPC struct {
	peer *Peer
}

// RecvPing signals to the sender that this peer is online.
func (r *RPC) RecvPing(req *MessageRequestPing, res *MessageResponsePing) error {
	fmt.Println("RecvPing")
	r.peer.UpdateTable(req.Sender)
	res.Sender = r.peer.Contact
	res.Nonce = req.Nonce
	return nil
}

// RecvStore stores a key-value pair at this peer.
func (r *RPC) RecvStore(req *MessageRequestStore, res *MessageResponseStore) error {
	fmt.Println("RecvStore")
	r.peer.UpdateTable(req.Sender)
	res.Sender = r.peer.Contact
	res.Nonce = req.Nonce
	key, err := r.peer.Put(req.Data) // TODO should also store req.Sender
	if err != nil {
		return err
	}
	fmt.Printf("stored data at %s", key)
	return nil
}

// RecvFindNode returns `k` closest nodes to requested key.
func (r *RPC) RecvFindNode(req *MessageRequestFindNode, res *MessageResponseFindNode) error {
	fmt.Printf("RecvFindNode from [ %s ].\n", req.Sender.Address())
	r.peer.UpdateTable(req.Sender)
	res.Sender = r.peer.Contact
	res.Nonce = req.Nonce
	res.Contacts = r.peer.FindClosest(req.Target, k)
	return nil
}

// RecvFindValue returns value at key if found, else returns `k` closest nodes to key.
func (r *RPC) RecvFindValue(req *MessageRequestFindValue, res *MessageResponseFindValue) error {
	fmt.Println("RecvFindValue")
	r.peer.UpdateTable(req.Sender)
	res.Sender = r.peer.Contact
	res.Nonce = req.Nonce
	data, err := r.peer.Get(encoding.EncodeHash(req.Target))
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

type Server struct {
	port     string
	addr     *net.TCPAddr
	listener *net.TCPListener
}

func NewServer(peer *Peer) (*Server, error) {
	err := rpc.Register(&RPC{peer})
	if err != nil {
		return nil, err
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+peer.Contact.Port)
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}

	return &Server{
		peer.Contact.Port,
		tcpAddr,
		listener,
	}, nil
}

// RunServer starts an RPC server that listens for RPC calls from other peers.
func (s *Server) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Starting RPC server on port %s.\n", s.port)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(errors.Wrap(err, "failed to connect"))
			continue
		}
		go rpc.ServeConn(conn)
	}
}
