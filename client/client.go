package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

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
	ip          = flag.String("ip", "127.0.0.1", "server IP")
	connections = flag.Int("conn", 1, "number of websocket connections")
	userName    = flag.String("username", "vic", "username")
)

type Client struct {
	conn     *websocket.Conn
	identity string
}

func stert() {
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: *ip + ":8000", Path: "/"}
	log.Printf("Connecting to %s", u.String())

	client := Client{}
	var err error
	client.identity = *userName

	for {
		client.conn, _, err = websocket.DefaultDialer.Dial(u.String(), http.Header{"X-Small-Chat-Id": {*userName}})
		if err != nil {
			log.Println("Failed to connect", err)
			log.Println("Reconnect in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	defer func() {
		client.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
		time.Sleep(time.Second)
		client.conn.Close()
	}()

	log.Printf("Finished initializing connection: %s is connecting to server", client.identity)

	go func() {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("Fail to receive msg ", err.Error())
		}
		if msg != nil {
			log.Printf("msg: %s", string(msg))
		}
	}()

	for {
		log.Println("Conn sending message")
		arr := make([]string, 0)

		if err := client.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second*5)); err != nil {
			for {
				fmt.Printf("Failed to receive pong: %v\n", err)
				client.conn, _, err = websocket.DefaultDialer.Dial(u.String(), http.Header{"X-Small-Chat": {*userName}})
				if err != nil {
					log.Println("Failed to connect", err)
					log.Println("Reconnect in 5 seconds...")
					time.Sleep(5 * time.Second)
					continue
				}

				if len(arr) != 0 {
					log.Println("Begin resend all unsent msgs")
					for _, msg := range arr {
						client.conn.WriteMessage(websocket.TextMessage, []byte(msg))
					}
					arr = arr[:0]
					break
				}
			}
		}

		//scanner := bufio.NewScanner(os.Stdin)

		fmt.Println("Enter text: ")
		//scanner.Scan()
		text := "test" //scanner.Text()
		if len(text) != 0 {
			client.conn.WriteMessage(websocket.TextMessage, []byte(text))
		}

		fmt.Println("Msg sended to server: ", arr)
	}
}

func main() {
	goopt.Author = "Dmitry Fofanov"
	goopt.Version = "Client 0.1"
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

	stert()

}
