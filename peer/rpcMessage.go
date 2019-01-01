package peer

import (
	"github.com/askft/kademlia/node"
)

type MessageCommon struct {
	Sender node.Contact
	Nonce  node.Key
}

func createCommon(sender node.Contact, nonce node.Key) MessageCommon {
	return MessageCommon{
		Sender: sender,
		Nonce:  nonce,
	}
}

func createCommonWithNonce(sender node.Contact) MessageCommon {
	return MessageCommon{
		Sender: sender,
		Nonce:  node.GenerateRandomKey(),
	}
}

type MessageRequestPing struct {
	MessageCommon
}

type MessageResponsePing struct {
	MessageCommon
}

type MessageRequestStore struct {
	MessageCommon
	Data []byte
}

type MessageResponseStore struct {
	MessageCommon
}

type MessageRequestFindNode struct {
	MessageCommon
	Target node.Key
}

type MessageResponseFindNode struct {
	MessageCommon
	Contacts []node.Contact
}

type MessageRequestFindValue struct {
	MessageCommon
	Target node.Key
}

// MessageResponseFindValue NOTE: Either Contacts or Data should be empty.
type MessageResponseFindValue struct {
	MessageCommon
	Contacts []node.Contact
	Data     []byte
}
