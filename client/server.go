package client

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/davidqhr/socccks/utils"
)

func proxyToServer(client *Client, serverAddr string) {
	remoteConn, err := net.Dial("tcp", serverAddr)

	if err != nil {
		log.Println(err)
		return
	}

	eConn := utils.NewEncryptedConn(remoteConn, client.Password)
	// remoteConn, err := net.Dial("tcp", "localhost:8112")

	defer eConn.Close()

	go utils.Copy(eConn, client.Conn)
	utils.Copy(client.Conn, eConn)
}

func handleClient(client *Client, serverAddr string) {
	conn := client.Conn
	defer conn.Close()

	methods, err := client.GetSupportAuthMethods()

	if err != nil {
		log.Println(err.Error())
		return
	}

	if len(methods) == 0 {
		log.Println("no auth methods")
		return
	}

	method := chooseAuthMethod(methods)
	err = client.SetAuthMethod(method)

	if err != nil {
		log.Println("Set Auth Method Failed", err)
		return
	}

	if method == utils.NoAcceptableMethods {
		log.Println("No acceptable methods", method)
		return
	}

	ok := authentication(client)

	if !ok {
		log.Println("Auth failed")
		return
	}

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
	log.Println("Quiting...")
	// wg.Wait()
}
