package itsy

import (
	"fmt"
)

type (
	// Hypermedia represents a set of hypermedia controls.
	Hypermedia struct {
		Controls map[string]HypermediaControl // The hypermedia controls.
		Fields   map[string]interface{}       // The fields.
	}
	// HypermediaControl is the interface that describes a hypermedia control.
	HypermediaControl interface {
		Render(Context) string
	}
	// Field describes a custom key value pair.
	Field struct {
		Key   string
		Value interface{}
	}
	// Link is a link to another resource.
	Link struct {
		Href string // The URL of the resource.
		Rel  string // The relationship of the resource to the current resource.
	}
)

func NewHypermedia() *Hypermedia {
	return &Hypermedia{
		Controls: make(map[string]HypermediaControl),
		Fields:   make(map[string]interface{}),
	}
}

// SetHref sets the href of the link.
func (l *Link) SetHref(href string) {
	l.Href = href
}

// Render renders the link.
func (l *Link) Render(c Context) string {
	return fmt.Sprintf("<a href=\"%s\">%s</a>", l.Href, l.Rel)
}
