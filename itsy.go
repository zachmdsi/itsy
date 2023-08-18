package itsy

import (
	"net/http"
	_ "net/http/pprof"

	"go.uber.org/zap"
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

// ResourceExistsWithMethod returns true if a resource exists given a path and method.
func (i *Itsy) ResourceExistsWithMethod(path, method string) bool {
	if !i.ResourceExists(path) {
		return false
	}
	resource := i.Resource(path)
	return resource.Handler(method) != nil
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
