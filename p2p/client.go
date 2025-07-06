package p2p

import (
	"fmt"
	"net"
)

const MAX_PEER_FAILURES = 5

// P2P must NOT be locked before calling this
func (p2 *P2P) startClient(addr string, port uint16, private bool) {
	conn, err := p2.connectClient(net.JoinHostPort(addr, fmt.Sprintf("%d", port)))
	if err != nil {
		Log.Net("error connecting to", addr+":", err)

		func() {
			p2.Lock()
			defer p2.Unlock()

			for i, kp := range p2.KnownPeers {
				if kp.IP == addr && kp.Port == port {
					p2.KnownPeers[i].Fails++

					if p2.KnownPeers[i].Fails > MAX_PEER_FAILURES {
						// Remove the peer, it has too many failures
						p2.KnownPeers[i] = p2.KnownPeers[len(p2.KnownPeers)-1]
						p2.KnownPeers = p2.KnownPeers[:len(p2.KnownPeers)-1]

						return
					}

					p2.KnownPeers[i] = kp
					return
				}
			}
		}()

		return
	} else {
		go func() {
			p2.Lock()
			defer p2.Unlock()

			for i, kp := range p2.KnownPeers {
				if kp.IP == addr && kp.Port == port {
					p2.KnownPeers[i].Fails = 0
					if p2.KnownPeers[i].Type == PEER_GRAY {
						p2.KnownPeers[i].Type = PEER_WHITE
					}

					p2.KnownPeers[i] = kp
					return
				}
			}
		}()
	}

	err = p2.handleConnection(conn, private)
	if err != nil {
		Log.Debug("P2P client connection error:", err)
	} else {
		Log.Debug("P2P client connection gracefully closed")
	}
}

func (p2 *P2P) connectClient(addr string) (*Connection, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		Log.Net(err)
		return &Connection{}, err
	}
	return NewConnection(conn, true), nil
}
