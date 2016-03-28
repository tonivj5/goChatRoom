package server

import (
	"fmt"
	"net"
	"strings"
)

type mensaje struct {
	msg   string
	conn  net.Conn
	apodo string
}

var (
	conexiones []net.Conn
	server     = make(chan *mensaje)
)

// Run arranca la sala de chat
func Run(protocolo string, direcionPuertoEscucha string) {
	go escuchador()
	listener, err := net.Listen(protocolo, direcionPuertoEscucha)
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
	conn.Write([]byte("¿Cuál es tu apodo?: "))
	n, err := conn.Read(buffer[:])
	apodo := strings.TrimSpace(string(buffer[:n]))

	for {
		conn.Write([]byte(apodo + ": "))
		n, err = conn.Read(buffer[:])

		if err != nil {
			break
		}
		server <- &mensaje{msg: strings.TrimSpace(string(buffer[:n])), conn: conn, apodo: apodo}
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
		mensaje := <-server

		for i := range conexiones {
			if mensaje.conn != conexiones[i] {
				conexiones[i].Write([]byte(mensaje.msg + "\n"))
			} else {
				fmt.Println("REPE!!")
			}
		}

		fmt.Println("El cliente", mensaje.conn.RemoteAddr(), "("+mensaje.apodo+")", "ha enviado el mensaje:", mensaje.msg)
	}
}
