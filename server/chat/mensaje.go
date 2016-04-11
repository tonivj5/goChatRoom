package chat

type tipo string

const (
	ID          tipo = "id"
	MENSAJE     tipo = "msg"
	UPDATEUSERS tipo = "updateUsers"
	UPDATEID    tipo = "updateID"
)

type MensajeBasico struct {
	Tipo  tipo   `json:"tipo"`
	Fecha string `json:"fecha"`
}

type MensajeConectados struct {
	MensajeBasico
}

type Mensaje struct {
	MensajeBasico
	Contenido string   `json:"msg"`
	Cliente   *Cliente `json:"cliente"`
}
