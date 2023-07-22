package itsy

import (
	"net/http"

	"go.uber.org/zap"
)

type Itsy struct {
	router *Router
	Logger *zap.Logger
}

// New creates a new Itsy instance.
func New() *Itsy {
	return &Itsy{
		router: NewRouter(),
	}
}

// GET registers a handler for the GET HTTP method.
func (i *Itsy) GET(route string, handler HandlerFunc) {
	i.router.Handle("GET", route, handler)
}

// POST registers a handler for the POST HTTP method.
func (i *Itsy) POST(route string, handler HandlerFunc) {
	i.router.Handle("POST", route, handler)
}

// Run starts the HTTP server.
func (i *Itsy) Run(addr string) {
	http.ListenAndServe(addr, i.router)
}