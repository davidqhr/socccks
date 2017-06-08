package client

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func bytesToIpv4String(ipBytes []byte) string {
	var addr bytes.Buffer

	for i := 0; i < len(ipBytes); i++ {
		if i != 0 {
			addr.WriteString(".")
		}
		addr.WriteString(fmt.Sprintf("%d", int(ipBytes[i])))
	}

	return addr.String()
}

func ipv4StringToBytes(ip string) []byte {
	parts := strings.Split(ip, ".")
	bytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		partInt, _ := strconv.Atoi(parts[i])
		bytes[i] = byte(partInt)
	}
	return bytes
}
