package middleware

import (
    "log"
    "net/http"
    "github.com/stevan/httpapp"
)

type SimpleLoggingHandler struct {
    App httpapp.App
}

func (m *SimpleLoggingHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
    m.Call(httpapp.NewEnv(r)).WriteTo(w)
}

func (m *SimpleLoggingHandler) Call (e *httpapp.Env) (resp *httpapp.Response) {
    log.Printf("%s %s %s", e.Request.RemoteAddr, e.Request.Method, e.Request.URL)
    return m.App.Call(e)
}
