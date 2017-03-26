package client

import (
	"errors"
	"fmt"
	"math/rand"
	"net"

	"github.com/davidqhr/sock5/helper"
)

type Client struct {
	Id         string
	Conn       net.Conn
	AuthMethod byte
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

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

	if version != helper.VERSION {
		return emptyBytes, errors.New(fmt.Sprintf("DO NOT SUPPORT PROXY VERSION %X", version))
	}

	methods_count := int(buf[1])
	methods := buf[2 : methods_count+2]

	return methods, nil
}

func (client *Client) SetAuthMethod(method byte) error {
	client.AuthMethod = method
	_, err := client.Conn.Write([]byte{helper.VERSION, method})
	return err
}

func (client *Client) AuthSuccess() error {
	_, err := client.Conn.Write([]byte("\x01\x00"))
	return err
}
