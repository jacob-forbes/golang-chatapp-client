package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

var scanner = bufio.NewScanner(os.Stdin)

func readMessage(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error: ", err)
			return
		}
		fmt.Println(string(msg))
	}
}

func writeMessage(conn *websocket.Conn) {
	for {
		// trying to stop the message from echoing locally...

		// 	fd := int(os.Stdin.Fd())
		// 	oldState, err := terminal.MakeRaw(fd)
		// 	if err != nil {
		// 		fmt.Println("Error setting raw terminal:", err)
		// 		return
		// 	}

		// 	var input []byte
		// 	for {
		// 		var b [1]byte
		// 		_, err := os.Stdin.Read(b[:])
		// 		if err != nil {
		// 			fmt.Println("Error reading input:", err)
		// 			return
		// 		}

		// 		// Break loop if Enter key is pressed
		// 		if b[0] == '\n' {
		// 			terminal.Restore(fd, oldState)
		// 			break
		// 		}

		// 		// Append to input buffer
		// 		input = append(input, b[0])
		// 		fmt.Print(b)

		// 	}

		// 	fmt.Print("\033[2K\033[1G") // Clear current line

		scanner.Scan()
		input := scanner.Text()
		if len(input) <= 0 {
			continue
		}

		err := conn.WriteMessage(websocket.TextMessage, []byte(input))
		if err != nil {
			log.Println("Write error: ", err)
			return
		}
	}
}

func main() {

	fmt.Print("Enter room name: ")
	scanner.Scan()
	roomName := scanner.Text()

	//add self as sender client
	url := "ws://localhost:5005/ws/" + string(roomName)
	log.Printf("Connecting to %s", url)
	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		if err == websocket.ErrBadHandshake {
			log.Printf("handshake failed with status %d", resp.StatusCode)
			return
		}
		log.Println("[DIAL]", err)
		return
	}
	go readMessage(conn)
	writeMessage(conn)
}
