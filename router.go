package itsy

import (
	"net/http"
	"strings"
)

// NewRouter creates a new router instance.
func NewRouter(itsy *Itsy) *Router {
	return &Router{
		Index: &Route{
			Handlers: make(map[string]HandlerFunc),
			Children: make(map[string]*Route),
		},
		itsy: itsy,
	}
}

// Handle registers a handler for a given HTTP method and path.
func (r *Router) Handle(method string, route string, handler HandlerFunc) {
	// Apply middlewares in reverse order so that the first middleware added is the first to be called.
	for i := len(r.itsy.middlewares) - 1; i >= 0; i-- {
		handler = r.itsy.middlewares[i](handler)
	}

	segments := strings.FieldsFunc(route, func(r rune) bool { return r == '/' })
	currentNode := r.Index
	for _, segment := range segments {
		// If the segment doesn't exist, create it.
		if currentNode.Children[segment] == nil {
			currentNode.Children[segment] = &Route{
				Path:     segment,
				Handlers: make(map[string]HandlerFunc),
				Children: make(map[string]*Route),
				IsParam:  strings.HasPrefix(segment, ":"),
			}
		}
		// Move to the next node.
		currentNode = currentNode.Children[segment]
	}

	// Register the handler for the given method.
	currentNode.Handlers[method] = handler
}

// ServeHTTP is the entry point for all HTTP requests.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	segments := strings.FieldsFunc(req.URL.Path, func(r rune) bool { return r == '/' })
	currentNode := r.Index
	params := make(map[string]string)
	for _, segment := range segments {
		// If a direct match is found, move to the next node
		if child, ok := currentNode.Children[segment]; ok {
			currentNode = child
		} else {
			// If a direct match isn't found, try to match a parameterized route
			found := false
			for key, child := range currentNode.Children {
				if strings.HasPrefix(key, ":") {
					params[key[1:]] = segment
					currentNode = child
					found = true
					break
				}
			}

			// If no parameterized route is found, return a 404.
			if !found {
				HTTPError(http.StatusNotFound, w, req)
				return
			}
		}
	}

	// If a handler is found for the given method, call it.
	if handler, ok := currentNode.Handlers[req.Method]; ok {
		ctx := r.itsy.newBaseContext(req, w)
		// TODO: add params to context
		handler(ctx)
	} else {
		// If no handler is found, return a 405.
		HTTPError(http.StatusMethodNotAllowed, w, req)
	}
}
