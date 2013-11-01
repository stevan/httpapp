package sessions

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/stevan/httpapp"
)

var validSessionID = regexp.MustCompile(`^[0-9a-f]{40}$`)

func NewCookieState(key string) *SessionCookieState {
	state := &SessionCookieState{}
	state.SessionKey = key
	return state
}

type SessionCookieState struct {
	SessionKey string
	Path       string
	Domain     string
	Secure     bool
	HttpOnly   bool
	Expires    time.Duration
}

func (state *SessionCookieState) GetSessionKey() string {
	return state.SessionKey
}

func (state *SessionCookieState) GetSessionId(e *httpapp.Env) (string, error) {
	cookie, err := e.Request.Cookie(state.GetSessionKey())
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (state *SessionCookieState) ValidateSessionId(id string) bool {
	return validSessionID.MatchString(id)
}

func (state *SessionCookieState) ExpireSession(session *Session, resp *httpapp.Response) {
	session.Expire()
	cookie := state.createCookie(session)
	resp.Headers.Add("Set-Cookie", cookie.String())
}

func (state *SessionCookieState) Extract(e *httpapp.Env) (string, error) {
	id, err := state.GetSessionId(e)
	if err != nil {
		return "", err
	}
	if !state.ValidateSessionId(id) {
		return "", errors.New("Invalid session id")
	}
	return id, nil
}

func (state *SessionCookieState) Generate() string {
	hash := sha1.New()

	t := time.Now().UnixNano()
	r := rand.New(rand.NewSource(t)).Int()

	io.WriteString(hash, fmt.Sprintf("%d", t)+fmt.Sprintf("%d", r))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (state *SessionCookieState) Finalize(session *Session, resp *httpapp.Response) {
	cookie := state.createCookie(session)
	resp.Headers.Add("Set-Cookie", cookie.String())
}

// ...

func (state *SessionCookieState) createCookie(session *Session) *http.Cookie {
	cookie := http.Cookie{}
	cookie.Name = state.SessionKey
	cookie.Value = session.Id
	cookie.Path = state.Path
	cookie.Domain = state.Domain
	cookie.Secure = state.Secure
	cookie.HttpOnly = state.HttpOnly

	if session.IsExpired() {
		cookie.Expires = time.Now()
	} else {
		if state.Expires != 0 {
			cookie.Expires = time.Now().Add(state.Expires)
		}
	}

	return &cookie
}
