package utils

import (
	"crypto/aes"
	"encoding/binary"
	"io"
	"net"
)

type EncryptedConn struct {
	Conn      net.Conn
	Encryptor *Encryptor
}

// data protocol (bytes)
// length 2, iv 32, data ...

func NewEncryptedConn(conn net.Conn, password string) *EncryptedConn {
	return &EncryptedConn{
		Conn:      conn,
		Encryptor: NewEncryptor(password),
	}
}

func (ec *EncryptedConn) Write(rawData []byte) (nw int, err error) {
	encryptor := ec.Encryptor

	writeBuf := make([]byte, 1024*65)
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

	nw = encryptBytesLength - 2 - aes.BlockSize
	return
}

func (ec *EncryptedConn) Read(buf []byte) (rn int, err error) {
	encryptor := ec.Encryptor

	readBuffer := make([]byte, 1024*65)

	if _, er := io.ReadFull(ec.Conn, readBuffer[:2]); er != nil {
		if er == io.EOF {
			return
		}

		println("read encrytped data length error")
		err = er
		return
	}

	dataLen := binary.BigEndian.Uint16(readBuffer[:2])

	_, er := io.ReadFull(ec.Conn, readBuffer[:dataLen])
	if er != nil {
		if er == io.EOF {
			return
		}
		println("can't read full encrytped data")
		err = er
		return
	}

	// fmt.Printf("[debug] encrypted: %v, ", readBuffer[:dataLen])
	rn = encryptor.CFBDecrypter(readBuffer[:dataLen], buf)
	// fmt.Printf("raw: %v\n", buf[:rn])

	return
}
