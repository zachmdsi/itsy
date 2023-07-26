package itsy

import (
	"strings"
)

type ContentNegotiator struct {
	renderers map[string]Renderer
}

func NewContentNegotiator() *ContentNegotiator {
	cn := &ContentNegotiator{
		renderers: make(map[string]Renderer),
	}
	cn.RegisterRenderer("application/json", &JSONRenderer{})
	cn.RegisterRenderer("application/xml", &XMLRenderer{})
	return cn
}

func (cn *ContentNegotiator) RegisterRenderer(contentType string, renderer Renderer) {
	cn.renderers[contentType] = renderer
}

func (cn *ContentNegotiator) GetRenderer(acceptHeader string) Renderer {
	acceptedTypes := strings.Split(acceptHeader, ",")
	for _, acceptedType := range acceptedTypes {
		if renderer, ok := cn.renderers[acceptedType]; ok {
			return renderer
		}
	}
	// Return a default renderer if no match is found.
	return cn.renderers["application/json"]
}
