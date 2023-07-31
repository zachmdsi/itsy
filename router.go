package itsy

import (
	"strings"
)

type (
	// router is the main router instance.
	router struct {
		index *node // The root node of the router.
		itsy  *Itsy // The main framework instance.
	}
	// node is a node in the router.
	node struct {
		path     string           // The path of the node.
		children map[string]*node // The child nodes of the node.
		resource Resource         // The resource of the node.
		param    string           // The name of the parameter, if the node is a parameter node.
	}
	// HandlerFunc is a function that handles a request.
	HandlerFunc func(Context)
)

// newRouter creates a new router instance.
func newRouter(i *Itsy) *router {
	return &router{
		index: &node{
			children: make(map[string]*node),
			resource: nil,
		},
		itsy: i,
	}
}

// addRoute adds a route to the router.
func (r *router) addRoute(path string, resource Resource) {
	// Split the router into segments
	segments := splitPath(path)

	// Start at the root node
	n := r.index

	// For each segment in the path
	for _, segment := range segments {
		// If the segment is not empty
		if segment != "" {
			// If the segment starts with a colon, it's a parameter.
			isParam := strings.HasPrefix(segment, ":")

			if isParam {
				// Store the parameter name wihout the colon.
				paramName := segment[1:]

				// If the child node for parameters doesn't exist, create it.
				if _, ok := n.children[":"]; !ok {
					n.children[":"] = &node{
						path:     segment,
						children: make(map[string]*node),
						param:    paramName,
					}
				}

				// Move to the child node for parameters.
				n = n.children[":"]
			} else {
				// If the child node does not exist, create it.
				if _, ok := n.children[segment]; !ok {
					n.children[segment] = &node{
						path:     segment,
						children: make(map[string]*node),
					}
				}

				// Move to the child node
				n = n.children[segment]
			}
		}
	}

	// Set the resource of the node
	n.resource = resource
}

// splitPath splits a path into segments.
func splitPath(path string) []string {
	return strings.FieldsFunc(path, func(r rune) bool { return r == '/' })
}
