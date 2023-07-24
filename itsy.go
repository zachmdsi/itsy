package itsy

import (
	"net/http"

	"go.uber.org/zap"
)

type (
	// Itsy is the main framework instance.
	Itsy struct {
		router      *Router      // router is the main router tree
		middlewares []Middleware // middlewares is a list of middleware

		Logger *zap.Logger       // Logger is a zap logger
	}
	// Router is the tree of routes
	Router struct {
		Index *Route // Index is the root node of the router tree.
		itsy  *Itsy  // itsy is a reference to the main framework instance.
	}
	// Route represents a node in a router
	Route struct {
		Path     string                 // Path is the path segment of the node.
		Handlers map[string]HandlerFunc // Handlers is a map of HTTP methods to handlers.
		Children map[string]*Route      // Children is a map of path segments to child nodes.
		IsParam  bool                   // IsParam is true if the path segment is a parameter.
	}
	Middleware func(HandlerFunc) HandlerFunc
	HandlerFunc func(Context) error
)

// New creates a new Itsy instance.
func New() *Itsy {
	itsy := &Itsy{
		Logger: SetupLogger(),
	}
	itsy.router = NewRouter(itsy)
	return itsy
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

// Use adds a new middleware
func(i *Itsy) Use(m Middleware) {
	i.middlewares = append(i.middlewares, m)
}

// Run starts the HTTP server.
func (i *Itsy) Run(addr string) {
	http.ListenAndServe(addr, i.router)
}