package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/davidqhr/socccks/utils"
)

func handleConn(conn net.Conn) {
	encryptor := utils.NewEncryptor("test")
	data, err := utils.ReadThenDecrypt(conn, nil, encryptor)

	if err != nil {
		println(err)
		return
	}

	cmd := data[1]
	fmt.Printf("[debug] command received: %v\n", cmd)

	switch cmd {
	case utils.CmdConnect:
		handleCmdConnection(conn, data[3:])
	}
}

// start socccks server
func Start(addr string) {
	println("Pid: ", os.Getpid())
	connections := utils.StartAccepter(addr, 100)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGINT)

	// serve conn from connections until connections closed
	go func(chan net.Conn) {
		for conn := range connections {
			go handleConn(conn)
		}
	}(connections)

	// wait signal to close connections
	<-quit
	close(connections)

	// graceful exit
	// TODO: client timeout
	println("Quiting...")
	// wg.Wait()
}
