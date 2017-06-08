package client

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/davidqhr/socccks/utils"
)

func proxyToServer(client *Client, serverAddr string) {
	remoteConn, err := net.Dial("tcp", serverAddr)

	if err != nil {
		logger.Info(err)
		return
	}

	eConn := utils.NewEncryptedConn(remoteConn, client.Password)
	// remoteConn, err := net.Dial("tcp", "localhost:8112")

	defer eConn.Conn.Close()

	go io.CopyBuffer(eConn, client.Conn)
	io.CopyBuffer(client.Conn, eConn)
}

func handleClient(client *Client, serverAddr string) {
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
	fmt.Printf("[debug] support auth methods: %v, choose auth method: %v\n", methods, method)

	err = client.SetAuthMethod(method)

	if err != nil {
		logger.Error("Set Auth Method Failed", err)
		return
	}

	if method == utils.NoAcceptableMethods {
		logger.Error("No acceptable methods", method)
		return
	}

	ok := authentication(client)

	if !ok {
		println("[debug] Auth failed\n")
		return
	}
	println("[debug] Auth success\n")

	proxyToServer(client, serverAddr)
}

func chooseAuthMethod(methods []byte) byte {
	methods_map := make(map[byte]bool)

	for i := 0; i < len(methods); i++ {
		methods_map[methods[i]] = true
	}

	// only support no_auth or username_password_auth
	if methods_map[utils.AuthNo] {
		return utils.AuthNo
	} else if methods_map[utils.AptyDomainName] {
		return utils.AptyDomainName
	} else {
		return utils.NoAcceptableMethods
	}
}

func Start(localAddr string, serverAddr string, password string) {
	connections := utils.StartAccepter(localAddr, 100)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGINT)

	// serve conn from connections until connections closed
	go func(*chan net.Conn) {
		for conn := range connections {
			client := NewClient(conn, password)
			go handleClient(client, serverAddr)
		}
	}(&connections)

	// wait signal to close connections
	<-quit
	close(connections)

	// graceful exit
	// TODO: client timeout
	logger.Info("Quiting...")
	// wg.Wait()
}
