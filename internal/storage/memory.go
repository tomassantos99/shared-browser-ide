package storage

import (
	"github.com/google/uuid"
	"github.com/tomassantos99/shared-browser-ide/internal/ws"
	"sync"
)

type SessionStorage struct {
	mut      sync.RWMutex
	sessions map[uuid.UUID]*ws.Session
}

func NewSessionStorage() *SessionStorage{
	return &SessionStorage{
		sessions: make(map[uuid.UUID] *ws.Session),
	}
}

func (s *SessionStorage) SaveSession(session *ws.Session) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.sessions[session.Id] = session
}

func (s *SessionStorage) GetSession(id uuid.UUID) (*ws.Session, bool) {
	s.mut.RLock()
	defer s.mut.RUnlock()
	session, ok := s.sessions[id]
	return session, ok
}

func (s *SessionStorage) RemoveSession(id uuid.UUID) {
	s.mut.Lock()
	defer s.mut.Unlock()
	delete(s.sessions, id)
}
