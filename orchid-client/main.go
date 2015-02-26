package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	//"os"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		fmt.Println("Missing arguments: <url> <command>")
		return
	}

	conn, err := connect(args[0])
	errorHandle(err)

	sendCommand(args[1], conn)
	readMessages(conn)
}

func connect(url string) (*websocket.Conn, error) {
	conn, _, err := (&websocket.Dialer{}).Dial(url, http.Header{})
	errorHandle(err)

	return conn, nil
}

func sendCommand(command string, conn *websocket.Conn) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(command))
	errorHandle(err)
}

func readMessages(conn *websocket.Conn) {
	_, p, err := conn.ReadMessage()
	for err == nil {
		fmt.Println(string(p))
		_, p, err = conn.ReadMessage()
	}
}

func errorHandle(err error) {
	if err != nil {
		panic(err)
	}
}
