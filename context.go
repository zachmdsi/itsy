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
	}
	baseContext struct {
		request           *http.Request
		responseWriter    http.ResponseWriter
		logger            *zap.Logger
	}
)

func (c *baseContext) Request() *http.Request              { return c.request }
func (c *baseContext) Logger() *zap.Logger                 { return c.logger }
func (c *baseContext) ResponseWriter() http.ResponseWriter { return c.responseWriter }

// newbaseContext creates a new base context.
func (i *Itsy) newBaseContext(r *http.Request, w http.ResponseWriter) *baseContext {
	clonedLogger := i.Logger.With(zap.String("request_id", r.Header.Get("X-Request-Id")))
	return &baseContext{
		request:           r,
		responseWriter:    w,
		logger:            clonedLogger,
	}
}
