package chat

import (
	"encoding/json"
	"fmt"
	"io"

	ws "golang.org/x/net/websocket"
)

type Cliente struct {
	ID    int    `json:"id"`
	Apodo string `json:"apodo"`
	conn  *ws.Conn
	sala  *sala
	ch    chan *Mensaje
}

// NewCliente crea un cliente y devuelve un puntero del mismo
func NewCliente(conn *ws.Conn, sala *sala) *Cliente {
	return &Cliente{
		conn: conn,
		sala: sala,
		ch:   make(chan *Mensaje),
	}
}

// Enviamos un mensaje al cliente
func (c *Cliente) WriteToCliente(msg *Mensaje) {
	json, err := CodificarJSON(msg)

	if err != nil {
		fmt.Println("Ocurrió un error en la codificiación del JSON")
	}

	c.conn.Write(json)
}

// Leemos todo lo que envíe el cliente y lo mandamos al websocket para que haga la difusión
func (c *Cliente) ReadFromCliente() {
	for {
		var buffer [10024]byte
		n, err := c.conn.Read(buffer[:])
		if err != nil {
			fmt.Printf("Cliente %s desconectado ", c.Apodo)
			if err == io.EOF {
				fmt.Println("del chat")
				return
			}
			fmt.Printf("con un error %v\n", err)
			return
		}
		msg := &Mensaje{Cliente: c}
		json.Unmarshal(buffer[:n], msg)
		c.sala.broadcastClientes(msg)
	}
}

func (c *Cliente) String() string {
	return fmt.Sprintf("ID: %d, Apodo: %s, IP/Puerto: %s", c.ID, c.Apodo, c.conn.Request().RemoteAddr)
}
