package chat

import (
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/net/websocket"
)

type sala struct {
	ids      int
	mensajes []*mensaje
	clientes []*cliente
	chDel    chan *cliente
	chAdd    chan *cliente
	chAll    chan *mensaje
}

// NewSala crea una nueva sala de chat
func NewSala() *sala {
	return &sala{
		ids:      0,
		mensajes: make([]*mensaje, 0, 150),
		clientes: make([]*cliente, 0, 10),
		chAdd:    make(chan *cliente),
		chDel:    make(chan *cliente),
		chAll:    make(chan *mensaje),
	}
}

func (ws *sala) Listen() {
	var buffer [512]byte

	conectar := func(conn *websocket.Conn) {
		// Aseguramos la desconexión del cliente
		defer desconectar(conn)
		// Pedimos el apodo al cliente
		n, _ := conn.Read(buffer[:])
		apodo := string(buffer[:n])
		fmt.Printf("El cliente %s realiza una nueva conexión desde: %s\n", apodo, conn.RemoteAddr().String())

		// Creamos el cliente
		newCliente := NewCliente(apodo, conn, ws)
		// Devolvemos su ID
		conn.Write([]byte(strconv.Itoa(newCliente.id)))
		ws.addCliente(newCliente)
		// Mandamos al cliente mensajes antiguos
		ws.ultimosMensajes(newCliente)
		// Todo lo que envíe el cliente lo trataremos
		newCliente.ReadFromCliente()
	}

	// Llamará a conectar cada vez que se conecte un cliente mediante el websocket en */chat
	http.Handle("/chat", websocket.Handler(conectar))

	// Tratamos toda la información que nos llega mediante los canales
	for {
		select {
		// Nuevo cliente
		case newCliente := <-ws.chAdd:
			fmt.Println("Nuevo cliente añadido:", newCliente)
			ws.clientes = append(ws.clientes, newCliente)
		// Mensaje de un cliente al todos
		case msg := <-ws.chAll:
			// Almacenamos el mensaje una única vez
			ws.mensajes = append(ws.mensajes, msg)
			// Log sobre quién y qué mensaje ha mandado
			fmt.Printf("El cliente %s ha mandado el mensaje: %s\n", msg.cliente.apodo, msg.contenido)
			// Enviamos el mensaje a todo cristo
			for i := range ws.clientes {
				ws.clientes[i].WriteToCliente(msg)
			}
		}
	}
}

// Úlimos  mensajes que el nuevo cliente no ha llegado a recibir
func (ws *sala) ultimosMensajes(c *cliente) {
	for i := range ws.mensajes {
		msg := ws.mensajes[i]
		c.conn.Write([]byte(msg.cliente.apodo + ": " + msg.contenido))
	}
}

// Añadir un cliente nuevo para el broadcast
func (ws *sala) addCliente(c *cliente) {
	ws.chAdd <- c
}

func (ws *sala) delCliente(c *cliente) {
	ws.chAdd <- c
}

// Difunde el mensaje por el resto de clientes
func (ws *sala) broadcastClientes(msg *mensaje) {
	ws.chAll <- msg
}

// Log de desconexión de un cliente
func desconectar(c *websocket.Conn) {
	fmt.Printf("El cliente %s se ha desconectado\n", c.RemoteAddr().String())
}

func (ws *sala) newID() (id int) {
	id = ws.ids
	ws.ids++

	return
}
