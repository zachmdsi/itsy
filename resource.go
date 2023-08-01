package itsy

type (
	// Resource is the interface that describes a RESTful resource.
	Resource interface {
		Link(href, rel string) string      // Link to another resource.
		Links() map[string]Link            // Get the links of the resource.
		Hypermedia() *Hypermedia           // Get the hypermedia of the resource.
		Handler(method string) HandlerFunc // Get the handler of the resource. 
		GET(handler HandlerFunc)           // GET the resource.
	}
	// BaseResource is the base implementation of the Resource interface.
	BaseResource struct {
		links      map[string]Link
		handlers   map[string]HandlerFunc
		hypermedia *Hypermedia
	}
	CustomResource struct {
		BaseResource
	}
)

func newBaseResource() *BaseResource {
	return &BaseResource{
		links:      make(map[string]Link),
		handlers:   make(map[string]HandlerFunc),
		hypermedia: &Hypermedia{Controls: make(map[string]HypermediaControl)},
	}
}

func newCustomResource() *CustomResource {
	return &CustomResource{
		BaseResource: *newBaseResource(),
	}
}

// Link links to another resource.
func (r *BaseResource) Link(href, rel string) string {
	return ""
}

// Links gets the links of the resource.
func (r *BaseResource) Links() map[string]Link {
	return r.links
}

// Hypermedia gets the hypermedia of the resource.
func (r *BaseResource) Hypermedia() *Hypermedia {
	return r.hypermedia
}

// Handler gets the handler of the resource for the given method.
func (r *BaseResource) Handler(method string) HandlerFunc {
	handler, ok := r.handlers[method]
	if !ok {
		return nil
	}
	return handler
}

// GET calls the handler when the resource is requested with the GET method.
func (r *BaseResource) GET(handler HandlerFunc) {
	r.handlers[GET] = handler
}
