package main

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	userID := 1
	groupID := "11111111-1111-1111-1111-111111111111"
	token := "dummytoken"

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/group/ws", RawQuery: "user_id=1&token=" + token}
	log.Printf("Connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	// send a test message
	message := `{
		"sender_id": 1,
		"group_id": "` + groupID + `",
		"content": "Hello group from Go client!"
	}`
	conn.WriteMessage(websocket.TextMessage, []byte(message))

	// listen for messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s", msg)
		time.Sleep(2 * time.Second)
	}
}
