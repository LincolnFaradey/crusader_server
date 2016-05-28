package main

import (
	"log"
	"net"
	"github.com/LincolnFaradey/ConsoleChat/chat"
)

var openConnections = make(map[string]net.Conn)

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

		addr := conn.RemoteAddr().String()
		log.Println("New connection established", addr)
		if _, ok := openConnections[addr]; !ok {
			openConnections[addr] = conn
		}

		go handleConnection(conn)
	}

}


func handleConnection(conn net.Conn) {
	for {
		m := chat.New()

		if _, err := m.ReadFrom(conn); err != nil {
			log.Println("EOF: ", string(m.Content))
			delete(openConnections, conn.RemoteAddr().String())
			break
		}
		addr := conn.RemoteAddr().String()
		log.Println("Response:", string(m.Content), addr)

		go func() {
			for k, v := range openConnections {
				if k == addr {
					self := chat.New()
					self.Kind = []byte{chat.DEBUG}
					self.Content = []byte("Sent")
					if _, err := self.WriteTo(conn); err != nil {
						log.Println("Write Error:", err)
						break
					}
					continue
				}

				if _, err := m.WriteTo(v); err != nil {
					log.Println("Write Error:", err)
					delete(openConnections, v.RemoteAddr().String())
					continue
				}
			}
		}()


	}
}
