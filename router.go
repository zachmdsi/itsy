package itsy

import "net/http"

type Router struct {
	// handlers maps a HTTP method + path to a handler function.
	handlers map[string]HandlerFunc

	// itsy is a reference to the Itsy application.
	itsy *Itsy
}

type HandlerFunc func(Context)

// NewRouter creates a new router instance.
func NewRouter(itsy *Itsy) *Router {
	return &Router{
		handlers: make(map[string]HandlerFunc),
		itsy: itsy,
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
		ctx := r.itsy.newBaseContext(req, w) 
		handler(ctx)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Not Found"))
	}
}
