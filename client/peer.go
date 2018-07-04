package client

import (
	"encoding/binary"
	"errors"
	"net"
)

type PeerMessageType byte

const (
	Choke   PeerMessageType = iota
	Unchoke                 = iota + 1
	Interested
	Uninterested
	Have
	Bitfield
	Request
	Piece
	Cancel
)

type Peer struct {
	Id      string
	Address net.TCPAddr

	Connection  net.Conn
	IsAvailable bool
}

// OpenPeerConnection attempts to open a connection to the given address
// and executes the given callback function if it succeeds. Non-blocking
func (c *Peer) OpenConnection(callback func(*Peer)) error {
	conn, err := net.DialTCP("tcp", nil, &c.Address)

	if err != nil {
		return errors.New("Failed to open connection to peer: " + err.Error())
	}

	conn.SetKeepAlive(true)
	c.Connection = conn

	go callback(c)
	return nil
}

func (c *Peer) IsChoked() bool {
	return true
}

func (c *Peer) IsInterested() bool {
	return false
}

func (c *Peer) SendStateUpdate(update PeerMessageType) error {
	return c.sendToPeer(update, nil)
}

func (c *Peer) SendBitfield() error {
	return c.sendToPeer(Bitfield, nil)
}

func (c *Peer) SendHave(index int) error {
	indexBytes, err := intToBigEndianBytes(index)
	if err != nil {
		return err
	}

	return c.sendToPeer(Have, indexBytes)
}

func (c *Peer) SendRequest(index, begin, length int) error {
	indexBytes, indErr := intToBigEndianBytes(index)
	beginBytes, begErr := intToBigEndianBytes(begin)
	lengthBytes, lenErr := intToBigEndianBytes(length)
	if indexBytes != nil {
		return indErr
	}

	if begErr != nil {
		return begErr
	}

	if lenErr != nil {
		return lenErr
	}

	message := append(append(indexBytes, beginBytes...), lengthBytes...)
	return c.sendToPeer(Request, message)
}

func (c *Peer) SendPiece(index, begin int, block []byte) error {
	indexBytes, indErr := intToBigEndianBytes(index)
	beginBytes, begErr := intToBigEndianBytes(begin)
	if indexBytes != nil {
		return indErr
	}

	if begErr != nil {
		return begErr
	}

	message := append(indexBytes, beginBytes...)
	message = append(message, block...)
	return c.sendToPeer(Piece, message)
}

func (c *Peer) SendCancel(index, begin, length int) error {
	indexBytes, indErr := intToBigEndianBytes(index)
	beginBytes, begErr := intToBigEndianBytes(begin)
	lengthBytes, lenErr := intToBigEndianBytes(length)
	if indexBytes != nil {
		return indErr
	}

	if begErr != nil {
		return begErr
	}

	if lenErr != nil {
		return lenErr
	}

	message := append(indexBytes, beginBytes...)
	message = append(message, lengthBytes...)
	return c.sendToPeer(Cancel, message)
}

func (c *Peer) sendToPeer(typeId PeerMessageType, data []byte) error {
	message := []byte{byte(typeId)}
	message = append(message, data...)
	payload, err := addLengthPrefix(message)
	if err != nil {
		return errors.New("Failed to add message length: " + err.Error())
	}

	addr := c.Connection.RemoteAddr().String()
	log.Debugf("Sending message to peer %s: %s", addr, string(payload))
	_, err = c.Connection.Write(payload)
	if err != nil {
		return errors.New("Comm error occurred with peer " + addr + ": " + err.Error())
	}

	return nil
}

func addLengthPrefix(message []byte) ([]byte, error) {
	prefix, err := intToBigEndianBytes(len(message))
	if err != nil {
		return nil, err
	}

	return append(prefix, message...), nil
}

func intToBigEndianBytes(data int) ([]byte, error) {
	bytes := make([]byte, 4)
	if uint64(data) > uint64(2<<32) {
		return nil, errors.New("Integer cannot be larger than 4 bytes")
	}
	binary.BigEndian.PutUint32(bytes, uint32(data))
	return bytes, nil
}
