package socks5

import (
	"log"
	"net"
)

func handleRequest(client *Client) {
	conn := client.Conn
	var buf = make([]byte, 3)
	conn.Read(buf)

	cmd := buf[1]

	switch cmd {
	case CmdConnect:
		handleCmdConnection(client)
	}
}

func handleClient(client *Client) {
	conn := client.Conn
	defer conn.Close()

	methods, err := client.GetSupportAuthMethods()

	if err != nil {
		logger.Error(client, "%s", err.Error())
		return
	}

	if len(methods) == 0 {
		logger.Error("no auth methods")
		return
	}

	method := chooseAuthMethod(methods)
	err = client.SetAuthMethod(method)

	if err != nil {
		logger.Error("Set Auth Method Failed", err)
		return
	}

	if method == NoAcceptableMethods {
		return
	}

	ok := authentication(client)

	if !ok {
		logger.Error("Authentication Failed")
		return
	}

	handleRequest(client)
}

func chooseAuthMethod(methods []byte) byte {
	methods_map := make(map[byte]bool)

	for i := 0; i < len(methods); i++ {
		methods_map[methods[i]] = true
	}

	// only support no_auth or username_password_auth
	if methods_map[AuthNo] {
		return AuthNo
	} else if methods_map[AptyDomainName] {
		return AptyDomainName
	} else {
		return NoAcceptableMethods
	}
}

func Serve(addr string) {
	listen, err := net.Listen("tcp", addr)

	if err != nil {
		log.Println("Server start failed")
	} else {
		log.Println("Server Started listen", addr)
	}

	for {
		conn, err := listen.Accept()

		if err != nil {
			break
		}

		client := NewClient(conn)
		go handleClient(client)
	}
}
