package component

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/stevan/httpapp"
)

type URLMapComponent struct {
	Mux *http.ServeMux
}

func (c *URLMapComponent) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Call(httpapp.NewEnv(r)).WriteTo(w)
}

func (c *URLMapComponent) Call(e *httpapp.Env) *httpapp.Response {
	handler, _ := c.Mux.Handler(e.Request)
	if app, ok := handler.(httpapp.App); ok {
		return app.Call(e)
	} else {
		resp := httpapp.NewResponse(http.StatusNotFound)
		resp.Body.WriteString(
			fmt.Sprintf("Unable to get handler for %s, got %v instead", e.Request.URL.Path, handler),
		)
		return resp
	}
}

func (c *URLMapComponent) AddApplication(pattern string, a httpapp.App) {
	c.Mux.Handle(pattern, &prefixStripper{a, pattern})
}

// ...

type prefixStripper struct {
	app    httpapp.App
	prefix string
}

func (c *prefixStripper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Call(httpapp.NewEnv(r)).WriteTo(w)
}

func (c *prefixStripper) Call(e *httpapp.Env) *httpapp.Response {
	if p := strings.TrimPrefix(e.Request.URL.Path, c.prefix); len(p) < len(e.Request.URL.Path) {
		e.Request.URL.Path = p
		return c.app.Call(e)
	} else {
		resp := httpapp.NewResponse(http.StatusNotFound)
		resp.Body.WriteString(fmt.Sprintf("Page not found: %s", e.Request.URL.Path))
		return resp
	}
}
