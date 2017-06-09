package utils

import (
	"crypto/aes"
	"encoding/binary"
	"io"
	"net"
)

// EncryptedConn is a sort of encrypted connection between socccks-client and socccks-server
type EncryptedConn struct {
	Conn      net.Conn
	Encryptor *Encryptor
}

// encrypted data protocol (bytes)
// header:
//   length 2
// body:
//   iv 32
//   encryptedData ...

func NewEncryptedConn(conn net.Conn, password string) *EncryptedConn {
	return &EncryptedConn{
		Conn:      conn,
		Encryptor: NewEncryptor(password),
	}
}

// encrypt plainText, then write them to the socket
func (ec *EncryptedConn) Write(rawData []byte) (nw int, err error) {
	encryptor := ec.Encryptor

	writeBuf := BufPool.Get()
	defer BufPool.Put(writeBuf)

	encryptBytesLength := encryptor.CFBEncrypter(rawData, writeBuf[2:])

	binary.BigEndian.PutUint16(writeBuf[:2], uint16(encryptBytesLength))
	encryptBytesLength += 2

	// fmt.Printf("[debug] raw: %v, encrypted: %v\n", rawData, writeBuf[:encryptBytesLength])

	nw, ew := ec.Conn.Write(writeBuf[:encryptBytesLength])

	if ew != nil {
		err = ew
		return
	}

	if encryptBytesLength != nw {
		err = io.ErrShortWrite
		return
	}

	// In other place, I use io.Copy to proxy data between conns,
	// io.Copy requires the length of data received from src conn to be as same as the length of data written to the dst conn.
	// So it's necessary to return "expected data" length instead of the "real data length"
	nw = encryptBytesLength - 2 - aes.BlockSize
	return
}

// read Encrypted data, fill buf with plainText
func (ec *EncryptedConn) Read(buf []byte) (rn int, err error) {
	encryptor := ec.Encryptor

	readBuffer := BufPool.Get()
	defer BufPool.Put(readBuffer)

	if _, er := io.ReadFull(ec.Conn, readBuffer[:2]); er != nil {
		err = er

		if er == io.EOF {
			return
		}

		println("read encrytped data length error")
		return
	}

	dataLen := binary.BigEndian.Uint16(readBuffer[:2])

	_, er := io.ReadFull(ec.Conn, readBuffer[:dataLen])
	if er != nil {
		err = er

		if er == io.EOF {
			return
		}
		println("can't read full encrytped data")
		return
	}

	// fmt.Printf("[debug] encrypted: %v, ", readBuffer[:dataLen])
	rn = encryptor.CFBDecrypter(readBuffer[:dataLen], buf)
	// fmt.Printf("raw: %v\n", buf[:rn])

	return
}
