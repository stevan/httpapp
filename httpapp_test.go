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
    path := e.Request.URL.Path
    switch path {
        case "/simple":
            resp := NewResponse(200)
            resp.Header.Add("Content-Length", "7")
            resp.Body.Write([]byte("HELLO WORLD"))
            return resp
        case "/empty":
            return NewResponse(204)
    }
    return nil
}

func TestSimpleApp (t *testing.T) {
    app := new(TestApp)

    req, err := http.NewRequest("GET", "http://example.com/simple", nil)
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

    if w.Header().Get("Content-Length") != "7" {
        t.Errorf(
            "got wrong header value - got: %v expected: %v",
            w.Header().Get("Content-Length"),
            "7",
        )
    }
}

func TestEmptyApp (t *testing.T) {
    app := new(TestApp)

    req, err := http.NewRequest("GET", "http://example.com/empty", nil)
    if err != nil {
        t.Logf("got an error: %s", err)
        t.FailNow()
    }

    w := httptest.NewRecorder()
    app.ServeHTTP(w, req)


    if w.Code != 204 {
        t.Errorf("got wrong status - got: %v expected: %v", w.Code, 204)
    }

    if w.Body.String() != "" {
        t.Errorf("expected empty body - got: (%v)", w.Body.String())
    }
}


