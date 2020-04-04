package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

// Message struct for messaging
type Message struct {
	Text string `json:"text"`
}

func ShowLogo() {
	fmt.Println("")
	fmt.Println("")
	fmt.Println("██████╗ ███████╗██████╗  █████╗  ██████╗████████╗       ██████╗██╗  ██╗ █████╗ ████████╗")
	fmt.Println("██╔══██╗██╔════╝██╔══██╗██╔══██╗██╔════╝╚══██╔══╝      ██╔════╝██║  ██║██╔══██╗╚══██╔══╝")
	fmt.Println("██████╔╝█████╗  ██║  ██║███████║██║        ██║   █████╗██║     ███████║███████║   ██║   ")
	fmt.Println("██╔══██╗██╔══╝  ██║  ██║██╔══██║██║        ██║   ╚════╝██║     ██╔══██║██╔══██║   ██║   ")
	fmt.Println("██║  ██║███████╗██████╔╝██║  ██║╚██████╗   ██║         ╚██████╗██║  ██║██║  ██║   ██║   ")
	fmt.Println("╚═╝  ╚═╝╚══════╝╚═════╝ ╚═╝  ╚═╝ ╚═════╝   ╚═╝          ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   ")
	fmt.Println("Made with ❤ by @iob_j")
}

var (
	port = flag.String("port", "9000", "ws connection port")
)

func connect() (*websocket.Conn, error) {
	return websocket.Dial(fmt.Sprintf("ws://localhost:%s", *port), "", mockedIP())
}

func mockedIP() string {
	var arr [4]int
	for i := 0; i < 4; i++ {
		rand.Seed(time.Now().UnixNano())
		arr[i] = rand.Int()
	}
	return fmt.Sprintf("http://%d.%d.%d.%d", arr[0], arr[1], arr[2], arr[3])
}

func main() {
	flag.Parse()
	ShowLogo()
	ws, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	var m Message
	go func() {
		for {
			err := websocket.JSON.Receive(ws, &m)
			if err != nil {
				fmt.Println("Can't obtain message: ", err.Error())
				break
			}
			fmt.Println("Redacted User: ", m)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		m := Message{
			Text: text,
		}
		err = websocket.JSON.Send(ws, m)
		if err != nil {
			fmt.Println("Can't send message: ", err.Error())
			break
		}
	}
}
