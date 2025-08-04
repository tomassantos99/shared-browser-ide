package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/tomassantos99/shared-browser-ide/internal/storage"
	"github.com/tomassantos99/shared-browser-ide/internal/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // TODO: Check actual buffer size needed
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func WsUpgrade(sessionStorage *storage.SessionStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var clientName = r.URL.Query().Get("name")
		if clientName == "" {
			errorMessage := "Missing client name"
			logrus.Error(errorMessage)
			http.Error(w, errorMessage, http.StatusNotFound)
		}
		
		sessionId := chi.URLParam(r, "id")
		idString, err := uuid.Parse(sessionId)
		
		if err != nil {
			errorMessage := "Missing session ID"
			logrus.Error(errorMessage)
			http.Error(w, errorMessage, http.StatusNotFound)
		}

		var clientSession, ok = sessionStorage.GetSession(idString)

		if !ok {
			errorMessage := fmt.Sprintf("Session with ID %s not found", sessionId)
			logrus.Error(errorMessage)
			http.Error(w, errorMessage, http.StatusNotFound)
			return
		}
		
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logrus.Error(err)
			return
		}

		var client *ws.Client = &ws.Client{
			Name:       clientName,
			Connection: conn,
			Send:       make(chan []byte, 256),
			Session:    clientSession,
		}
		clientSession.Register <- client

		go client.ReadPump()
		go client.WritePump()

	}
}
