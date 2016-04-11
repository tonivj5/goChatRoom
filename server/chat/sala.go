package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/net/websocket"
)

type sala struct {
	ids      int
	mensajes []*Mensaje
	clientes map[int]*Cliente
	chDel    chan *Cliente
	chAdd    chan *Cliente
	chAll    chan *Mensaje
}

// NewSala crea una nueva sala de chat
func NewSala() *sala {
	return &sala{
		ids:      0,
		mensajes: make([]*Mensaje, 0, 150),
		clientes: make(map[int]*Cliente, 0),
		chAdd:    make(chan *Cliente),
		chDel:    make(chan *Cliente),
		chAll:    make(chan *Mensaje),
	}
}

func (ws *sala) Listen() {
	var buffer [1024]byte

	conectar := func(conn *websocket.Conn) {
		// Aseguramos la desconexión del cliente
		defer desconectar(conn)
		// Pedimos el apodo al cliente o si ya estaba inscrito nos mandará el id
		n, _ := conn.Read(buffer[:])
		// Creamos el cliente
		newCliente := NewCliente(conn, ws)
		// Decodificamos el JSON y rellenamos el cliente recién creado
		json.Unmarshal(buffer[:n], newCliente)
		fmt.Printf("El cliente %s realiza una nueva conexión desde: %s\n", newCliente.Apodo, conn.Request().RemoteAddr)
		if newCliente.ID == -1 {
			newCliente.ID = ws.newID()
			// Devolvemos su ID
			conn.Write([]byte(strconv.Itoa(newCliente.ID)))
		}

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
			ws.clientes[newCliente.ID] = newCliente
		// Mensaje de un cliente al todos
		case msg := <-ws.chAll:
			// Almacenamos el mensaje una única vez
			ws.mensajes = append(ws.mensajes, msg)
			// Log sobre quién y qué mensaje ha mandado
			fmt.Printf("El cliente %s ha mandado el mensaje: %s\n", msg.Cliente.Apodo, msg.Contenido)
			// Enviamos el mensaje a todo cristo
			for i := range ws.clientes {
				ws.clientes[i].WriteToCliente(msg)
			}
			/*case id := <-ws.chDel:
			delete(ws.clientes, id)*/
		}
	}
}

// Úlimos  mensajes que el nuevo cliente no ha llegado a recibir
func (ws *sala) ultimosMensajes(c *Cliente) {
	for i := range ws.mensajes {
		msg := ws.mensajes[i]
		json, err := CodificarJSON(msg)

		if err != nil {
			fmt.Println("Ocurrió un error en la codificiación del JSON")
		}

		c.conn.Write(json)
	}
}

/*func (ws *sala) clientesConectados(c *Cliente) {
	for i := range ws.clientes {

	}
}*/

func (ws *sala) nuevoClienteConectado(c *Cliente) {

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
func (ws *sala) broadcastClientes(msg *Mensaje) {
	ws.chAll <- msg
}

// Log de desconexión de un cliente
func desconectar(c *websocket.Conn) {
	fmt.Printf("El cliente %s se ha desconectado\n", c.Request().RemoteAddr)
}

func (ws *sala) newID() (id int) {
	id = ws.ids
	ws.ids++

	return
}

func CodificarJSON(msg *Mensaje) ([]byte, error) {
	return json.Marshal(msg)
}
