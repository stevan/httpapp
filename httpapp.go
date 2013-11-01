package httpapp

import (
	"bytes"
	"net/http"
)

// App Interface

type App interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Call(*Env) *Response
}

// Request

func NewEnv(req *http.Request) *Env {
	stash := make(map[string]interface{})
	return &Env{req, stash}
}

type Env struct {
	Request *http.Request
	stash   map[string]interface{}
}

func (e *Env) Get(key string) interface{} {
	return e.stash[key]
}

func (e *Env) Set(key string, value interface{}) {
	e.stash[key] = value
}

// Response

func NewResponse(status int) *Response {
	resp := new(Response)
	resp.Status = status
	resp.Headers = make(http.Header)
	resp.Body = new(bytes.Buffer)
	return resp
}

type Response struct {
	Status  int
	Headers http.Header
	Body    *bytes.Buffer
}

func (r *Response) WriteTo(w http.ResponseWriter) {
	out := r.Header()
	in := w.Header()
	for k, v := range out {
		in[k] = v
	}
	w.WriteHeader(r.Status)
	r.Body.WriteTo(w)
}

// ResponseWriter compatible API
// -----------------------------
// this allows you to pass in a httpapp.Response
// to places which require a http.ResponseWriter
// such as many of the http utility functions
// like http.ServeFile, etc.

func (r *Response) Header() http.Header {
	return r.Headers
}

func (r *Response) Write(buf []byte) (int, error) {
	return r.Body.Write(buf)
}

func (r *Response) WriteHeader(status int) {
	r.Status = status
}
