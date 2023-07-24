package itsy

import (
	"net/http"

	"go.uber.org/zap"
)

type Context interface {
	Request()        *http.Request       // Request returns the HTTP request.
	ResponseWriter() http.ResponseWriter // ResponseWriter returns the HTTP response writer.
	Logger()         *zap.Logger         // Logger returns the logger instance.
}

type BaseContext struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	logger         *zap.Logger
}

func (c *BaseContext) Request() *http.Request {
	return c.request
}

func (c *BaseContext) Logger() *zap.Logger {
	return c.logger
}

func (c *BaseContext) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

func (i *Itsy) newBaseContext(r *http.Request, w http.ResponseWriter) *BaseContext {
	clonedLogger := i.Logger.With(zap.String("request_id", r.Header.Get("X-Request-Id")))
	return &BaseContext{
		request:        r,
		responseWriter: w,
		logger:         clonedLogger,
	}
}
