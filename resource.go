package itsy

type (
	// Resource is the interface that describes a RESTful resource.
	Resource interface {
		Render(Context) string        // Render the resource.
		Link(href, rel string) string // Link to another resource.
		Links() map[string]Link       // Get the links of the resource.
		Hypermedia() *Hypermedia      // Get the hypermedia of the resource.
	}
	// BaseResource is the base implementation of the Resource interface.
	BaseResource struct {
		links      map[string]Link
		renderFunc func(Context) string
		hypermedia *Hypermedia
	}
)

func (r *BaseResource) Render(c Context) string {
	return r.renderFunc(c)
}

func (r *BaseResource) Link(href, rel string) string {
	return ""
}

func (r *BaseResource) Links() map[string]Link {
	return r.links
}

func (r *BaseResource) Hypermedia() *Hypermedia {
	return r.hypermedia
}
