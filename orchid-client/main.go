package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"os"
)

type Settings struct {
	ServerUrl string
}

func Usage() {
	fmt.Println(`
  The orchid-client sends commands to a running orchid-server.
  
  "orchid-client <cmd>" will run the command on the remote orchid-server
  defined in the config file in the current dirrectory.
  
  remote.json configuration file format:
  {
	  "ServerUrl": "ws://<url>"
  }
  
  The following addition flags are available:
	`)
	flag.PrintDefaults()
}

func loadSettings(path string) (Settings, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Settings{}, err
	}

	var s Settings
	err = json.Unmarshal(data, &s)
	return s, err
}

func main() {
	path := flag.String("config", "remote.json", "Where to find the configuration file")
	flag.Usage = Usage
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		fmt.Println("Missing arguments: <command>")
		return
	}

	settings, err := loadSettings((*path))
	errorHandle(err)

	conn, err := connect(settings.ServerUrl)
	errorHandle(err)

	sendCommand(args[0]+" "+args[1], conn)
	readMessages(conn)
}

func connect(url string) (*websocket.Conn, error) {
	fmt.Println("Connecting to " + url)
	conn, _, err := (&websocket.Dialer{}).Dial(url, http.Header{})
	errorHandle(err)

	return conn, nil
}

func sendCommand(command string, conn *websocket.Conn) {
	fmt.Println("Sending command: " + command)
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
		fmt.Printf("An error occured: %s\n", err.Error())
		os.Exit(1)
	}
}
