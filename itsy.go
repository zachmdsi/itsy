package itsy

import (
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Itsy is the main framework instance.
	Itsy struct {
		router    *router             // Used to route requests to resources.
		resources map[string]Resource // A map of resource names to resources.

		Logger *zap.Logger  // Uses zap for logging.
		Server *http.Server // The HTTP server.
	}
	// HandlerFunc is a function that handles a request.
	HandlerFunc func(Context) error
	// Middleware is a function that wraps a handler.
	Middleware func(Context, HandlerFunc) HandlerFunc
)

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

// sendHTTPError writes an error message to the response given a status code
func (i *Itsy) sendHTTPError(statusCode int, message string, res http.ResponseWriter, logger *zap.Logger) {
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
	currentNode := i.router.index

	for _, segment := range segments {
		if segment != "" {
			found := false
			for _, child := range currentNode.children {
				if child.path == segment || (child.regex != nil && child.regex.MatchString(segment)) {
					if child.regex != nil && child.regex.MatchString(segment) {
						c.SetParam(child.param, segment)
						c.SetResource(child.resource)
						child.resource.SetParam(child.param, segment)
					}
					currentNode = child
					found = true
					break
				}
			}
			if !found {
				i.sendHTTPError(StatusNotFound, "Resource does not exist", c.Response().Writer, i.Logger)
				return nil
			}
		}
	}
	return currentNode
}

// handleRequestNode handles the request node by calling the appropriate handler.
func (i *Itsy) handleRequestNode(n *node, c Context, req *http.Request, res http.ResponseWriter) {
	if n == nil {
		i.sendHTTPError(StatusNotFound, "Resource does not exist", res, i.Logger)
		return
	}

	if n.resource != nil {
		handler := n.resource.Handler(req.Method)
		if handler == nil {
			i.sendHTTPError(StatusMethodNotAllowed, "Handler does not exist for the request method", res, i.Logger)
			return
		}
		handlerWithHypermedia := HypermediaMiddleware(handler)
		if handlerWithHypermedia != nil {
			handlerWithHypermedia(c)
		} else {
			i.sendHTTPError(StatusMethodNotAllowed, "Handler does not exist for the request method", res, i.Logger)
		}
	} else {
		i.sendHTTPError(StatusNotFound, "Resource does not exist", res, i.Logger)
	}
}

// Register registers a resource to the Itsy instance.
func (i *Itsy) Register(path string) Resource {
	r := newBaseResource(path, i)
	i.resources[path] = r
	i.router.addRoute(path, r)
	return r
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

	// Configure the server
	serverPort := DefaultPort
	if len(port) > 0 {
		serverPort = port[0]
	}
	i.Server = &http.Server{
		Addr:         serverPort,
		Handler:      i,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	i.Logger.Info("Listening on port " + serverPort)
	err := i.Server.ListenAndServe()
	if err != nil {
		i.Logger.Fatal("Server error", zap.Error(err))
	}
}
