package itsy

import (
	"errors"
)

type (
	// Resource is the interface that describes a RESTful resource.
	Resource interface {
		GET(HandlerFunc)                     // Set the GET handler of the resource.
		POST(HandlerFunc)                    // Set the POST handler of the resource.
		PUT(HandlerFunc)                     // Set the PUT handler of the resource.
		PATCH(HandlerFunc)                   // Set the PATCH handler of the resource.
		DELETE(HandlerFunc)                  // Set the DELETE handler of the resource.
		Hypermedia() *Hypermedia             // Get the hypermedia of the resource.
		Handler(method string) HandlerFunc   // Get the handler of the resource.
		Itsy() *Itsy                         // Get the main framework instance.
		Link(res Resource, rel string) error // Link to another resource.
		Links() map[string]*Link             // Get the links of the resource.
		Path() string                        // Get the path of the resource.
	}
	// BaseResource is the base implementation of the Resource interface.
	baseResource struct {
		handlers   map[string]HandlerFunc
		hypermedia *Hypermedia
		itsy       *Itsy
		path       string
	}
)

// newBaseResource creates a new base resource.
func newBaseResource(path string, i *Itsy) *baseResource {
	return &baseResource{
		handlers:   make(map[string]HandlerFunc),
		hypermedia: newHypermedia(),
		itsy:       i,
		path:       path,
	}
}

// Link links to another resource.
func (r *baseResource) Link(res Resource, rel string) error {
	path := res.Path()
	if !r.Itsy().ResourceExists(path) {
		return errors.New("Resource does not exist")
	}

	link := Link{Href: path, Rel: rel}
	r.hypermedia.Links[rel] = &link

	return nil
}

// Links gets the links of the resource.
func (r *baseResource) Links() map[string]*Link {
	return r.hypermedia.Links
}

// Hypermedia gets the hypermedia of the resource.
func (r *baseResource) Hypermedia() *Hypermedia {
	return r.hypermedia
}

// Handler gets the handler of the resource for the given method.
func (r *baseResource) Handler(method string) HandlerFunc {
	handler, ok := r.handlers[method]
	if !ok {
		return nil
	}
	return handler
}

// GET calls the handler when the resource is requested with the GET method.
func (r *baseResource) GET(handler HandlerFunc) {
	r.handlers[GET] = handler
}

// POST calls the handler when the resource is requested with the POST method.
func (r *baseResource) POST(handler HandlerFunc) {
	r.handlers[POST] = handler
}

// PUT calls the handler when the resource is requested with the PUT method.
func (r *baseResource) PUT(handler HandlerFunc) {
	r.handlers[PUT] = handler
}

// PATCH calls the handler when the resource is requested with the PATCH method.
func (r *baseResource) PATCH(handler HandlerFunc) {
	r.handlers[PATCH] = handler
}

// DELETE calls the handler when the resource is requested with the DELETE method.
func (r *baseResource) DELETE(handler HandlerFunc) {
	r.handlers[DELETE] = handler
}

// Path gets the path of the resource.
func (r *baseResource) Path() string {
	return r.path
}

// Itsy gets the main framework instance.
func (r *baseResource) Itsy() *Itsy {
	return r.itsy
}
