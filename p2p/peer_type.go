package p2p

type PeerType uint8

const (
	// a "white" peer is a peer that was working on last connection
	PEER_WHITE PeerType = iota
	// a "gray" peer is a peer that hasn't been contacted by this node yet
	PEER_GRAY
	// a "black" peer is a peer that was unreachable on last connection attempt
	PEER_BLACK
	// a "red" peer is a banned peer
	PEER_RED
)
