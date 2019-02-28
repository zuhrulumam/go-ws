package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var updagrader = websocket.Upgrader{}

func main() {
	indexFile, err := os.Open("index.html")
	if err != nil {
		log.Printf("error on : %s", err)
	}

	index, err := ioutil.ReadAll(indexFile)
	if err != nil {
		log.Printf("error on : %s", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(index))
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := updagrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("error on : %s", err)
			return
		}

		fmt.Println("Client Subscribed")

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("error on : %s", err)
				break
			}

			fmt.Println(string(msg), msgType)
			err = conn.WriteMessage(msgType, []byte("pong"))
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		fmt.Println("Client Unsubscribed")

	})

	http.ListenAndServe(":3000", nil)
}
