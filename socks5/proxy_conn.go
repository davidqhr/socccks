package socks5

import "net"

func proxyTcp(src net.Conn, dst net.Conn) {
	buf := bufferPool.Get()
	defer bufferPool.Put(buf)

	for {
		n, err := src.Read(buf)

		if err != nil {
			break
		}

		_, err = dst.Write(buf[:n])

		if err != nil {
			break
		}
	}
}
