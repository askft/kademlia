package main

type MessageCommon struct {
	Sender NodeID
	Nonce  NodeID
}

func createCommon(sender, nonce NodeID) MessageCommon {
	return MessageCommon{
		Sender: sender,
		Nonce:  nonce,
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
	Target NodeID
}

type MessageResponseFindNode struct {
	MessageCommon
	Contacts []Contact
}

type MessageRequestFindValue struct {
	MessageCommon
	Target NodeID
}

// MessageResponseFindValue NOTE: Either Contacts or Data should be empty.
type MessageResponseFindValue struct {
	MessageCommon
	Contacts []Contact
	Data     []byte
}

// This is unnecessary, will be returned from the RPC call instead.
type MessageResponseError struct {
	MessageCommon
	Err error
}
