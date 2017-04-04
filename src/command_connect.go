package socks5

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

func handleCmdConnection(client *Client) {
	conn := client.Conn
	// conn.SetDeadline(time.Now().Add(5000))
	buf := make([]byte, 100)

	_, err := conn.Read(buf)

	if err != nil {
		return
	}

	addrType := buf[0]
	var addr string
	portBytes := make([]byte, 2)

	// logger.Debug(client, "addrType: %X", addrType)
	switch addrType {
	case AptyIPV4:
		ipv4Bytes := buf[1:5]
		portBytes = buf[5:7]
		addr = bytesToIpv4String(ipv4Bytes)
	case AptyIPV6:
		logger.Error(client, "NOT IMPLEMENTED APTY_IPV6")
		return
	case AptyDomainName:
		domainLen := uint8(buf[1])
		addr = string(buf[2 : 2+domainLen])
		portBytes = buf[2+domainLen : 2+domainLen+2]
		// logger.Error(client, "NOT IMPLEMENTED APTY_DOMAINNAME")
		// return
	}

	port := binary.BigEndian.Uint16(portBytes)

	logger.Debug("addr", addr, "port", port)

	remoteConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", string(addr), port))
	// remoteConn.SetDeadline(time.Now().Add(5000))
	defer remoteConn.Close()

	if err != nil {
		logger.Info(client, "Connect remote Failed %s", conn.LocalAddr().String())
		return
	}

	addrAndPort := strings.Split(remoteConn.LocalAddr().String(), ":")
	dstAddr := addrAndPort[0]
	dstPort := addrAndPort[1]

	dstPortBytes := make([]byte, 2)
	dstPortInt, err := strconv.Atoi(dstPort)

	if err != nil {
		logger.Error(err)
		return
	}

	binary.BigEndian.PutUint16(dstPortBytes, uint16(dstPortInt))

	data := []byte{Version, ReplySuccess, Rsv, AptyIPV4}
	data = append(data, ipv4StringToBytes(dstAddr)...)
	data = append(data, dstPortBytes...)

	logger.Debug("send data", data)
	_, err = conn.Write(data)

	if err != nil {
		logger.Error(err)
		return
	}

	go io.Copy(conn, remoteConn)
	io.Copy(remoteConn, conn)
}
