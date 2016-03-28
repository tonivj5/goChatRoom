package main

import (
	"fmt"
	"net"
	"strings"
)

type datos struct {
	n   int
	err error
}

type mensaje struct {
	msg  string
	conn net.Conn
}

var (
	conexiones []net.Conn
	server     = make(chan *mensaje)
)

func main() {
	go escuchador()
	listener, err := net.Listen("tcp", ":5000")
	checkError(err)
	fmt.Println("Se está escuchando en", listener.Addr())

	for {
		conn, err := listener.Accept()
		checkError(err)

		fmt.Println("El cliente", conn.RemoteAddr(), "se ha conectado")

		conexiones = append(conexiones, conn)

		go handlerConn(conn)
	}
}

func handlerConn(conn net.Conn) {
	defer desconectar(conn)
	var buffer [512]byte

	for {
		n, err := conn.Read(buffer[:])

		if err != nil {
			break
		}
		server <- &mensaje{msg: strings.TrimSpace(string(buffer[:n])), conn: conn}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func desconectar(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Este cliente", conn.RemoteAddr(), "se ha desconectado.")
}

func escuchador() {
	for {
		dato := <-server

		for i := range conexiones {
			if dato.conn != conexiones[i] {
				conexiones[i].Write([]byte(dato.msg + "\n"))
			} else {
				fmt.Println("REPE!!")
			}
		}

		fmt.Println("El cliente", dato.conn.RemoteAddr(), "ha enviado el mensaje:", dato.msg)
	}
}
