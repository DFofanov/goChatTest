package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ip       = flag.String("ip", "127.0.0.1", "server IP")
	userName = flag.String("username", "vasa", "username")
)

type (
	Client struct {
		Conn     *websocket.Conn
		Identity string
	}
)

type (
	Message struct {
		Recipient string
		Text      string
	}
)

func main() {
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: *ip + ":8000", Path: "/"}
	log.Printf("Connecting to %s", u.String())

	client := Client{}
	var err error
	client.Identity = *userName

	for {
		client.Conn, _, err = websocket.DefaultDialer.Dial(u.String(), http.Header{"X-Small-Chat-Id": {*userName}})
		if err != nil {
			log.Println("Failed to connect", err)
			log.Println("Reconnect in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	defer func() {
		client.Conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
		time.Sleep(time.Second)
		client.Conn.Close()
	}()

	log.Printf("Finished initializing connection: %s is connecting to server", client.Identity)

	go func() {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Fail to receive msg ", err.Error())
		}
		if msg != nil {
			dec := gob.NewDecoder(bytes.NewBuffer(msg))
			var req Message
			if dec.Decode(&req) != nil {
				log.Fatal("decode:", err)
			}
			log.Printf("msg from %s: Message %s", req.Recipient, req.Text)
		}
	}()

	for {
		arr := make([]string, 0)

		if err := client.Conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second*5)); err != nil {
			for {
				fmt.Printf("Failed to receive pong: %v\n", err)
				client.Conn, _, err = websocket.DefaultDialer.Dial(u.String(), http.Header{"X-Small-Chat": {*userName}})
				if err != nil {
					log.Println("Failed to connect", err)
					log.Println("Reconnect in 5 seconds...")
					time.Sleep(5 * time.Second)
					continue
				}

				if len(arr) != 0 {
					log.Println("Begin resend all unsent msgs")
					for _, msg := range arr {
						client.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
					}
					arr = arr[:0]
					break
				}
			}
		}
		fmt.Print("Enter User: ")
		var user string
		fmt.Scanln(&user)
		if user != "" {
			fmt.Print("Enter Message: ")
			var message string
			fmt.Scanln(&message)

			msg := Message{user, message}
			var buffer bytes.Buffer
			enc := gob.NewEncoder(&buffer)
			err := enc.Encode(msg)
			if err != nil {
				log.Fatal("encode:", err)
			}

			if buffer.Len() != 0 {
				client.Conn.WriteMessage(websocket.TextMessage, buffer.Bytes())
			}

			log.Println("Conn sending message")
		}
	}
}
