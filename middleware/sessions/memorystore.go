package sessions

import (
	"errors"
)

func NewMemoryStore() *SessionMemoryStore {
	stash := make(map[string]*Session)
	return &SessionMemoryStore{stash}
}

type SessionMemoryStore struct {
	stash map[string]*Session
}

func (s *SessionMemoryStore) Fetch(id string) (*Session, error) {
	session, ok := s.stash[id]
	if !ok {
		return nil, errors.New("No Session Found")
	}
	return session, nil
}

func (s *SessionMemoryStore) Store(session *Session) {
	s.stash[session.Id] = session
}

func (s *SessionMemoryStore) Remove(id string) {
	delete(s.stash, id)
}
