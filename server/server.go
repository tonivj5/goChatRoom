package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/xxxtonixxx/chatRoom/server/chat"
)

// Run arranca la sala de chat
func Run(direcionPuertoEscucha string) {
	http.HandleFunc("/", handlerConn)
	go chat.NewSala().Listen()

	fmt.Println("Se está escuchando en", direcionPuertoEscucha)
	http.ListenAndServe(direcionPuertoEscucha, nil)
}

func handlerConn(w http.ResponseWriter, r *http.Request) {
	var buf [4096]byte
	file, err := os.Open("./resources/html/chat.html")
	checkError(err)

	n, err := file.Read(buf[:])
	checkError(err)

	w.Write(buf[:n])
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Ha ocurrido un error:", err)
	}
}
