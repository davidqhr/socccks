package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/davidqhr/sock5/client"
	"github.com/davidqhr/sock5/logger"
)

var VERSION = byte('\x05')
var AUTH_NO = byte('\x00')
var AUTH_USERNAME_PASSWORD = byte('\x02')

var NO_ACCEPTABLE_METHODS = byte('\xFF')

func handleConn(client *client.Client) {
	conn := client.Conn
	defer conn.Close()

	methods, err := handshake(client)

	if err != nil {
		logger.Error(client, "[ERROR]: ", err)
		return
	}

	method := chooseAuthMethod(conn, methods)
	client.AuthMethod = method

	logger.Debug(client, "choose auth method: %X", method)
	conn.Write([]byte{VERSION, method})

	if method == NO_ACCEPTABLE_METHODS {
		return
	}

	ok := authentication(client)

	if !ok {
		logger.Info(client, "authentication failed")
		return
	}
}

func authentication(client *client.Client) bool {
	conn := client.Conn
	switch client.AuthMethod {
	case AUTH_NO:
		return true
	case AUTH_USERNAME_PASSWORD:
		// http://www.rfc-base.org/txt/rfc-1929.txt
		var buf = make([]byte, 2)
		n, err := conn.Read(buf)

		if err != nil {
			logger.Info(client, "read AUTH_USERNAME_PASSWORD error")
			return false
		}

		if n != 2 {
			logger.Info(client, "read username and password len error")
			return false
		}

		username_len := int(buf[1])

		if username_len > 255 || username_len < 1 {
			logger.Error(client, "username size error [1-255]")
			return false
		}

		var username = make([]byte, username_len)
		n, err = conn.Read(username)

		if err != nil {
			logger.Info(client, "read username error: ", err)
			return false
		}

		buf = make([]byte, 1)
		n, err = conn.Read(buf)
		password_len := int(buf[0])

		if password_len > 255 || password_len < 1 {
			logger.Error(client, "password size error [1-255]")
			return false
		}

		var password = make([]byte, password_len)
		n, err = conn.Read(password)

		logger.Info(client, "username: %s, password: %s", string(username), string(password))

		// TODO auth
		conn.Write([]byte("\x01\x00"))

		return true
	}

	return false
}

func chooseAuthMethod(conn net.Conn, methods []byte) byte {
	methods_map := make(map[byte]bool)

	for i := 0; i < len(methods); i++ {
		methods_map[methods[i]] = true
	}

	if methods_map[AUTH_NO] {
		return AUTH_NO
	} else if methods_map[byte('\x01')] {
		return byte('\x01')
	} else if methods_map[AUTH_USERNAME_PASSWORD] {
		return AUTH_USERNAME_PASSWORD
	} else if methods_map[byte('\x03')] {
		return byte('\x03')
	} else if methods_map[byte('\x80')] {
		return byte('\x80')
	} else {
		return byte('\xFF')
	}
}

func handshake(client *client.Client) ([]byte, error) {
	conn := client.Conn
	var buf = make([]byte, 2)

	_, err := conn.Read(buf)

	if err != nil {
		log.Println("conn read error: ", err)
		return make([]byte, 0), err
	}

	// log.Printf("read %d bytes, content is %X\n", string(buf[:n]))

	version := buf[0]

	if version != VERSION {
		return make([]byte, 0), errors.New(fmt.Sprintf("DO NOT SUPPORT VERSION %X", version))
	}

	methods_count := int(buf[1])

	if methods_count == 0 {
		return make([]byte, 0), errors.New(fmt.Sprintf("no auth methods"))
	}

	var methods = make([]byte, methods_count)

	n, err := conn.Read(methods)

	if err != nil {
		// log.Println("read methods error: ", err)
		return make([]byte, 0), err
	}

	if n != methods_count {
		// log.Println("read methods count error: expect(%d) actual(%d)", methods_count, n)
		return make([]byte, 0), err
	}

	log.Printf("methods count: %d, methods: %X", methods_count, methods)

	return methods, nil
}

func serve(addr string) {
	listen, err := net.Listen("tcp", addr)

	if err != nil {
		log.Println("Server start failed")
	} else {
		log.Println("Server Started listen", addr)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept error", err)
			break
		}
		client := client.NewClient(conn)
		go handleConn(client)
	}
}

func main() {
	serve("localhost:8080")
}
