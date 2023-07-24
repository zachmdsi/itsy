package itsy

import (
	"encoding/json"
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
	BaseContext struct {
		request        *http.Request
		responseWriter http.ResponseWriter
		logger         *zap.Logger
	}
)

func (c *BaseContext) Request() *http.Request { return c.request }
func (c *BaseContext) Logger() *zap.Logger { return c.logger }
func (c *BaseContext) ResponseWriter() http.ResponseWriter { return c.responseWriter }

// RenderResource renders a resource as JSON.
func (c *BaseContext) RenderResource(res Resource) error {
	c.ResponseWriter().Header().Set("Content-Type", "application/json")

	// Convert to JSON
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return err
	}

	// Write the JSON to the response
	_, err = c.ResponseWriter().Write(jsonBytes)
	if err != nil {
		return err
	}

	return nil
}

// newBaseContext creates a new base context.
func (i *Itsy) newBaseContext(r *http.Request, w http.ResponseWriter) *BaseContext {
	clonedLogger := i.Logger.With(zap.String("request_id", r.Header.Get("X-Request-Id")))
	return &BaseContext{
		request:        r,
		responseWriter: w,
		logger:         clonedLogger,
	}
}
