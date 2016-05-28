package chat

import (
	"encoding/binary"
	"net"
)

type Kind byte

const (
	TEXT = iota
	FILE
	INFO
	DEBUG
)

type Message struct {
	Header  []byte
	Kind    []byte
	Content []byte
}

func New() *Message {
	return &Message{
		Header: make([]byte, 8),
		Kind: make([]byte, 1),
	}
}

func (m *Message)create() {
	m.Header = make([]byte, 8)
	m.Kind = make([]byte, 1)
}


func (m *Message)ReadFrom(conn net.Conn) (int, error) {
	m.create()
	if _, err := conn.Read(m.Header); err != nil {
		return 0, err
	}

	if _, err := conn.Read(m.Kind); err != nil {
		return len(m.Header), err
	}
	l := binary.BigEndian.Uint64(m.Header)
	m.Content = make([]byte, l)
	if _, err := conn.Read(m.Content); err != nil {
		return len(m.Header) + len(m.Kind), err
	}

	return len(m.Header) + len(m.Kind) + len(m.Content), nil
}

func (m *Message)WriteTo(conn net.Conn) (int, error) {
	binary.BigEndian.PutUint64(m.Header, uint64(len(m.Content)))
	ht := append(m.Header, m.Kind[:]...)
	data := append(ht, m.Content[:]...)
	return conn.Write(data)
}