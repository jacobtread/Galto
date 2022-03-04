package main

import (
	"Galto/net"
	"embed"
)

//go:embed assets
var assets embed.FS

func main() {
	net.TestEncode()
	net.TestDecode()
}
