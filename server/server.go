package server

import (
	"fmt"
	"net"
	"os"
)

func Run() {

	l, err := net.Listen("tcp", ":5001")

	if err != nil {
		fmt.Fprintln(os.Stderr, "Puerto en uso")
		os.Exit(1)
	}
	fmt.Println("Socket funcionando")
	for {
		conn, err := l.Accept()

		if err != nil {
			continue
		}

		handler(conn)
	}
}

func handler(c net.Conn) {
	var buf []byte

	n, _ := c.Read(buf)

	c.Write(buf[:n])

	c.Close()

	return
}
