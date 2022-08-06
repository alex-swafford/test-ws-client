// client.go
package main

import (
	"bufio"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

var done chan interface{}
var interrupt chan os.Signal

var reader = bufio.NewReader(os.Stdin)

func readKeys(input chan rune, interrupt chan os.Signal) {
	for {
		select {
		case <-interrupt:
			break
		default:
			char, _, err := reader.ReadRune()
			if err != nil {
				log.Fatal(err)
			}
			input <- char
		}
	}
}

func receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		log.Printf("Received: %s\n", msg)
	}
}

func main() {

	host := os.Getenv("ws_host_to_test")

	if host == "" {
		host = "localhost:8080"
	}

	done = make(chan interface{}) // Channel to indicate that the receiverHandler is done

	socketUrl := "ws://" + host + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()
	go receiveHandler(conn)

	// Our main loop for the client
	// We send our relevant packets here
	for {
		inputLine, _, _ := reader.ReadLine()
		err := conn.WriteMessage(websocket.TextMessage, []byte(string(inputLine)))
		if err != nil {
			log.Fatal(err)
		}
	}
}
