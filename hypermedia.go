package itsy

import "regexp"

type (
	// Hypermedia represents a set of hypermedia controls.
	Hypermedia struct {
		Links []Link // The links.
	}
	// Link is a link to another resource.
	Link struct {
		re   *regexp.Regexp
		Href string // The URL of the resource.
		Rel  string // The relationship of the resource to the current resource.
	}
)

// newHypermedia creates a new hypermedia instance.
func newHypermedia() *Hypermedia {
	return &Hypermedia{
		Links: make([]Link, 0),
	}
}

// newLink creates a new link.
func newLink(href, rel string) Link {
	re := regexp.MustCompile(`:(\w+)`) // Pre-compile the regular expression.
	return Link{
		re:   re,
		Href: href,
		Rel:  rel,
	}
}
