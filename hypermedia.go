package itsy

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type (
	// Hypermedia represents a set of hypermedia controls.
	Hypermedia struct {
		Controls map[string]HypermediaControl
	}
	// HypermediaControl is the interface that describes a hypermedia control.
	HypermediaControl interface {
		Render(Context) string
	}
	// Link is a link to another resource.
	Link struct {
		Href string // The URL of the resource.
		Rel  string // The relationship of the resource to the current resource.
	}
	// responseWriterWrapper is a wrapper around the response writer.
	responseWriterWrapper struct {
		writer     io.Writer
		statusCode int
		original   http.ResponseWriter
	}
)

// Write writes to the response.
func (res *responseWriterWrapper) Write(b []byte) (n int, err error) {
	return res.writer.Write(b)
}

// Header gets the header of the response.
func (res *responseWriterWrapper) Header() http.Header {
	return res.original.Header()
}

// WriteHeader writes the header of the response.
func (res *responseWriterWrapper) WriteHeader(code int) {
	res.statusCode = code
	res.original.WriteHeader(code)
}

// HypermediaMiddlware is a middleware that processes a handler and adds hypermedia controls to the response.
func HypermediaMiddleware(next HandlerFunc) HandlerFunc {
	return func(c Context) error {
		originalWriter, buffer, err := prepareResponse(c)
		if err != nil {
			return err
		}

		if err := processHandler(c, next, buffer); err != nil {
			return err
		}

		if err := writeHypermediaControls(c, buffer); err != nil {
			return err
		}

		return finalizeResponse(c, originalWriter, buffer)
	}
}

// prepareResponse writes the initial HTML to the response.
func prepareResponse(c Context) (http.ResponseWriter, *bytes.Buffer, error) {
	originalWriter := c.Response().Writer
	if originalWriter == nil {
		return nil, nil, errors.New("Response writer is nil")
	}

	buffer := new(bytes.Buffer)
	wrapper := &responseWriterWrapper{writer: buffer, original: originalWriter}
	c.Response().Writer = wrapper
	buffer.Write([]byte("<html><body>"))
	return originalWriter, buffer, nil
}

// processHandler processes the next handler in the chain. 
func processHandler(c Context, next HandlerFunc, buffer *bytes.Buffer) error {
	err := next(c)
	if err != nil {
		c.Itsy().Logger.Error("Handler Error", zap.Error(err))
		c.Response().Writer.WriteHeader(StatusInternalServerError)
		return err
	}
	return nil
}

// writeHypermediaControls writes the hypermedia controls to the response.
func writeHypermediaControls(c Context, buffer *bytes.Buffer) error {
	if resource := c.Resource(); resource != nil {
		hypermedia := resource.Hypermedia()
		if hypermedia != nil && len(hypermedia.Controls) > 0 {
			buffer.Write([]byte("<div>Links:"))
			for _, control := range hypermedia.Controls {
				if err := writeLink(c, control, buffer); err != nil {
					return err
				}
			}
			buffer.Write([]byte("</div>"))
		}
	}
	return nil
}

// writeLink writes a link to the response.
func writeLink(c Context, control HypermediaControl, buffer *bytes.Buffer) error {
	if link, ok := control.(*Link); ok {
		if params := c.Resource().GetParams(); params != nil {
			for param, value := range params {
				placeholder := fmt.Sprintf(":%s", param)
				link.Href = strings.Replace(link.Href, placeholder, value, -1)
			}
		}
		buffer.Write([]byte(fmt.Sprintf("<a href=\"%s\">%s</a>", link.Href, link.Rel)))
	}
	return nil
}

// finalizeResponse writes the final HTML to the response.
func finalizeResponse(c Context, originalWriter http.ResponseWriter, buffer *bytes.Buffer) error {
	buffer.Write([]byte("</body></html>"))
	_, err := originalWriter.Write(buffer.Bytes())
	if err != nil {
		c.Itsy().Logger.Error("Original Writer Error", zap.Error(err))
		c.Response().Writer.WriteHeader(StatusInternalServerError)
		return err
	}

	wrapper, ok := c.Response().Writer.(*responseWriterWrapper)
	if ok {
		wrapper.statusCode = StatusOK
	}
	c.Itsy().Logger.Info("Response", zap.Int("status", wrapper.statusCode))
	return nil
}
