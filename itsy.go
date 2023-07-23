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
		Logger: SetupLogger(),
	}
}

// GET registers a handler for a GET request.
func (i *Itsy) GET(route string, handler HandlerFunc) {
	i.router.Handle("GET", route, handler)
}

// POST registers a handler for a POST request.
func (i *Itsy) POST(route string, handler HandlerFunc) {
	i.router.Handle("POST", route, handler)
}

// PUT registers a handler for a PUT request.
func (i *Itsy) PUT(route string, handler HandlerFunc) {
	i.router.Handle("PUT", route, handler)
}

// DELETE registers a handler for a DELETE request.
func (i *Itsy) DELETE(route string, handler HandlerFunc) {
	i.router.Handle("DELETE", route, handler)
}

// PATCH registers a handler for a PATCH request.
func (i *Itsy) PATCH(route string, handler HandlerFunc) {
	i.router.Handle("PATCH", route, handler)
}

// Run starts the HTTP server. 
func (i *Itsy) Run(addr string) {
	http.ListenAndServe(addr, i.router)
}