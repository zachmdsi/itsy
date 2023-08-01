package itsy

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Itsy is the main framework instance.
	Itsy struct {
		router    *router
		resources map[string]Resource

		Logger *zap.Logger
	}
	// Middleware is a function that wraps a handler.
	Middleware func(Context, HandlerFunc) HandlerFunc
)

// Define HTTP Methods
const (
	GET    = http.MethodGet
	POST   = http.MethodPost
	PUT    = http.MethodPut
	PATCH  = http.MethodPatch
	DELETE = http.MethodDelete
)

// DefaultPort is the default port to listen on.
const DefaultPort = ":8080"

// Define HTTP Status Codes
const (
	StatusOK                  = http.StatusOK
	StatusBadRequest          = http.StatusBadRequest
	StatusUnauthorized        = http.StatusUnauthorized
	StatusForbidden           = http.StatusForbidden
	StatusNotFound            = http.StatusNotFound
	StatusMethodNotAllowed    = http.StatusMethodNotAllowed
	StatusInternalServerError = http.StatusInternalServerError
)

// Define a map of HTTP status codes to error messages.
var httpErrors = map[int]string{
	StatusOK:                  "OK",
	StatusBadRequest:          "Bad Request",
	StatusUnauthorized:        "Unauthorized",
	StatusForbidden:           "Forbidden",
	StatusNotFound:            "Not Found",
	StatusMethodNotAllowed:    "Method Not Allowed",
	StatusInternalServerError: "Internal Server Error",
}

// setupLogger creates a new logger instance.
func setupLogger() *zap.Logger {
	// Create a new encoder config.
	encoderConfig := zap.NewProductionEncoderConfig()

	// Configure the logger to output ISO8601 time and capital level names.
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Write to stdout.
	writeSyncer := zapcore.Lock(os.Stdout)

	// Create a new core.
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	// Create a new logger.
	logger := zap.New(core)

	return logger
}

// HTTPError writes an error message to the response given a status code
func HTTPError(statusCode int, res http.ResponseWriter) {
	// Get the status text from the map. If not found, default to "Internal Server Error".
	statusText, ok := httpErrors[statusCode]
	if !ok {
		statusText = httpErrors[StatusInternalServerError]
	}

	// Write the status code and message to the response.
	res.WriteHeader(statusCode)
	res.Write([]byte(strconv.Itoa(statusCode) + " " + statusText))
}

// New creates a new Itsy instance.
func New() *Itsy {
	i := &Itsy{
		resources: make(map[string]Resource),
		Logger:    setupLogger(),
	}
	i.router = newRouter(i)
	return i
}

// ServeHTTP is the entry point for all HTTP requests.
func (i *Itsy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	i.Logger.Info("Request received", zap.String("method", req.Method), zap.String("path", req.URL.Path))

	// Split the path into segments
	segments := splitPath(req.URL.Path)

	// Start at the root node
	n := i.router.index

	// Create a new base context
	c := &baseContext{
		req:    req,
		res:    NewResponse(res, i),
		resource: i.Resource(req.URL.Path),
		params: make(map[string]string),
		path:   req.URL.Path,
		itsy:   i,
	}

	// For each segment in the path
	for _, segment := range segments {
		// If the segment is not empty
		if segment != "" {
			// If the child node exists
			if child, ok := n.children[segment]; ok {
				// Move to the child node
				n = child
			} else { // Otherwise try to find a child node with a parameter
				found := false

				// For each child node
				for key, child := range n.children {
					// If the child node has a parameter
					if strings.HasPrefix(key, ":") {
						// Add the parameter to the context
						c.params[key[1:]] = segment

						// Mark the parameter as found
						found = true

						// Move to the child node
						n = child
						break
					}
				}

				if !found {
					// If no child node was found, return a 404 error
					HTTPError(http.StatusNotFound, res)
					return
				}
			}
		}
	}

	// Add the parameters to the context
	for param, value := range c.params {
		c.SetParam(param, value)
	}

	// If the node has a resource
	if n.resource != nil {
		i.Logger.Info("Resource found", zap.String("path", n.path))

		handler := n.resource.Handler(req.Method)

		// Check for a handler for the request method
		if handler != nil {
			handler(c)
		} else {
			HTTPError(StatusMethodNotAllowed, res)
			return
		}
	} else {
		// If the node does not have a resource, return a 404 error
		HTTPError(StatusNotFound, res)
		return
	}
}

// Register registers a resource to the Itsy instance.
func (i *Itsy) Register(path string, resource Resource) Resource {
	customResource := newCustomResource()

	// Register the custom resources with the router
	i.resources[path] = customResource
	i.router.addRoute(path, customResource)
	return customResource 
}

// Resource returns a resource given a path.
func (i *Itsy) Resource(path string) Resource {
	return i.resources[path]
}

// Run runs the Itsy instance.
func (i *Itsy) Run(port ...string) {
	i.Logger.Info("Starting server...")

	if len(port) > 0 {
		i.Logger.Info("Listening on port " + port[0])
		i.Logger.Fatal("Server stopped", zap.Error(http.ListenAndServe(port[0], i)))
	} else {
		i.Logger.Info("Listening on port " + DefaultPort)
		i.Logger.Fatal("Server stopped", zap.Error(http.ListenAndServe(DefaultPort, i)))
	}
}
