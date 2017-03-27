package socks5

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func handleCmdConnection(client *Client) {
	conn := client.Conn
	buf := make([]byte, 100)

	_, err := conn.Read(buf)

	if err != nil {
		return
	}

	addrType := buf[0]
	var addr string
	port_bytes := make([]byte, 2)

	// logger.Debug(client, "addrType: %X", addrType)
	switch addrType {
	case ATYP_IPV4:
		ipv4_bytes := buf[1:5]
		port_bytes = buf[5:7]
		addr = bytesToIpv4String(ipv4_bytes)
	case APTY_IPV6:
		logger.Error(client, "NOT IMPLEMENTED APTY_IPV6")
		return
	case APTY_DOMAINNAME:
		domainLen := uint8(buf[1])
		addr = string(buf[2 : 2+domainLen])
		port_bytes = buf[2+domainLen : 2+domainLen+2]
		// logger.Error(client, "NOT IMPLEMENTED APTY_DOMAINNAME")
		// return
	}

	port := binary.BigEndian.Uint16(port_bytes)

	logger.Debug(client, "addr: %s, port: %d", addr, port)

	remoteConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", string(addr), port))
	defer remoteConn.Close()

	if err != nil {
		logger.Info(client, "Connect remote Failed %s", conn.LocalAddr().String())
		return
	}

	addr_and_port := make([]string, 2)
	addr_and_port = strings.Split(conn.LocalAddr().String(), ":")
	dstAddr := addr_and_port[0]
	dstPort := addr_and_port[1]

	dstPortBytes := make([]byte, 2)
	dstPortInt, err := strconv.Atoi(dstPort)

	binary.BigEndian.PutUint16(dstPortBytes, uint16(dstPortInt))

	data := []byte{VERSION, REPLY_SUCCESS, RSV, ATYP_IPV4}
	data = append(data, ipv4StringToBytes(dstAddr)...)
	data = append(data, dstPortBytes...)

	logger.Debug(client, "send data %X", data)
	conn.Write(data)

	go proxyTcp(conn, remoteConn)
	proxyTcp(remoteConn, conn)
}

func proxyTcp(src net.Conn, dst net.Conn) {
	buf := bufferPool.Get()
	defer bufferPool.Put(buf)

	for {
		n, err := src.Read(buf)

		if err != nil {
			break
		}

		n, err = dst.Write(buf[:n])

		if err != nil {
			break
		}
	}
}

func handleRequest(client *Client) {
	conn := client.Conn
	var buf = make([]byte, 3)
	conn.Read(buf)

	cmd := buf[1]

	switch cmd {
	case CMD_CONNECT:
		handleCmdConnection(client)
	}
}

func handleConn(client *Client) {
	conn := client.Conn
	defer conn.Close()

	methods, err := client.GetSupportAuthMethods()

	if err != nil {
		logger.Error(client, "%s", err.Error())
		return
	}

	if len(methods) == 0 {
		logger.Error(client, "no auth methods")
	}

	method := chooseAuthMethod(methods)
	err = client.SetAuthMethod(method)

	if method == NO_ACCEPTABLE_METHODS {
		return
	}

	ok := authentication(client)

	if !ok {
		logger.Info(client, "Authentication Failed")
		return
	}

	handleRequest(client)
}

func authentication(client *Client) bool {
	conn := client.Conn
	switch client.AuthMethod {
	case AUTH_NO:
		return true
	case AUTH_USERNAME_PASSWORD:
		// http://www.rfc-base.org/txt/rfc-1929.txt
		var buf = make([]byte, 513)
		_, err := conn.Read(buf)

		if err != nil {
			logger.Info(client, "read AUTH_USERNAME_PASSWORD error")
			return false
		}

		username_len := int(buf[1])

		if username_len > 255 || username_len < 1 {
			logger.Error(client, "username size error [1-255]")
			return false
		}

		username := buf[1+1 : 1+1+username_len]

		password_len := int(buf[1+1+username_len])

		if password_len > 255 || password_len < 1 {
			logger.Error(client, "password size error [1-255]")
			return false
		}

		password := buf[1+1+username_len+1 : 1+1+username_len+1+password_len]

		logger.Info(client, "username: %s, password: %s", string(username), string(password))

		// TODO real auth
		err = client.AuthSuccess()

		if err != nil {
			return false
		}

		return true
	}

	return false
}

func chooseAuthMethod(methods []byte) byte {
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
		return NO_ACCEPTABLE_METHODS
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
			fmt.Printf("accept error", err)
			break
		}

		client := NewClient(conn)
		go handleConn(client)
	}
}
