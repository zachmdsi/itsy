package itsy

import (
	"net/http"

	"go.uber.org/zap"
)

type (
	// Context is the interface that provides access to the information about the current request.
	Context interface {
		Request() *http.Request              // Request returns the HTTP request.
		ResponseWriter() http.ResponseWriter // ResponseWriter returns the HTTP response writer.
		Logger() *zap.Logger                 // Logger returns the logger instance.
		RenderResource(res Resource) error   // RenderResource renders a resource as JSON.
	}
	baseContext struct {
		request           *http.Request
		responseWriter    http.ResponseWriter
		logger            *zap.Logger
		contentNegotiator *ContentNegotiator
	}
)

// Request returns the HTTP request.
func (c *baseContext) Request() *http.Request { return c.request }

// Logger returns the logger instance.
func (c *baseContext) Logger() *zap.Logger { return c.logger }

// ResponseWriter returns the HTTP response writer.
func (c *baseContext) ResponseWriter() http.ResponseWriter { return c.responseWriter }

// RenderResource renders a resource as JSON.
func (c *baseContext) RenderResource(resource Resource) error {
	renderer := c.contentNegotiator.GetRenderer(c.Request().Header.Get("Accept"))
	return renderer.Render(c.ResponseWriter(), resource)
}

// newbaseContext creates a new base context.
func (i *Itsy) newBaseContext(r *http.Request, w http.ResponseWriter) *baseContext {
	clonedLogger := i.Logger.With(zap.String("request_id", r.Header.Get("X-Request-Id")))
	negotiator := i.newContentNegotiator()
	return &baseContext{
		request:           r,
		responseWriter:    w,
		logger:            clonedLogger,
		contentNegotiator: negotiator,
	}
}

// newContentNegotiator creates a new content negotiator.
func (i *Itsy) newContentNegotiator() *ContentNegotiator {
	negotiator := &ContentNegotiator{
		renderers: map[string]Renderer{
			"application/json": &JSONRenderer{},
			"application/xml":  &XMLRenderer{},
		},
	}
	return negotiator
}
