package utils

import (
	"io"
	"net"
	"time"
)

// this code is copy from https://golang.org/src/io/io.go
// add timeout and buffer pool support

func Copy(dst net.Conn, src net.Conn) (written int64, err error) {
	timeoutDuration := 15 * time.Second
	buf := Pool32K.Get()
	defer Pool32K.Put(buf)

	for {
		src.SetReadDeadline(time.Now().Add(timeoutDuration))
		nr, er := src.Read(buf)
		if nr > 0 {

			dst.SetWriteDeadline(time.Now().Add(timeoutDuration))
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
