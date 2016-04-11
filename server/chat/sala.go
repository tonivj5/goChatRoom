package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

type sala struct {
	ids      int
	mensajes []*Mensaje
	clientes map[int]*Cliente
	chDel    chan int
	chAdd    chan *Cliente
	chAll    chan IMensaje
}

// NewSala crea una nueva sala de chat
func NewSala() *sala {
	return &sala{
		ids:      0,
		mensajes: make([]*Mensaje, 0, 150),
		clientes: make(map[int]*Cliente, 0),
		chAdd:    make(chan *Cliente),
		chDel:    make(chan int),
		chAll:    make(chan IMensaje),
	}
}

func (ws *sala) Listen() {
	var buffer [1024]byte

	conectar := func(conn *websocket.Conn) {
		// Pedimos el apodo al cliente o si ya estaba inscrito nos mandará el ID
		n, _ := conn.Read(buffer[:])
		// Creamos el cliente
		newCliente := NewCliente(conn, ws)
		// Aseguramos la desconexión del cliente
		defer ws.desconectar(newCliente)
		// Decodificamos el JSON y rellenamos el cliente recién creado
		json.Unmarshal(buffer[:n], newCliente)
		fmt.Printf("El cliente %s realiza una nueva conexión desde: %s\n", newCliente.Apodo, conn.Request().RemoteAddr)
		// Si el cliente ya tiene un ID o esa ID ya está cogida, le damos una nueva
		if _, ok := ws.clientes[newCliente.ID]; ok || newCliente.ID == -1 {
			// Reemplazamos su ID por una nueva
			id := ws.newID()
			for _, ok := ws.clientes[id]; ok; id = ws.newID() {
				_, ok = ws.clientes[id]
			}
			fmt.Printf("Su nuevo ID es %d\n", id)
			newCliente.ID = id

			json, err := json.Marshal(&MensajeID{MensajeBasico: &MensajeBasico{Tipo: ID}, ID: id})

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error al generar el json de nuevo ID: %v\n", err)
			}

			conn.Write(json)
		}

		// Actualizamos su estado para el resto de clientes
		ws.actualizarEstadoCliente(newCliente, CONECTADO)
		// Añadimos el cliente
		ws.addCliente(newCliente)
		// Mandamos al cliente mensajes antiguos
		ws.ultimosMensajes(newCliente)
		// Mandamos los clientes que se encuentran conectados en este momento
		ws.todosClientesConectados(newCliente)
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
			ws.clientes[newCliente.ID] = newCliente
		// Mensaje de un cliente a todos
		case imsg := <-ws.chAll:
			if imsg.getTipo() == MENSAJE {
				msg, _ := imsg.(*Mensaje)
				// Almacenamos el mensaje una única vez
				ws.mensajes = append(ws.mensajes, msg)
				// Log sobre quién y qué mensaje ha mandado
				fmt.Printf("El cliente %s ha mandado el mensaje: %s\n", msg.Cliente.Apodo, msg.Contenido)
			}

			// Enviamos el mensaje a todo cristo
			for i := range ws.clientes {
				ws.clientes[i].WriteToCliente(imsg)
			}
		case id := <-ws.chDel:
			delete(ws.clientes, id)
		}
	}
}

// Úlimos  mensajes que el nuevo cliente no ha llegado a recibir
func (ws *sala) ultimosMensajes(c *Cliente) {
	for i := range ws.mensajes {
		msg := ws.mensajes[i]
		json, err := json.Marshal(msg)

		if err != nil {
			fmt.Println("Ocurrió un error en la codificiación del JSON")
		}

		c.conn.Write(json)
	}
}

func (ws *sala) todosClientesConectados(c *Cliente) {
	msg := &MensajeUserEvent{MensajeBasico: &MensajeBasico{Tipo: CONECTADO}, Clientes: ws.mapToSlice(ws.clientes)}

	c.WriteToCliente(msg)
}

func (ws *sala) actualizarEstadoCliente(c *Cliente, evento tipo) {
	msg := &MensajeUserEvent{MensajeBasico: &MensajeBasico{Tipo: evento}, Clientes: []*Cliente{c}}

	ws.broadcastClientes(msg)
}

// Añadir un cliente nuevo para el broadcast
func (ws *sala) addCliente(c *Cliente) {
	ws.chAdd <- c
}

func (ws *sala) delCliente(c *Cliente) {
	ws.chAdd <- c
}

// Difunde el mensaje por el resto de clientes
// TODO: Enviar al cliente además de la hora la fecha (quizá separar Fecha de hora...)
func (ws *sala) broadcastClientes(msg IMensaje) {
	ws.chAll <- msg
}

// Log de desconexión de un cliente
func (ws *sala) desconectar(c *Cliente) {
	fmt.Printf("El cliente %s con IP %s se ha desconectado\n", c.Apodo, c.conn.Request().RemoteAddr)
	ws.actualizarEstadoCliente(c, DESCONECTADO)
	ws.chDel <- c.ID
}

func (ws *sala) newID() (id int) {
	id = ws.ids
	ws.ids++

	return
}

func (ws *sala) mapToSlice(map[int]*Cliente) []*Cliente {
	i, clientes := 0, make([]*Cliente, len(ws.clientes))
	for key := range ws.clientes {
		clientes[i] = ws.clientes[key]
		i++
	}

	return clientes
}
