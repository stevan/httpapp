package sessions

import (
    "github.com/stevan/httpapp"
)

// --------------------------------------
// Interfaces
// --------------------------------------

type SessionState interface {
    GetSessionKey () string
    GetSessionId ( *httpapp.Env ) (string, error)
    ValidateSessionId ( string ) bool
    ExpireSession ( *Session, *httpapp.Response )
    Extract ( *httpapp.Env ) (string, error)
    Generate () string
    Finalize ( *Session, *httpapp.Response )
}

type SessionStore interface {
    Fetch  ( string ) (*Session, error)
    Store  ( *Session )
    Remove ( string )
}