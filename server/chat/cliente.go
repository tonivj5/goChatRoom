package chat

import (
	"fmt"
	"io"

	ws "golang.org/x/net/websocket"
)

type cliente struct {
	id    int
	apodo string
	conn  *ws.Conn
	sala  *sala
	ch    chan *mensaje
}

// NewCliente crea un cliente y devuelve un puntero del mismo
func NewCliente(apodo string, conn *ws.Conn, sala *sala) *cliente {
	return &cliente{
		id:    sala.newID(),
		apodo: apodo,
		conn:  conn,
		sala:  sala,
		ch:    make(chan *mensaje),
	}
}

// Enviamos un mensaje al cliente
func (c *cliente) WriteToCliente(msg *mensaje) {
	c.conn.Write([]byte(msg.cliente.apodo + ": " + msg.contenido))
}

// Leemos todo lo que envíe el cliente y lo mandamos al websocket para que haga la difusión
func (c *cliente) ReadFromCliente() {
	for {
		var buffer [10024]byte
		n, err := c.conn.Read(buffer[:])
		if err != nil {
			fmt.Printf("Cliente %s desconectado ", c.apodo)
			if err == io.EOF {
				fmt.Println("del chat")
				return
			}
			fmt.Printf("con un error %v\n", err)
			return
		}

		c.sala.broadcastClientes(&mensaje{cliente: c, contenido: string(buffer[:n])})
	}
}
