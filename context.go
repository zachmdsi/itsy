package itsy

import (
	"net/http"

	"go.uber.org/zap"
)

type Context interface {
	// Request returns the HTTP request.
	Request() *http.Request

	// ResponseWriter returns the HTTP response writer.
	ResponseWriter() http.ResponseWriter

	// Logger returns the logger instance.
	Logger() *zap.Logger
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

func (i *Itsy) newBaseContext(req *http.Request, w http.ResponseWriter) *BaseContext {
	clonedLogger := i.Logger.With(zap.String("request_id", req.Header.Get("X-Request-Id")))
	return &BaseContext{
		request:        req,
		responseWriter: w,
		logger:         clonedLogger,
	}
}
