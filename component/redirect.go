package component

import (
	"github.com/stevan/httpapp"
	"net/http"
)

type RedirectComponent struct {
	Location string
}

func (c *RedirectComponent) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Call(httpapp.NewEnv(r)).WriteTo(w)
}

func (c *RedirectComponent) Call(e *httpapp.Env) *httpapp.Response {
	resp := httpapp.NewResponse(http.StatusMovedPermanently)
	resp.Headers.Add("Location", c.Location)
	return resp
}
