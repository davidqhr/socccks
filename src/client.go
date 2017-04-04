package socks5

import (
	"errors"
	"fmt"
	"net"
)

type Client struct {
	Id         string
	Conn       net.Conn
	AuthMethod byte
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func NewClient(conn net.Conn) *Client {
	return &Client{
		Conn: conn,
		Id:   randStringRunes(32),
	}
}

func (client *Client) GetSupportAuthMethods() ([]byte, error) {
	conn := client.Conn
	var buf = make([]byte, 100)
	var emptyBytes = make([]byte, 0)

	_, err := conn.Read(buf)

	if err != nil {
		return emptyBytes, err
	}

	version := buf[0]

	if version != Version {
		return emptyBytes, errors.New(fmt.Sprintf("DO NOT SUPPORT PROXY Version %X", version))
	}

	methodsCount := int(buf[1])
	methods := buf[2 : methodsCount+2]

	return methods, nil
}

func (client *Client) SetAuthMethod(method byte) error {
	client.AuthMethod = method
	_, err := client.Conn.Write([]byte{Version, method})
	return err
}

func (client *Client) AuthSuccess() error {
	_, err := client.Conn.Write([]byte("\x01\x00"))
	return err
}
