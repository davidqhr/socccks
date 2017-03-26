package client

import (
	"math/rand"
	"net"
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
	return &Client{Conn: conn, Id: randStringRunes(32)}
}
