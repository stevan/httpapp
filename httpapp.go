package httpapp

import(
    "net/http"
    "bytes"
)

// App Interface

type App interface {
    ServeHTTP (w http.ResponseWriter, r *http.Request)
    Call (*Env) *Response
}

// Request

func NewEnv (req *http.Request) *Env {
    stash := make(map[string]interface{})
    return &Env{req, stash}
}

type Env struct {
    Request *http.Request
    stash   map[string]interface{}
}

func (e *Env) Get (key string) interface{} {
    return e.stash[key]
}

func (e *Env) Set (key string, value interface{}) {
    e.stash[key] = value
}

// Response

func NewResponse (status int) *Response {
    resp := new(Response)
    resp.Status  = status
    resp.Headers = make(http.Header)
    resp.Body    = new(bytes.Buffer)
    return resp
}

type Response struct {
    Status  int
    Headers http.Header
    Body    *bytes.Buffer
}

func (r *Response) WriteTo (w http.ResponseWriter) {
    out := map[string][]string(r.Headers)
    in  := map[string][]string(w.Header())
    for k, v := range out { in[k] = v }
    w.WriteHeader(r.Status)
    r.Body.WriteTo(w)
}

