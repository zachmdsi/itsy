package itsy

import "net/http"

type Router struct {
	handlers map[string]HandlerFunc
}

type HandlerFunc func(*Context)

// NewRouter creates a new router instance.
func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]HandlerFunc),
	}
}

// Handle registers a handler for a given HTTP method and path.
func (r *Router) Handle(method string, route string, handler HandlerFunc) {
	key := method + "-" + route
	r.handlers[key] = handler
}

// ServeHTTP is the entry point for all HTTP requests.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := r.handlers[key]; ok {
		ctx := &Context{Request: req, ResponseWriter: w}
		handler(ctx)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Not Found"))
	}
}