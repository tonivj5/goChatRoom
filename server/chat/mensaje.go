package chat

type Mensaje struct {
	Contenido string   `json:"msg"`
	Cliente   *Cliente `json:"cliente"`
	Fecha     string   `json:"fecha"`
}
