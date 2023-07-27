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
		router      *router             // router is the main router tree
		resources   map[string]Resource // map of path to Resource

		Logger *zap.Logger // Logger is a zap logger
	}
	HandlerFunc func(Context) error // HandlerFunc is a function that handles a request.
)

// New creates a new Itsy instance.
func New() *Itsy {
	itsy := &Itsy{
		Logger: setupLogger(),
		resources: make(map[string]Resource),
	}
	itsy.router = newRouter(itsy)
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

// setupLogger creates a new zap logger instance.
func setupLogger() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	writeSyncer := zapcore.Lock(os.Stdout)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core)
	return logger
}

// ServeHTTP is the entry point for all HTTP requests.
func (i *Itsy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	segments := strings.FieldsFunc(req.URL.Path, func(r rune) bool { return r == '/' })
	currentNode := i.router.index
	params := make(map[string]string)
	for _, segment := range segments {
		// If a direct match is found, move to the next node
		if child, ok := currentNode.children[segment]; ok {
			currentNode = child
		} else {
			// If a direct match isn't found, try to match a parameterized route
			found := false
			for key, child := range currentNode.children {
				if strings.HasPrefix(key, ":") {
					params[key[1:]] = segment
					currentNode = child
					found = true
					break
				}
			}

			// If no parameterized route is found, return a 404.
			if !found {
				HTTPError(http.StatusNotFound, w)
				return
			}
		}
	}

	// Fetch handler function based on method
	handler, ok := currentNode.handlers[req.Method]
	if !ok {
		HTTPError(http.StatusMethodNotAllowed, w)
		return
	}

	// Create a new context
	ctx := i.newBaseContext(req, w)
	// Call the handler function
	err := handler(ctx)
	if err != nil {
		HTTPError(http.StatusInternalServerError, w)
		return
	}	
}

// AddResource adds a resource to the Itsy instance.
func (i *Itsy) AddResource(method, path string, resource Resource) {
	// Add resource to the resources map.
	i.resources[path] = resource
	// Add a handler function for the resource to the router.
	handler := func(ctx Context) error {
		return i.handleResource(ctx, resource)
	}
	i.router.addRoute(http.MethodGet, path, handler)
}

// handleResource handles a resource based on the HTTP method.
func (i *Itsy) handleResource(ctx Context, resource Resource) error {
	switch ctx.Request().Method {
	case http.MethodGet:
		html := resource.Render()
		ctx.ResponseWriter().Write([]byte(html))
	case http.MethodPost:
		// logic for POST requests
	case http.MethodPut:
		// logic for PUT requests
	case http.MethodDelete:
		// logic for DELETE requests
	default:
		// logic for unsupported methods 
	}
	return nil
}
// Run starts the HTTP server.
func (i *Itsy) Run(addr string) {
	http.ListenAndServe(addr, i)
}
