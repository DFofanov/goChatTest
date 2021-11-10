package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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
		Length    []byte
		Recipient []byte
		Text      []byte
	}
)

func ws(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	clientName := r.Header.Get("X-Small-Chat-Id")
	clientIdentities[clientName] = conn

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}
		message := Message{msg[:2], []byte(strings.Split(string(msg[2:]), "->")[0]), []byte(strings.Split(string(msg[2:]), "->")[1])}
		log.Printf("msg from %s -> to %s: Len %s, Message %s", string(clientName), string(message.Recipient), string(message.Length), string(message.Text))

		// switch clientName {
		// case "vic":
		// 	go clientIdentities["judy"].WriteMessage(websocket.TextMessage, message)
		// case "judy":
		// 	go clientIdentities["vic"].WriteMessage(websocket.TextMessage, message)
		// }

	}
}

func newUserMap() {
	clientIdentities = map[string]*websocket.Conn{}
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
	//	goopt.Parse(nil)
	//	if len(goopt.Args) != 1 {
	//		PrintUsage()
	//	}

	//	bs := []byte(strconv.Itoa(31415926))
	//	fmt.Println(bs)

	newUserMap()
	http.HandleFunc("/", ws)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
