package middleware

import (
	"net/http"

	"github.com/stevan/httpapp"
	"github.com/stevan/httpapp/middleware/sessions"
)

// --------------------------------------
// Session Handler
// --------------------------------------

type SessionHandler struct {
	App   httpapp.App
	State sessions.SessionState
	Store sessions.SessionStore
}

func (m *SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Call(httpapp.NewEnv(r)).WriteTo(w)
}

func (m *SessionHandler) Call(e *httpapp.Env) *httpapp.Response {
	session, err := m.GetSession(e)
	if err != nil {
		session = sessions.NewSession(m.State.Generate())
	}
	e.Set("session", session)
	resp := m.App.Call(e)
	m.Finalize(session, resp)
	return resp
}

func (m *SessionHandler) GetSession(e *httpapp.Env) (*sessions.Session, error) {
	id, err := m.State.Extract(e)
	if err != nil {
		return nil, err
	}
	session, err := m.Store.Fetch(id)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (m *SessionHandler) Finalize(session *sessions.Session, resp *httpapp.Response) {
	m.commit(session)
	if session.IsExpired() {
		m.expireSession(session, resp)
	} else {
		m.saveState(session, resp)
	}
}

// ...

func (m *SessionHandler) commit(session *sessions.Session) {
	if session.IsExpired() {
		m.Store.Remove(session.Id)
	} else {
		m.Store.Store(session)
	}
}

func (m *SessionHandler) expireSession(session *sessions.Session, resp *httpapp.Response) {
	m.State.ExpireSession(session, resp)
}

func (m *SessionHandler) saveState(session *sessions.Session, resp *httpapp.Response) {
	m.State.Finalize(session, resp)
}
