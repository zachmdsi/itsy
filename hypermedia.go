package itsy

import (
	"fmt"
	"strings"
)

type (
	// Hypermedia represents a set of hypermedia controls.
	Hypermedia struct {
		Links map[string]*Link // The links.
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
		Links: make(map[string]*Link),
	}
}

// Render renders the link.
func (l *Link) Render(c Context) string {
	href := l.Href
	if params := c.GetParams(); params != nil {
		for _, param := range params {
			placeholder := fmt.Sprintf(":%s", param.Name)
			href = strings.Replace(href, placeholder, param.Value, -1)
		}
	}
	return fmt.Sprintf("<a href=\"%s\">%s</a>", href, l.Rel)
}
