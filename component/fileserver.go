package component

import (
	"net/http"
	"path"
	"strings"

	"github.com/stevan/httpapp"
)

type FileServerComponent struct {
	Root string
}

func (c *FileServerComponent) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Call(httpapp.NewEnv(r)).WriteTo(w)
}

func (c *FileServerComponent) Call(e *httpapp.Env) *httpapp.Response {
	upath := e.Request.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		e.Request.URL.Path = upath
	}

	resp := httpapp.NewResponse(200)
	http.ServeFile(resp, e.Request, path.Join(c.Root, path.Clean(upath)))
	return resp
}
