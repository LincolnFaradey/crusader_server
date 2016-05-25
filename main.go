package main

import (
	"log"
	"net"
	"bufio"
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
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')

	return []byte(line), err
}

func writeMsg(conn net.Conn, msg []byte) error {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(len(msg)))

	data := append(bs, msg[:]...)
	log.Println("Data:", bs, data)
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
