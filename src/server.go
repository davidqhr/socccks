package socks5

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var wg sync.WaitGroup

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
	defer wg.Done()

	conn := client.Conn
	defer conn.Close()

	methods, err := client.GetSupportAuthMethods()

	if err != nil {
		logger.Error(err.Error())
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
		logger.Error("No acceptable methods", method)
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

func serve(connections chan net.Conn) {

	for conn := range connections {
		client := NewClient(conn)
		go handleClient(client)
	}
}

func startAccepter(addr string) chan net.Conn {
	connections := make(chan net.Conn)
	listen, err := net.Listen("tcp", addr)

	if err != nil {
		log.Println("Listen failed", err)
		panic("Listen failed")
	} else {
		log.Println("Listen on", addr)
	}

	go func() {
		for {
			conn, err := listen.Accept()

			if err != nil {
				log.Println("Accept Error")
				break
			}

			connections <- conn
			wg.Add(1)
		}
	}()

	return connections
}

func Start(addr string) {
	println("Pid: ", os.Getpid())
	connections := startAccepter(addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGINT)

	go serve(connections)

	<-quit
	close(connections)

	logger.Info("Quiting...")
	wg.Wait()
}
