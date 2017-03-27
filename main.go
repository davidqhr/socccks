package main

import socks5 "github.com/davidqhr/sock5/socks5"

func main() {
	socks5.Serve("localhost:8080")
}
