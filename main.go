package main

import "github.com/xxxtonixxx/chatRoom/server"

func main() {
	server.Run("tcp", ":5000")
}
