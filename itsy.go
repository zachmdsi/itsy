package itsy

import (
	"bufio"
	"net/http"
	_ "net/http/pprof"
	"os"

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

// New creates a new Itsy instance.
func New() *Itsy {
	i := &Itsy{
		resources: make(map[string]Resource),
		Logger:    setupLogger(),
	}
	i.router = newRouter(i)
	return i
}

// Register registers a resource to the Itsy instance.
func (i *Itsy) Register(path string) Resource {
	baseResource := newBaseResource(path, i)
	i.resources[path] = baseResource
	i.router.addRoute(path, baseResource)
	return baseResource
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
						c.AddParam(child.param, segment)
						c.SetResource(child.resource)
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
		switch req.Method {
		case GET:
			if n.resource.Handler(GET) == nil {
				i.sendHTTPError(StatusMethodNotAllowed, "Handler does not exist for the request method", res, i.Logger)
				return
			}
			callHandler(n.resource, GET, c)
		case POST:
			if n.resource.Handler(POST) == nil {
				i.sendHTTPError(StatusMethodNotAllowed, "Handler does not exist for the request method", res, i.Logger)
				return
			}
			callHandler(n.resource, POST, c)
		default:
			i.sendHTTPError(StatusMethodNotAllowed, "Handler does not exist for the request method", res, i.Logger)
		}
	} else {
		i.sendHTTPError(StatusNotFound, "Resource does not exist", res, i.Logger)
	}
}

// callHandler calls the handler of the resource.
func callHandler(resource Resource, method string, c Context) error {
	handler := resource.Handler(method)
	if handler == nil {
		return nil
	}
	return handler(c)
}

// setupLogger sets up the logger for the Itsy instance.
func setupLogger() *zap.Logger {
	// Encoder Configuration
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.EpochTimeEncoder // Optimized time encoding
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	logWriter := zapcore.AddSync(bufio.NewWriter(os.Stdout))

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), logWriter), zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

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
