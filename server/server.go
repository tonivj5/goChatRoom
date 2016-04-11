package server

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/xxxtonixxx/chatRoom/server/chat"
)

const (
	htmlPATH = "./resources/html"
	jsPATH   = "./resources/js"
)

// Run arranca la sala de chat
func Run(direcionPuertoEscucha string) {
	http.HandleFunc("/", handlerConn)
	go chat.NewSala().Listen()

	fmt.Println("Se está escuchando en", direcionPuertoEscucha)
	http.ListenAndServe(direcionPuertoEscucha, nil)
}

func handlerConn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		handlerHTML(w, r, htmlPATH+"chat.html")

		return
	}

	regex := regexp.MustCompile(".*.html$")
	path := htmlPATH + r.URL.Path
	if regex.MatchString(path) {
		handlerHTML(w, r, path)

		return
	}

	path = jsPATH + r.URL.Path
	regex = regexp.MustCompile(".*.js$")

	if regex.MatchString(path) {
		handlerJS(w, r, path)

		return
	}
}

func handlerHTML(w http.ResponseWriter, r *http.Request, path string) {
	var buf [10024]byte
	if !exists(path) {
		fmt.Println("No existe el HTML", path)
		path = htmlPATH + "/chat.html"
	}

	file, err := os.Open(path)
	checkError(err)

	n, err := file.Read(buf[:])
	checkError(err)

	w.Write(buf[:n])
}

func handlerJS(w http.ResponseWriter, r *http.Request, path string) {
	var buf [10024]byte
	if !exists(path) {
		fmt.Println("No existe el JS", path)
		path = jsPATH + "/chat.js"
	}

	file, err := os.Open(path)
	checkError(err)

	n, err := file.Read(buf[:])
	checkError(err)
	w.Header().Set("Content-Type; charset=utf-8", "text/javascript")

	w.Write(buf[:n])
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Ha ocurrido un error:", err)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return true
}
