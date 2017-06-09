package utils

import (
	"log"
	"net"
)

func StartAccepter(addr string, connectionsPoolSize int) (connections chan net.Conn) {
	listen, err := net.Listen("tcp", addr)
	connections = make(chan net.Conn, connectionsPoolSize)

	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			conn, err := listen.Accept()

			if err != nil {
				log.Println("[DEBUG] accept error")
				break
			}

			log.Println("[DEBUG] new connections in")

			connections <- conn
		}
	}()

	return
}
