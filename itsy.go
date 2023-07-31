package itsy

import (
	"net/http"
	"os"
	"strconv"

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

// DefaultPort is the default port to listen on.
const DefaultPort = ":8080"

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
	// Split the path into segments
	segments := splitPath(req.URL.Path)

	// Start at the root node
	n := i.router.index

	// Create a new base context
	c := &baseContext{
		req:    req,
		res:    &res,
		params: make(map[string]string),
	}

	// For each segment in the path
	for _, segment := range segments {
		// If the segment is not empty
		if segment != "" {
			// If the child node exists
			if child, ok := n.children[segment]; ok {
				// Move to the child node
				n = child

			} else if child, ok := n.children[":"]; ok { // If the child node for parameters exists, move to it.
				// Move to the child node
				n = child

				// Store the parameter
				c.params[n.param] = segment
			} else {
				// If the child node does not exist, return a 404 error
				HTTPError(http.StatusNotFound, res)
				return
			}
		}
	}

	// If the node has a resource
	if n.resource != nil {
		// Render the resource
		res.Write([]byte(n.resource.Render(c)))
	} else {
		// If the node does not have a resource, return a 404 error
		http.Error(res, "Not Found", http.StatusNotFound)
		return
	}
}

// Register registers a resource to the Itsy instance.
func (i *Itsy) Register(path string, resource Resource) {
	i.resources[path] = resource
	i.router.addRoute(path, resource)
}

// Run runs the Itsy instance.
func (i *Itsy) Run() {
	i.Logger.Info("Starting server...")
	i.Logger.Fatal("Server stopped", zap.Error(http.ListenAndServe(DefaultPort, i)))
}
