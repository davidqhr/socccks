package helper

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func BytesToIpv4String(ip_bytes []byte) string {
	var addr bytes.Buffer

	for i := 0; i < len(ip_bytes); i++ {
		if i != 0 {
			addr.WriteString(".")
		}
		addr.WriteString(fmt.Sprintf("%d", int(ip_bytes[i])))
	}

	return addr.String()
}

func Ipv4StringToBytes(ip string) []byte {
	parts := strings.Split(ip, ".")
	bytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		partInt, _ := strconv.Atoi(parts[i])
		bytes[i] = byte(partInt)
	}
	return bytes
}
