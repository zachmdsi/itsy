package itsy

import (
	"net/http"

	"go.uber.org/zap"
)

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
