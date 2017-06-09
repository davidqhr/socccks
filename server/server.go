package server

import (
	"log"

	"github.com/davidqhr/socccks/utils"
)

func handleConn(eConn *utils.EncryptedConn) {
	defer eConn.Close()

	buf := utils.Pool33K.Get()
	defer utils.Pool33K.Put(buf)

	_, err := eConn.Read(buf)

	if err != nil {
		log.Println(err)
		return
	}

	cmd := buf[1]

	switch cmd {
	case utils.CmdConnect:
		handleCmdConnection(eConn, buf[3:])
	}
}

// start socccks server
func Start(addr string, password string) {
	connections := utils.StartAccepter(addr, 100) // poolSize

	log.Printf("Listen on: %s\n", addr)

	for conn := range connections {
		eConn := utils.NewEncryptedConn(conn, password)
		go handleConn(eConn)
	}
}
