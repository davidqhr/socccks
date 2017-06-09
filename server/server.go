package server

import (
	"fmt"
	"log"

	"github.com/davidqhr/socccks/utils"
)

func handleConn(eConn *utils.EncryptedConn) {
	buf := utils.BufPool.Get()
	defer utils.BufPool.Put(buf)

	_, err := eConn.Read(buf)

	if err != nil {
		println(err)
		return
	}

	cmd := buf[1]
	fmt.Printf("[debug] command received: %v\n", cmd)

	switch cmd {
	case utils.CmdConnect:
		handleCmdConnection(eConn, buf[3:])
	}
}

// start socccks server
func Start(addr string, password string) {
	connections := utils.StartAccepter(addr, 100)

	log.Printf("Listen on: %s, ( poolSize: %d )\n", addr, 100)

	for conn := range connections {
		eConn := utils.NewEncryptedConn(conn, password)
		go handleConn(eConn)
	}
}
