package main

import (
	"io"
	"net/http"

	"code.google.com/p/go.net/websocket"
)

func EchoServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func main() {
	http.Handle("/echo", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic("ListenAndServer: " + err.Error())
	}
}
