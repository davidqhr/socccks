package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/davidqhr/socccks/utils"
)

func handleCmdConnection(eConn *utils.EncryptedConn, buf []byte) {
	addrType := buf[0]
	var addr string
	portBytes := make([]byte, 2)

	switch addrType {
	case utils.AptyIPV4:
		ipv4Bytes := buf[1:5]
		portBytes = buf[5:7]
		addr = net.IP(ipv4Bytes).String()
	case utils.AptyIPV6:
		println("NOT IMPLEMENTED APTY_IPV6")
		return
	case utils.AptyDomainName:
		domainLen := uint8(buf[1])
		addr = string(buf[2 : 2+domainLen])
		portBytes = buf[2+domainLen : 2+domainLen+2]
	}

	port := binary.BigEndian.Uint16(portBytes)
	println("addr", addr, "port", port)

	remoteConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", string(addr), port))
	if err != nil {
		println("Connect remote Failed %s", eConn.Conn.LocalAddr().String())
		return
	}

	defer remoteConn.Close()

	// addrAndPort := strings.Split(remoteConn.LocalAddr().String(), ":")
	// dstAddr := addrAndPort[0]
	// dstPort := addrAndPort[1]
	// //
	// dstPortBytes := make([]byte, 2)
	// dstPortInt, err := strconv.Atoi(dstPort)
	// //
	// if err != nil {
	// 	println(err)
	// 	return
	// }
	//
	// binary.BigEndian.PutUint16(dstPortBytes, uint16(dstPortInt))

	data := []byte{utils.Version, utils.ReplySuccess, utils.Rsv, utils.AptyIPV4, 0, 0, 0, 0, 0, 0}
	// fmt.Printf("ip: %v %s", net.ParseIP(dstAddr), net.ParseIP(dstAddr))
	// data = append(data, net.ParseIP(dstAddr)...)
	// data = append(data, dstPortBytes...)

	// _, err = utils.EncryptThenWrite(conn, data, encryptor)
	_, err = eConn.Write(data)

	if err != nil {
		println(err)
		return
	}

	go io.Copy(remoteConn, eConn)
	io.Copy(eConn, remoteConn)
}
