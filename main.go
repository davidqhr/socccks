package main

import socks5 "github.com/davidqhr/socks5/src"

func main() {
	socks5.Serve("localhost:8080")
}
