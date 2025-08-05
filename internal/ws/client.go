package ws

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Session *Session

	// The websocket Connection.
	Connection *websocket.Conn

	Send chan Message

	Name string
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func (client *Client) ReadPump() {
	defer func() {
		client.Session.unregister <- client
		client.Connection.Close()
	}()
	client.Connection.SetReadLimit(maxMessageSize)
	client.Connection.SetReadDeadline(time.Now().Add(pongWait))
	client.Connection.SetPongHandler(func(string) error { client.Connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := client.Connection.ReadMessage()
		message = bytes.TrimSpace(message)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Warn(err)
			}
			break
		}

		var convertedMessage, conversionError = convertMessage(message)
		if conversionError != nil {
			break
		}

		client.Session.broadcast <- convertedMessage
	}
}

func (client *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Connection.Close()
	}()
	for {
		select {
		case message, ok := <-client.Send:
			client.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The session closed the channel.
				client.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			var byteMessage, marshalErr = json.Marshal(message)
			if marshalErr != nil {
				logrus.Error("Error converting message from session")
			} else {
				w.Write(byteMessage)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func convertMessage(message []byte) (Message, error) {
	var convertedMessage Message

	var err = json.Unmarshal(message, &convertedMessage)

	if err != nil {
		logrus.Error("Error converting message from client. ", err)
		return DefaultMessage(), err
	}

	return convertedMessage, nil
}
