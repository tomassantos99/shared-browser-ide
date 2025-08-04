package ws

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tomassantos99/shared-browser-ide/pkg"
)

type Session struct {
	Id         uuid.UUID
	clients    map[*Client]bool
	broadcast  chan []byte
	Register   chan *Client
	unregister chan *Client
	Password   string
	onEmpty    func(sessionId uuid.UUID)
	LastActive time.Time
}

const PASSWORD_LENGTH int = 5

func NewSession(onEmpty func(sessionId uuid.UUID)) *Session {
	return &Session{
		Id:         uuid.New(),
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		Password:   pkg.RandomString(PASSWORD_LENGTH),
		LastActive: time.Now(),
		onEmpty:    onEmpty,
	}
}

func (s *Session) Run() {
	for {
		select {
		case client := <-s.Register:
			s.registerClient(client)

		case client := <-s.unregister:
			s.unregisterClient(client)
		case message := <-s.broadcast:
			s.broadcastMessage(message)
		}

		if len(s.clients) == 0 {
			logrus.Info(fmt.Sprintf("Session with id %s is empty. Deleting session...", s.Id.String()))
			s.onEmpty(s.Id)
			return
		}
	}
}

func (s *Session) registerClient(client *Client) {
	s.clients[client] = true
	logrus.Info(fmt.Sprintf("Registered a new client for session %s: %s", s.Id.String(), client.Name))
}

func (s *Session) unregisterClient(client *Client) {
	_, ok := s.clients[client]
	if ok {
		delete(s.clients, client)
		close(client.Send)
		logrus.Info(fmt.Sprintf("Unregistered a client for session %s: %s", s.Id.String(), client.Name))
	}
}

func (s *Session) broadcastMessage(message []byte) {
	logrus.Info(fmt.Sprintf("Broadcasting message to %d clients", len(s.clients)))
	for client := range s.clients {
		select {
		case client.Send <- message:
		default:
			s.unregisterClient(client)
		}
	}
}
