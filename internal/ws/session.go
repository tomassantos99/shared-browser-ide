package ws

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tomassantos99/shared-browser-ide/pkg"
)

type Session struct {
	Mut          sync.RWMutex
	Id           uuid.UUID
	clients      map[*Client]bool
	codeUpdate   chan Message
	Register     chan *Client
	unregister   chan *Client
	Close        chan bool
	Password     string
	onEmpty      func(sessionId uuid.UUID)
	editorState  EditorState
	LastActive   time.Time
	sessionAdmin *Client
}

type EditorState struct {
	programmingLanguage string
	content             string
}

const PASSWORD_LENGTH int = 5

func NewSession(onEmpty func(sessionId uuid.UUID)) *Session {
	return &Session{
		Id:           uuid.New(),
		codeUpdate:   make(chan Message),
		Register:     make(chan *Client),
		unregister:   make(chan *Client),
		Close:        make(chan bool),
		clients:      make(map[*Client]bool),
		Password:     pkg.RandomString(PASSWORD_LENGTH),
		LastActive:   time.Now(),
		editorState:  EditorState{},
		onEmpty:      onEmpty,
		sessionAdmin: nil,
	}
}

func (s *Session) Run() {
	for {
		select {
		case client := <-s.Register:
			s.registerClient(client)
			s.sendSessionClientsUpdate()
		case client := <-s.unregister:
			s.unregisterClient(client)
			s.sendSessionClientsUpdate()
		case message := <-s.codeUpdate:
			s.updateEditorState(message.ProgrammingLanguage, message.EditorContent, message.Sender)
			s.sendEditorStateUpdate(message)
		case _ = <-s.Close:
			s.closeClientConnections()
			s.onEmpty(s.Id)
			return
		}

		if len(s.clients) == 0 {
			logrus.Info(fmt.Sprintf("Session with id %s is empty. Deleting session...", s.Id.String()))
			s.onEmpty(s.Id)
			return
		}
	}
}

func (s *Session) registerClient(client *Client) {
	s.Mut.Lock()
	defer s.Mut.Unlock()

	s.clients[client] = true
	if s.sessionAdmin == nil {
		s.sessionAdmin = client
	}

	logrus.Info(fmt.Sprintf("Registered a new client for session %s: %s", s.Id.String(), client.Name))

	if s.editorState.programmingLanguage != "" { // for now it serves to check if it is the first connection to the session
		var message, err = CreateMessage(SessionCodeUpdate, s.editorState.programmingLanguage, s.editorState.content, nil, nil)
		if err != nil {
			logrus.Error(err)
		} else {
			client.Send <- message
		}
	}
}

func (s *Session) unregisterClient(client *Client) {
	s.Mut.Lock()
	defer s.Mut.Unlock()

	_, ok := s.clients[client]
	if ok {
		if s.sessionAdmin == client {
			s.sessionAdmin = nil
			for k := range s.clients {
				s.sessionAdmin = k
			}
		}
		delete(s.clients, client)
		close(client.Send)
		logrus.Info(fmt.Sprintf("Unregistered a client for session %s: %s", s.Id.String(), client.Name))
	}
}

func (s *Session) sendSessionClientsUpdate() {
	s.Mut.RLock()
	defer s.Mut.RUnlock()

	var clients []string
	for client := range s.clients {
		clients = append(clients, client.Name)
	}

	var message, err = CreateMessage(ClientsUpdate, "", "", clients, nil)
	if err != nil {
		logrus.Error(err)
		return
	}
	for client := range s.clients {
		client.Send <- message
	}
}

func (s *Session) sendEditorStateUpdate(message Message) {
	s.Mut.RLock()
	clients := make([]*Client, 0, len(s.clients))
	for client := range s.clients {
		clients = append(clients, client)
	}
	s.Mut.RUnlock()

	var updateMessage, err = CreateMessage(SessionCodeUpdate, message.ProgrammingLanguage, message.EditorContent, nil, message.Sender)
	if err != nil {
		logrus.Error(err)
		return
	}

	for _, client := range clients {
		select {
		case client.Send <- updateMessage:
		default:
			s.unregisterClient(client)
		}
	}
}

func (s *Session) updateEditorState(programmingLanguage string, content string, client *Client) {
	s.Mut.Lock()
	defer s.Mut.Unlock()

	if s.sessionAdmin == client {
		s.editorState.programmingLanguage = programmingLanguage
	}
	s.editorState.content = content
	s.LastActive = time.Now()
}

func (s *Session) closeClientConnections() {
	s.Mut.RLock()
	defer s.Mut.RLock()

	for client := range s.clients {
		close(client.Send)
	}
}
