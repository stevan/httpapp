package httpapp

import (
    "testing"

    "net/http"
    "net/http/httptest"
)

type TestApp struct {}

func (app *TestApp) ServeHTTP (w http.ResponseWriter, r *http.Request) {
    app.Call(NewEnv(r)).WriteTo(w)
}

func (app *TestApp) Call (e *Env) *Response {
    resp := NewResponse(200)
    resp.Body.Write([]byte("HELLO WORLD"))
    return resp
}

func TestSimpleApp (t *testing.T) {
    app := new(TestApp)

    req, err := http.NewRequest("GET", "http://example.com/", nil)
    if err != nil {
        t.Logf("got an error: %s", err)
        t.FailNow()
    }

    w := httptest.NewRecorder()
    app.ServeHTTP(w, req)


    if w.Code != 200 {
        t.Errorf("got wrong status - got: %v expected: %v", w.Code, 200)
    }

    if w.Body.String() != "HELLO WORLD" {
        t.Errorf("got wrong body - got: %v expected: %v", w.Body.String(), "HELLO WORLD")
    }
}

