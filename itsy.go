package itsy

import (
	"net/http"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Itsy is the main framework instance.
	Itsy struct {
		router    *router             // Used to route requests to resources.
		resources map[string]Resource // A map of resource names to resources.

		Logger *zap.Logger // Uses zap for logging.
	}
	// HandlerFunc is a function that handles a request.
	HandlerFunc func(Context) error
	// Middleware is a function that wraps a handler.
	Middleware func(Context, HandlerFunc) HandlerFunc
)

const (
	// DefaultPort is the default port to listen on.
	DefaultPort = ":8080"

	// Define HTTP Methods
	GET    = http.MethodGet
	POST   = http.MethodPost
	PUT    = http.MethodPut
	PATCH  = http.MethodPatch
	DELETE = http.MethodDelete

	// Define HTTP Status Codes
	StatusOK                  = http.StatusOK                  // 200
	StatusBadRequest          = http.StatusBadRequest          // 400
	StatusUnauthorized        = http.StatusUnauthorized        // 401
	StatusForbidden           = http.StatusForbidden           // 403
	StatusNotFound            = http.StatusNotFound            // 404
	StatusMethodNotAllowed    = http.StatusMethodNotAllowed    // 405
	StatusInternalServerError = http.StatusInternalServerError // 500

	// Define HTTP Header Names
	HeaderAccept      = "Accept"
	HeaderContentType = "Content-Type"

	// Define MIME Types
	MIMETextHTML = "text/html"
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
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	writeSyncer := zapcore.Lock(os.Stdout)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core)
	return logger
}

// HTTPError writes an error message to the response given a status code
func HTTPError(statusCode int, message string, res http.ResponseWriter, logger *zap.Logger) {
    statusText, ok := httpErrors[statusCode]
    if !ok {
        statusText = httpErrors[StatusInternalServerError]
    }

    errorMessage := statusText + ": " + message
    logger.Error("HTTP Error", zap.Int("status", statusCode), zap.String("message", errorMessage))

    res.WriteHeader(statusCode)
    res.Write([]byte(errorMessage))
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

// ServeHTTP is the main entry point for the Itsy instance.
func (i *Itsy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	c := i.prepareRequestContext(res, req, path)
	i.Logger.Info("Request", zap.String("method", req.Method), zap.String("path", path))

	n := i.processRouteSegments(c, path)
	if n == nil {
		i.Logger.Error("No route found", zap.String("path", path))
		return
	}
	i.handleRequestNode(n, c, req, res)
}

// prepareRequestContext creates a new context for the request.
func (i *Itsy) prepareRequestContext(res http.ResponseWriter, req *http.Request, path string) Context {
	wrapper := &responseWriterWrapper{original: res}
	c := newBaseContext(req, NewResponse(wrapper.original, i), i.Resource(path), path, i)
	if c.Request().Header.Get(HeaderAccept) == "" {
		c.Request().Header.Set(HeaderContentType, MIMETextHTML)
	}
	return c
}

// processRouteSegments processes the route segments of the request path.
func (i *Itsy) processRouteSegments(c Context, path string) *node {
	segments := splitPath(path)
	n := i.router.index
	for _, segment := range segments {
		if segment != "" {
			if child, ok := n.children[segment]; ok {
				n = child
			} else {
				found := false
				for key, child := range n.children {
					if strings.HasPrefix(key, ":") {
						c.SetResource(child.resource)
						c.SetParam(key[1:], segment)
						c.Resource().SetParam(key[1:], segment)
						found = true
						n = child
						break
					}
				}
				if !found {
					HTTPError(StatusNotFound, "Resource does not exist", c.Response().Writer, i.Logger)
					return nil
				}
			}
		}
	}
	return n
}

// handleRequestNode handles the request node by calling the appropriate handler.
func (i *Itsy) handleRequestNode(n *node, c Context, req *http.Request, res http.ResponseWriter) {
	if n == nil {
		HTTPError(StatusNotFound, "Resource does not exist", res, i.Logger)
		return
	}

	if n.resource != nil {
		handler := n.resource.Handler(req.Method)
		if handler == nil {
			HTTPError(StatusMethodNotAllowed, "Handler does not exist for the request method", res, i.Logger)
			return
		}
		handlerWithHypermedia := HypermediaMiddleware(handler)
		if handlerWithHypermedia != nil {
			handlerWithHypermedia(c)
		} else {
			HTTPError(StatusMethodNotAllowed, "Handler does not exist for the request method", res, i.Logger)
		}
	} else {
		HTTPError(StatusNotFound, "Resource does not exist", res, i.Logger)
	}
}

// Register registers a resource to the Itsy instance.
func (i *Itsy) Register(path string, resource Resource) Resource {
	customResource := newCustomResource(path, i)
	i.resources[path] = customResource
	i.router.addRoute(path, customResource)
	return customResource
}

// SetResource sets a resource given a path.
func (i *Itsy) SetResource(path string, resource Resource) {
	i.resources[path] = resource
}

// Resource returns a resource given a path.
func (i *Itsy) Resource(path string) Resource {
	return i.resources[path]
}

// ResourceExists returns true if a resource exists given a path.
func (i *Itsy) ResourceExists(path string) bool {
	_, exists := i.resources[path]
	return exists
}

// Run runs the Itsy instance.
func (i *Itsy) Run(port ...string) {
	i.Logger.Info("Starting server...")

	if len(port) > 0 {
		i.Logger.Info("Listening on port " + port[0])
		err := http.ListenAndServe(port[0], i)
		if err != nil {
			i.Logger.Fatal("Server stopped", zap.Error(err))
		}
	} else {
		i.Logger.Info("Listening on port " + DefaultPort)
		err := http.ListenAndServe(DefaultPort, i)
		if err != nil {
			i.Logger.Fatal("Server stopped", zap.Error(err))
		}
	}
}
