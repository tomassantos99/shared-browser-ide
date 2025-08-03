package handlers

import (
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024, // TODO: Check actual buffer size needed
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func WsUpgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if (err != nil){
		logrus.Error(err)
	}
	defer conn.Close()

	for{ // Test connection
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			logrus.Println("Read error:", err)
			break
		}

		logrus.Printf("Received: %+v\n", msg)
	}

	// TODO: create client and goroutines

}
