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
)

// newRouter creates a new router instance.
func newRouter(i *Itsy) *router {
	return &router{
		index: &node{
			children: make(map[string]*node),
		},
		itsy: i,
	}
}

// addRoute adds a route to the router.
func (r *router) addRoute(path string, resource Resource) {
	if path == "/" {
		r.index.resource = resource
		return
	}
	segments := splitPath(path)
	n := r.index

	for _, segment := range segments {
		if segment != "" {
			// If a direct match is found, move to the next node
			if child, ok := n.children[segment]; ok {
				n = child
			} else { // If no direct match is found, create a new node
				node := &node{
					path:     segment,
					children: make(map[string]*node),
					resource: resource,
				}

				// If the segments starts with ":", it's a parameterized route
				if strings.HasPrefix(segment, ":") {
					node.param = segment[1:]
				}

				// Add the node to the parent node
				n.children[segment] = node

				// Move to the new node
				n = node
			}
		}
	}
}

// splitPath splits a path into segments.
func splitPath(path string) []string {
	return strings.FieldsFunc(path, func(r rune) bool { return r == '/' })
}
