package itsy

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type (
	// Resource is the interface that describes a RESTful resource.
	Resource interface {
		GET(handler HandlerFunc)             // Set the GET handler of the resource.
		GetParams() map[string]string        // Get the parameters of the resource.
		SetParam(name, value string)         // Set a parameter.
		Hypermedia() *Hypermedia             // Get the hypermedia of the resource.
		Handler(method string) HandlerFunc   // Get the handler of the resource.
		Itsy() *Itsy                         // Get the main framework instance.
		Link(res Resource, rel string) error // Link to another resource.
		Links() map[string]*Link             // Get the links of the resource.
		Path() string                        // Get the path of the resource.
	}
	// BaseResource is the base implementation of the Resource interface.
	baseResource struct {
		mu         sync.RWMutex
		handlers   map[string]HandlerFunc
		hypermedia *Hypermedia
		itsy       *Itsy
		params     map[string]string
		path       string
	}
)

// newBaseResource creates a new base resource.
func newBaseResource(path string, i *Itsy) *baseResource {
	return &baseResource{
		handlers:   make(map[string]HandlerFunc),
		hypermedia: &Hypermedia{Controls: make(map[string]HypermediaControl)},
		itsy:       i,
		params:     make(map[string]string),
		path:       path,
	}
}

// Link links to another resource.
func (r *baseResource) Link(res Resource, rel string) error {
	path := res.Path()
	if !r.Itsy().ResourceExists(path) {
		return errors.New("Resource does not exist")
	}

	if params := res.GetParams(); params != nil {
		for param, value := range params {
			placeholder := fmt.Sprintf(":%s", param)
			path = strings.Replace(path, placeholder, value, -1)
		}
	}

	link := Link{Href: path, Rel: rel}
	r.hypermedia.Controls[rel] = &link

	return nil
}

// Links gets the links of the resource.
func (r *baseResource) Links() map[string]*Link {
	links := make(map[string]*Link)
	for rel, control := range r.hypermedia.Controls {
		if link, ok := control.(*Link); ok {
			links[rel] = link
		}
	}
	return links
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

// Path gets the path of the resource.
func (r *baseResource) Path() string {
	return r.path
}

// Itsy gets the main framework instance.
func (r *baseResource) Itsy() *Itsy {
	return r.itsy
}

// GetParams gets the parameters of the resource.
func (r *baseResource) GetParams() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.params
}

// SetParam sets a parameter.
func (r *baseResource) SetParam(name, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.params[name] = value
}
