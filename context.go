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
		Params() map[string]string           // Params returns the parameters extracted from the URL path.
		SetParam(name, value string)         // SetParam sets the value of a parameter.
	}
	baseContext struct {
		request        *http.Request
		responseWriter http.ResponseWriter
		logger         *zap.Logger
		params         map[string]string
	}
)

func (c *baseContext) Request() *http.Request              { return c.request }
func (c *baseContext) Logger() *zap.Logger                 { return c.logger }
func (c *baseContext) ResponseWriter() http.ResponseWriter { return c.responseWriter }
func (c *baseContext) Params() map[string]string           { return c.params }
func (c *baseContext) SetParam(name, value string)         { c.params[name] = value }

// newbaseContext creates a new base context.
func (i *Itsy) newBaseContext(r *http.Request, w http.ResponseWriter) *baseContext {
	clonedLogger := i.Logger.With(zap.String("request_id", r.Header.Get("X-Request-Id")))
	return &baseContext{
		request:        r,
		responseWriter: w,
		logger:         clonedLogger,
		params:         make(map[string]string),
	}
}
