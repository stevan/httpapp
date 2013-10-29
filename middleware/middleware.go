package middleware

import (
    "github.com/stevan/httpapp"
    "github.com/stevan/httpapp/middleware/sessions"

    "code.google.com/p/goauth2/oauth"
)

func HandleErrors (app httpapp.App) httpapp.App {
    return &ErrorHandler{app}
}

func HandleSessions (app httpapp.App, state sessions.SessionState, store sessions.SessionStore) httpapp.App {
    return &SessionHandler{app, state, store}
}

func HandleGoogleOAuthAuthentication (app httpapp.App, oauth_cfg *oauth.Config) httpapp.App {
    return &GoogleOAuthHandler{app, oauth_cfg}
}

func HandleSimpleLogging (app httpapp.App) httpapp.App {
    return &SimpleLoggingHandler{app}
}

