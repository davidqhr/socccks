package client

import (
	"fmt"
	"io"
	"net"

	"github.com/davidqhr/socccks/utils"
)

type Client struct {
	Conn       net.Conn
	AuthMethod byte
	Password   string
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func NewClient(conn net.Conn, password string) *Client {
	return &Client{
		Conn:     conn,
		Password: password,
	}
}

func (client *Client) GetSupportAuthMethods() (methods []byte, err error) {
	conn := client.Conn

	buf := utils.BufPool.Get()
	defer utils.BufPool.Put(buf)

	_, er := io.ReadFull(conn, buf[:2])

	if er != nil {
		err = er
		return
	}

	version := buf[0]

	if version != utils.Version {
		err = fmt.Errorf("DO NOT SUPPORT PROXY Version %X", version)
		return
	}

	methodsCount := int(buf[1])

	_, er = io.ReadFull(conn, buf[:methodsCount])

	if er != nil {
		err = er
		return
	}

	methods = buf[:methodsCount]
	return
}

func (client *Client) SetAuthMethod(method byte) error {
	client.AuthMethod = method
	_, err := client.Conn.Write([]byte{utils.Version, method})
	return err
}

func (client *Client) AuthSuccess() error {
	_, err := client.Conn.Write([]byte("\x01\x00"))
	return err
}
