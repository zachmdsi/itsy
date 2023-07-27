package itsy

import "strings" 

type (
	// router is the main router tree
	router struct {
		index *node // index is the root node of the router tree
		itsy  *Itsy // itsy is the main framework instance
	}

	// node is a node in the router tree
	node struct {
		handlers map[string]HandlerFunc // handlers is a map of handlers for each method
		children map[string]*node       // children is a map of child nodes
	}
)

func newRouter(itsy *Itsy) *router {
	return &router{
		index: &node{
			handlers: make(map[string]HandlerFunc),
			children: make(map[string]*node),
		},
		itsy: itsy,
	}
}

func (r *router) addRoute(method, path string, handler HandlerFunc) {
	segments := strings.FieldsFunc(path, func(r rune) bool { return r == '/' })
	currentNode := r.index
	for _, segment := range segments {
		// If a direct match is found, move to the next node
		if child, ok := currentNode.children[segment]; ok {
			currentNode = child
		} else {
			// If no direct match is found, create a new node
			newNode := &node{
				handlers: make(map[string]HandlerFunc),
				children: make(map[string]*node),
			}
			currentNode.children[segment] = newNode
			currentNode = newNode
		}
	}
	currentNode.handlers[method] = handler
}