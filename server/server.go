package server

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type mensaje struct {
	msg   string
	conn  net.Conn
	apodo string
}

var (
	mensajes []*mensaje
	server   = make(chan *mensaje)
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

		go handlerConn(conn)
	}
}

func handlerConn(conn net.Conn) {
	defer desconectar(conn)
	var buffer [10024]byte
	n, err := conn.Read(buffer[:])
	checkError(err)
	headers := strings.Split(string(buffer[:n]), "\r\n")
	fmt.Println(headers)
	query, _ := url.QueryUnescape(strings.Split(headers[0], " ")[1])
	query = query[1:]

	var nuevo = true
	params := strings.Split(query, "&")
	fmt.Println("Parámetros:", params, len(params))
	if params[0] != "" && !strings.Contains(params[0], "favicon.ico") {
		params[0] = params[0][1:]
		fmt.Printf("%#v\n", params)
		nuevo = false
	}

	headerResponse := strings.Split(headers[0], " ")[2] + " 200 OK\r\nX-Name: Soy Toni\r\nContent-Type: text/html; charset=utf8\r\n\r\n"
	conn.Write([]byte(headerResponse))

	if nuevo {
		contenido, _ := os.Open("./resources/html/chat.html")
		n, _ = contenido.Read(buffer[:])
		conn.Write(buffer[:n])
	} else {
		if len(params) == 1 {
			outOfDateMessage, err := strconv.Atoi(strings.Split(params[0], "=")[1])
			if err != nil {
				fmt.Print("Error en outOfDateMessage: ")
				checkError(err)

				return
			}

			upOfDateMessage := len(mensajes)
			lastMessge := strconv.Itoa(upOfDateMessage) + ".---*"

			fmt.Println(outOfDateMessage, upOfDateMessage)
			if outOfDateMessage > upOfDateMessage {
				conn.Write([]byte(lastMessge))

				return
			}

			listaMensajes := mensajes[outOfDateMessage:upOfDateMessage]
			mensajesAEnviar := make([]string, len(listaMensajes))
			for i := range listaMensajes {
				mensajesAEnviar[i] = listaMensajes[i].apodo + ": " + listaMensajes[i].msg
			}

			conn.Write([]byte(lastMessge + strings.Join(mensajesAEnviar, "***___***")))
		} else {
			apodo := strings.Split(params[0], "=")[1]
			msg := strings.Split(params[1], "=")[1]
			server <- &mensaje{msg: strings.TrimSpace(string(msg)), conn: conn, apodo: apodo}
		}
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
		mensajes = append(mensajes, mensaje)
		fmt.Println("El cliente", mensaje.conn.RemoteAddr(), "("+mensaje.apodo+")", "ha enviado el mensaje:", mensaje.msg)
	}
}
