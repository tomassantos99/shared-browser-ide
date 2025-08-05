package ws

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tomassantos99/shared-browser-ide/pkg"
)

type Session struct {
	Id          uuid.UUID
	clients     map[*Client]bool
	codeUpdate  chan Message
	Register    chan *Client
	unregister  chan *Client
	Password    string
	onEmpty     func(sessionId uuid.UUID)
	editorState EditorState
	LastActive  time.Time
}

type EditorState struct {
	programmingLanguage string
	content  string
}

const PASSWORD_LENGTH int = 5

func NewSession(onEmpty func(sessionId uuid.UUID)) *Session {
	return &Session{
		Id:          uuid.New(),
		codeUpdate:  make(chan Message),
		Register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		Password:    pkg.RandomString(PASSWORD_LENGTH),
		LastActive:  time.Now(),
		editorState: EditorState{},
		onEmpty:     onEmpty,
	}
}

func (s *Session) Run() {
	for {
		select {
		case client := <-s.Register:
			s.registerClient(client)

		case client := <-s.unregister:
			s.unregisterClient(client)
		case message := <-s.codeUpdate:
			s.updateClientsEditorState(message)
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

	if (s.editorState.programmingLanguage != "") { // for now it serves to check if it is the first connection to the session
		var message, err = CreateMessage(SessionCodeUpdate, s.editorState.programmingLanguage, s.editorState.content)
		if (err != nil) {
			logrus.Error(err)
		}else{
			client.Send <- message
		}
	}
}

func (s *Session) unregisterClient(client *Client) {
	_, ok := s.clients[client]
	if ok {
		delete(s.clients, client)
		close(client.Send)
		logrus.Info(fmt.Sprintf("Unregistered a client for session %s: %s", s.Id.String(), client.Name))
	}
}

func (s *Session) updateClientsEditorState(message Message) {
	s.updateEditorState(message.ProgrammingLanguage, message.Content)
	var updateMessage, err = CreateMessage(SessionCodeUpdate, message.ProgrammingLanguage, message.Content)
	if err != nil {
		logrus.Error(err)
		return
	}
	
	logrus.Info(fmt.Sprintf("Broadcasting message to %d clients", len(s.clients)))
	for client := range s.clients {
		select {
		case client.Send <- updateMessage:
		default:
			s.unregisterClient(client)
		}
	}
}

func (s *Session) updateEditorState(programmingLanguage string, content string) {
	if (s.editorState.programmingLanguage == "") {
		s.editorState.programmingLanguage = programmingLanguage
	}
	s.editorState.content = content
}
