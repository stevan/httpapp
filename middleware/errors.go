package middleware

import (
    "os"
    "net/http"

    "github.com/stevan/httpapp"
)

type ErrorHandler struct {
    App httpapp.App
}

func (m *ErrorHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
    m.Call(httpapp.NewEnv(r)).WriteTo(w)
}

func (m *ErrorHandler) Call (e *httpapp.Env) (resp *httpapp.Response) {
    defer func() {
        if err := recover(); err != nil {
            var status int
            if os.IsNotExist(err.(error)) {
                status = http.StatusNotFound
            } else {
                status = http.StatusInternalServerError
            }
            resp = httpapp.NewResponse(status)
            resp.Body.WriteString(err.(error).Error())
        }
    }()
    return m.App.Call(e)
}