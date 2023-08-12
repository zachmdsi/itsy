package itsy

import (
	"fmt"
	"strings"
)

type (
	// Hypermedia represents a set of hypermedia controls.
	Hypermedia struct {
		Links  map[string]*Link       // The links.
		Fields map[string]interface{} // The fields.
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

// newHypermedia creates a new hypermedia instance.
func newHypermedia() *Hypermedia {
	return &Hypermedia{
		Links:  make(map[string]*Link),
		Fields: make(map[string]interface{}),
	}
}

// SetHref sets the href of the link.
func (l *Link) SetHref(href string) {
	l.Href = href
}

// Render renders the link.
func (l *Link) Render(c Context) string {
	href := l.Href
	if params := c.Resource().GetParams(); params != nil {
		for param, value := range params {
			placeholder := fmt.Sprintf(":%s", param)
			href = strings.Replace(href, placeholder, value, -1)
		}
	}
	return fmt.Sprintf("<a href=\"%s\">%s</a>", href, l.Rel)
}
