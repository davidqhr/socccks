package utils

import (
	"fmt"
	"log"
	"net"
)

func StartAccepter(addr string, connectionsPoolSize int) (connections chan net.Conn) {
	listen, err := net.Listen("tcp", addr)
	connections = make(chan net.Conn, connectionsPoolSize)

	if err != nil {
		panic(fmt.Sprintf("Listen failed: %s", addr))
	} else {
		log.Println("Listen on", addr, "(poolSize:", connectionsPoolSize, ")")
	}

	go func() {
		for {
			conn, err := listen.Accept()

			if err != nil {
				log.Println("Accept Error")
				break
			}

			println("new connections in")

			connections <- conn
			// (*wg).Add(1)
		}
	}()

	return
}
