package main

import (
	"net"
)

const (
	idLength = 20           // The byte length of a node ID.
	α        = 3            // Parallelism parameter for RPC calls.
	b        = idLength * 8 // Bit length of node ID and key.
	k        = 20           // Bucket size.

	updateTimeout = 1000 // Milliseconds
)

// Options contains general configuration parameters for a peer.
type Options struct {
	id        NodeID
	host      net.IP
	port      string
	store     Store
	networkID string
}

var defaultOptions = &Options{
	id:        GenerateRandomNodeID(),
	host:      getLocalIP(),
	port:      "4001",
	store:     NewLocalStore(),
	networkID: "v1",
}

// TimeOptions contains time-specific configuration parameters for a peer.
type TimeOptions struct {
	Expire    uint // TTL for KV pair from original publication date
	Refresh   uint // Time until an unaccessed bucket must be refreshed
	Replicate uint // Interval between replication events, when a node is required to publish its entire database
	Republish uint // Time after which original publisher must republish a KV pair
}

var timeOptions = TimeOptions{
	Expire:    86410, // Longer than tReplublish as per the xlattice spec
	Refresh:   3600,
	Replicate: 3600,
	Republish: 86400,
}
