package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	clientIdentities map[string]*websocket.Conn
)

type (
	Message struct {
		Recipient string
		Text      string
	}
)

func newUserMap() {
	clientIdentities = map[string]*websocket.Conn{}
}

func ws(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := r.Header.Get("X-Small-Chat-Id")
	clientIdentities[client] = conn

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}

		dec := gob.NewDecoder(bytes.NewBuffer(msg))
		var req Message
		if dec.Decode(&req) != nil {
			log.Fatal("decode:", err)
		}

		log.Printf("msg from %s -> to %s: Message %s", client, req.Recipient, req.Text)

		res := Message{req.Recipient, req.Text}
		var buffer bytes.Buffer

		enc := gob.NewEncoder(&buffer)
		if enc.Encode(res) != nil {
			log.Fatal("encode:", err)
		}

		wsc, ok := clientIdentities[res.Recipient]
		if ok {
			wsc.WriteMessage(websocket.TextMessage, buffer.Bytes())
		} else {
			if strings.ToLower(req.Recipient) == "all" {
				for client, wsc := range clientIdentities {
					if client != "" {
						wsc.WriteMessage(websocket.TextMessage, buffer.Bytes())
					}
				}
			}
		}
	}
}

func main() {
	newUserMap()
	http.HandleFunc("/", ws)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
