package itsy

import (
	"regexp"
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
		path     string         // The path of the node.
		regex    *regexp.Regexp // The regex of the segment, if it's a parameterized route.
		children []*node        // The child nodes of the node.
		resource Resource       // The resource of the node.
		param    string         // The name of the parameter, if the node is a parameter node.
	}
)

// newRouter creates a new router instance.
func newRouter(i *Itsy) *router {
	return &router{
		index: &node{
			children: make([]*node, 0),
		},
		itsy: i,
	}
}

// addRoute adds a route to the router.
func (r *router) addRoute(path string, resource Resource) {
	segments := splitPath(path)

	if len(segments) == 0 {
		r.index.resource = resource
		return
	}

	currentNode := r.index

	for _, segment := range segments {
		if segment != "" {
			found := false
			for _, child := range currentNode.children {
				if child.path == segment || (child.regex != nil && child.regex.MatchString(segment)) {
					currentNode = child
					found = true
					break
				}
			}

			if !found {
				newNode := &node{
					path:     segment,
					children: make([]*node, 0),
					resource: resource,
				}

				// If the segments starts with ":", it's a parameterized route
				if strings.HasPrefix(segment, ":") {
					newNode.param = segment[1:]
					newNode.regex = regexp.MustCompile("^[a-zA-Z0-9_]+$") // Compile regex for the parameter
				}

				// Add the node to the parent node
				currentNode.children = append(currentNode.children, newNode)

				// Move to the new node
				currentNode = newNode
			}
		}
	}
}

// splitPath splits a path into segments.
func splitPath(path string) []string {
	return strings.FieldsFunc(path, func(r rune) bool { return r == '/' })
}
