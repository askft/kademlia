package peer

import (
	"net"
	"time"

	"github.com/askft/kademlia/node"
	"github.com/askft/kademlia/store"
)

const (
	α = 3  // Parallelism parameter for RPC calls.
	k = 20 // Bucket size.

	updateTimeout = 1000 // Milliseconds
)

// Options contains general configuration parameters for a peer.
type Options struct {
	Key       node.Key
	Host      net.IP
	Port      string
	Store     store.Store
	NetworkID string
}

// TimeOptions contains time-specific configuration parameters for a peer.
type TimeOptions struct {
	Expire    time.Duration // TTL for KV pair from original publication date
	Refresh   time.Duration // Time until an unaccessed bucket must be refreshed
	Replicate time.Duration // Interval between replication events, when a node is required to publish its entire database
	Republish time.Duration // Time after which original publisher must republish a KV pair
}

var timeOptions = TimeOptions{
	Expire:    86410, // Longer than tReplublish as per the xlattice spec
	Refresh:   3600,
	Replicate: 3600,
	Republish: 86400,
}
