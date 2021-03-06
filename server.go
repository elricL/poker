package main

import (
	"flag"
	"fmt"
	"github.com/elricL/poker/board"
	"github.com/elricL/poker/game"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var connections map[*websocket.Conn]bool

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Taken from gorilla's website
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	log.Println("Succesfully upgraded connection")
	connections[conn] = true
	if err := conn.WriteMessage(websocket.TextMessage, []byte("Go to https://github.com/elricL/poker/blob/master/README.md for usage instructions")); err != nil {
		log.Println("EERROEROERO WHile sending message")
	}
	for {
		// Blocks until a message is read
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(connections, conn)
			conn.Close()
			return
		}
		game.HandlePokerMessage(msg, &poker_board, conn)
	}
}

func main() {
	// command line flags
	poker_board = board.MakeNewBoard()
	port := flag.Int("port", 9000, "port to serve on")
	dir := flag.String("directory", "web/", "directory of web files")
	flag.Parse()

	connections = make(map[*websocket.Conn]bool)

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)
	http.HandleFunc("/ws", wsHandler)

	log.Printf("Running on port %d\n", *port)

	addr := fmt.Sprintf(":%d", *port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
