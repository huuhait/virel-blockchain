package p2p

import (
	"net"
	"time"
	"virel-blockchain/binary"
	"virel-blockchain/bitcrypto"
	"virel-blockchain/util"
)

func NewConnection(c net.Conn, outgoing bool) *Connection {
	return &Connection{
		data: &ConnData{
			Conn:     c,
			Outgoing: outgoing,
			LastPing: time.Now().Unix(),
		},
		peerData: &PeerData{},
	}
}

// a concurrency-safe wrapper for ConnData
type Connection struct {
	data     *ConnData
	peerData *PeerData

	mut   util.RWMutex
	pdMut util.Mutex
}

func (c *Connection) View(f func(c *ConnData) error) error {
	c.mut.RLock()
	defer c.mut.RUnlock()

	return f(c.data)
}
func (c *Connection) Update(f func(c *ConnData) error) error {
	c.mut.Lock()
	defer c.mut.Unlock()

	return f(c.data)
}
func (c *Connection) PeerData(f func(d *PeerData)) {
	c.pdMut.Lock()
	defer c.pdMut.Unlock()

	f(c.peerData)
}

type ConnData struct {
	Outgoing      bool             // true if this connection is outgoing, false if it is incoming
	PeerId        bitcrypto.Pubkey // the other peer's ID
	LastPing      int64            // last packet received from peer
	LastOutPacket int64            // when the last packet has been sent to the peer
	Cipher        bitcrypto.Cipher

	Conn     net.Conn
	writeMut util.Mutex
}

// returns the connection IP address (without port)
func (c *ConnData) IP() string {
	return c.Conn.RemoteAddr().(*net.TCPAddr).IP.String()
}
func (c *ConnData) Close() {
	c.Conn.Close()
}

func (c *ConnData) sendPacket(p pack) error {
	c.Conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	Log.Debug("sendPacket:", p.String())

	c.LastOutPacket = time.Now().Unix()

	ser := binary.Ser{}

	ser.AddUint16(p.Type)
	ser.AddFixedByteArray(p.Data)

	data, err := c.Cipher.Encrypt(ser.Output())

	if err != nil {
		return err
	}

	ser = binary.Ser{}
	ser.AddUint32(uint32(len(data)))
	ser.AddFixedByteArray(data)

	func() {
		c.writeMut.Lock()
		defer c.writeMut.Unlock()
		_, err = c.Conn.Write(ser.Output())
		if err != nil {
			Log.Warn(err)
		}
	}()
	if err != nil {
		return err
	}
	Log.Debugf("packet sent (%0.3f KB)", float64(len(ser.Output()))/1000)

	return nil
}
func (c *ConnData) SendPacket(p *Packet) error {
	Log.Devf("Sending packet of type %s to peer %x", p.Type.String(), c.PeerId)
	return c.sendPacket(pack{Data: p.Data, Type: uint16(p.Type) + 2})
}

func (c *Connection) SendPacket(p *Packet) error {
	return c.Update(func(c *ConnData) error {
		return c.SendPacket(p)
	})
}
