package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/droundy/goopt"
	"github.com/gorilla/websocket"
)

var License = `License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law`

func Version() error {
	fmt.Printf("GoPoloniexTest 0.1 %s\n\nCopyright (C) 2021 %s\n%s\n", goopt.Version, goopt.Author, License)
	os.Exit(0)
	return nil
}

func PrintUsage() {
	fmt.Fprintf(os.Stderr, goopt.Usage())
	os.Exit(1)
}

var (
	clientIdentities map[string]*websocket.Conn
)

type (
	Message struct {
		recipient string
		text      string
	}

	Response struct {
		length []byte
		data   []byte
	}
)

func newUserMap() {
	clientIdentities = map[string]*websocket.Conn{}
}

func IntToByteArray(num int, size int) []byte {
	tmp_arr := make([]byte, size)
	binary.LittleEndian.PutUint16(tmp_arr, uint16(num))

	result_arr := make([]byte, size)
	for i := size - 1; i >= 0; i-- {
		result_arr[i] = tmp_arr[(size-1)-i]
	}
	return result_arr
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

		//js, _ := json.Marshal(map[string]string{})
		//var buf bytes.Buffer
		//buf.ReadByte(&msg)
		//dec := gob.NewDecoder(buf)
		err := enc.Decode(&msg)
		if err != nil {
			log.Println("Error: %s", err)
		}

		length, err := strconv.ParseInt(string(msg[:2]), 2, 64)
		if err != nil {
			fmt.Printf("Error: %s", err)
		}

		if int(length) == len(msg[2:]) {
			req := Message{}
			json.Unmarshal(msg[2:], &req)

			log.Printf("msg from %s -> to %s: Len %v, Message %s", client, req.recipient, length, req.text)

			js, _ := json.Marshal(Message{client, req.text})
			resp := Response{IntToByteArray(len(js), 2), js}
			text := string(resp.length) + string(resp.data)
			wsc, ok := clientIdentities[req.recipient]
			if ok {
				wsc.WriteMessage(websocket.TextMessage, []byte(text))
			} else {
				if req.recipient == "all" {
					for _, wsc := range clientIdentities {
						wsc.WriteMessage(websocket.TextMessage, []byte(text))
					}
				}
			}
		}
	}
}

func main() {
	goopt.Author = "Dmitry Fofanov"
	goopt.Version = "Server 0.1"
	goopt.Summary = "Implementation of the test task, chat in the goland language (Details: https://github.com/DFofanov/goChatTest)"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s Port\n OPTION\n", os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless an option is passed to it."
	}
	goopt.NoArg([]string{"-v", "--version"}, "outputs version information and exits", Version)

	//	if err != nil {
	//		fmt.Println(err)
	//	}

	newUserMap()
	http.HandleFunc("/", ws)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
