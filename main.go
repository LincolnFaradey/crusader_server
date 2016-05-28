package main

import (
	"log"
	"net"
	"encoding/binary"
)


func main() {
	log.Println("Running")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}

}

func readMsg(conn net.Conn) ([]byte, error)  {
	header := make([]byte, 8)
	typeBuf := make([]byte, 1)
	if _, err := conn.Read(header); err != nil {
		return nil, err
	}
	if _, err := conn.Read(typeBuf); err != nil {
		return nil, err
	}
	log.Println("Message type:", typeBuf[0])

	l := binary.BigEndian.Uint64(header)

	msg := make([]byte, l)
	if _, err := conn.Read(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func writeMsg(conn net.Conn, msg []byte) error {
	header := make([]byte, 8)
	binary.BigEndian.PutUint64(header, uint64(len(msg)))
	tp := []byte{0}
	ht := append(header, tp[:]...)
	data := append(ht, msg[:]...)
	log.Println("Data:", header, ht, data)
	if _, err := conn.Write(data); err != nil {
		return err
	}

	return nil
}

func handleConnection(conn net.Conn) {
	msgReceived := []byte("Message recieved\n")
	for {
		resp, err := readMsg(conn)
		if err != nil {
			log.Println("EOF: ", string(resp))

			break
		}
		log.Println("Response:", string(resp))

		if err := writeMsg(conn, msgReceived); err != nil {
			log.Println("Write Error:", err)
			break
		}
	}
}
