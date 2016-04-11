package chat

import "encoding/json"

type tipo string

const (
	ID           tipo = "id"
	MENSAJE      tipo = "msg"
	CONECTADO    tipo = "conectado"
	DESCONECTADO tipo = "desconectado"
	UPDATEID     tipo = "updateID"
)

type IMensaje interface {
	getJSON() ([]byte, error)
	getTipo() tipo
}

type MensajeBasico struct {
	Tipo  tipo   `json:"tipo"`
	Fecha string `json:"fecha"`
}

func (m *MensajeBasico) getJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MensajeBasico) getTipo() tipo {
	return m.Tipo
}

type MensajeUserEvent struct {
	*MensajeBasico
	Clientes []*Cliente `json:"clientes"`
}

func (m *MensajeUserEvent) getJSON() ([]byte, error) {
	return json.Marshal(m)
}

type Mensaje struct {
	*MensajeBasico
	Contenido string   `json:"msg"`
	Cliente   *Cliente `json:"cliente"`
}

func (m *Mensaje) getJSON() ([]byte, error) {
	return json.Marshal(m)
}

type MensajeID struct {
	*MensajeBasico
	ID int `json:"id"`
}

func (m *MensajeID) getJSON() ([]byte, error) {
	return json.Marshal(m)
}
