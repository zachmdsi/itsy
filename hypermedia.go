package itsy

type (
	// Hypermedia represents a set of hypermedia controls.
	Hypermedia struct {
		Controls map[string]HypermediaControl
	}
	// HypermediaControl is the interface that describes a hypermedia control.
	HypermediaControl interface {
		Render(Context) string
	}
	// Link is a link to another resource.
	Link struct {
		Href string // The URL of the resource.
		Rel  string // The relationship of the resource to the current resource.
	}
)
