package itsy

import (
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type (
	// Itsy is the main framework instance.
	Itsy struct {
		router      *Router      // router is the main router tree
		middlewares []Middleware // middlewares is a list of middleware

		Logger *zap.Logger // Logger is a zap logger
	}
	Middleware  func(HandlerFunc) HandlerFunc // Middleware is a function that wraps a handler.
	HandlerFunc func(Context) error           // HandlerFunc is a function that handles a request.
)

// New creates a new Itsy instance.
func New() *Itsy {
	itsy := &Itsy{
		Logger: SetupLogger(),
	}
	itsy.router = NewRouter(itsy)
	return itsy
}

// Define a map of HTTP status codes to error messages.
var httpErrors = map[int]string{
	http.StatusOK:                  "OK",
	http.StatusBadRequest:          "Bad Request",
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusForbidden:           "Forbidden",
	http.StatusNotFound:            "Not Found",
	http.StatusMethodNotAllowed:    "Method Not Allowed",
	http.StatusInternalServerError: "Internal Server Error",
}

// HTTPError writes an error message to the response given a status code
func HTTPError(statusCode int, w http.ResponseWriter) {
	// Get the status text from the map. If not found, default to "Internal Server Error".
	statusText, ok := httpErrors[statusCode]
	if !ok {
		statusCode = http.StatusInternalServerError
		statusText = httpErrors[statusCode]
	}

	// Write the status code and message to the response.
	w.WriteHeader(statusCode)
	w.Write([]byte(strconv.Itoa(statusCode) + " " + statusText))
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
func (i *Itsy) Use(m Middleware) {
	i.middlewares = append(i.middlewares, m)
}

// Run starts the HTTP server.
func (i *Itsy) Run(addr string) {
	http.ListenAndServe(addr, i.router)
}
