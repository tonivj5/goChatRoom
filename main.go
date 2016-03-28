package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type conexion struct {
	conn *net.Conn
	msg  string
}

func main() {
	l, err := net.Listen("tcp", ":5000")
	defer l.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error al escuchar en el puerto:", err)

		os.Exit(1)
	}

	var connections []net.Conn
	c2 := make(chan conexion)
	go func() {
		for {
			connection := <-c2
			for _, connect := range connections {
				if connection.conn != &connect {
					connect.Write([]byte(connection.msg))
				} else {
					fmt.Println("Hola caracola")
				}
			}
		}
	}()

	fmt.Println("Socket funcionando")
	c := make(chan string)

	for {
		conn, err := l.Accept()
		fmt.Println("Cliente acceptado:", conn.RemoteAddr())

		if err != nil {
			fmt.Fprintln(os.Stderr, "Ha ocurrido un error:", err)
			continue
		}
		connections = append(connections, conn)
		go handler(conn, c)
		go func() {
			for {
				c2 <- conexion{
					msg:  (conn.RemoteAddr().String() + " escribió: " + strings.TrimSpace(<-c) + "\n"),
					conn: &conn,
				}
			}
		}()
	}
}

func handler(c net.Conn, channel chan string) {
	defer c.Close()

	var buf [512]byte
	c.Write([]byte("Bienvenido al chat!\n"))

	for {
		n, err := c.Read(buf[:])
		if err != nil {
			break
		}
		channel <- string(buf[0:n])
	}
}
