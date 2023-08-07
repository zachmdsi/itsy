package itsy

import (
	"bytes"
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

// HypermediaMiddleware is a middleware that adds hypermedia controls to the response.
func HypermediaMiddleware(next HandlerFunc) HandlerFunc {
	return func(c Context) error {
		originalWriter := c.Response().Writer
		if originalWriter == nil {
			return fmt.Errorf("response writer is nil")
		}

		buffer := new(bytes.Buffer)
		wrapper := &responseWriterWrapper{writer: buffer, original: originalWriter}
		c.Response().Writer = wrapper
		buffer.Write([]byte("<html><body>"))

		err := next(c)
		if err != nil {
			c.Itsy().Logger.Error("Handler Error", zap.Error(err))
			c.Response().Writer.WriteHeader(StatusInternalServerError)
			return err
		}

		if resource := c.Resource(); resource != nil {
			hypermedia := resource.Hypermedia()
			if hypermedia != nil && len(hypermedia.Controls) > 0 {
				buffer.Write([]byte("<div>Links:"))
				for _, control := range hypermedia.Controls {
					if link, ok := control.(*Link); ok {
						if params := c.Resource().GetParams(); params != nil {
							for param, value := range params {
								placeholder := fmt.Sprintf(":%s", param)
								link.Href = strings.Replace(link.Href, placeholder, value, -1)
							}
						}
						buffer.Write([]byte(
							fmt.Sprintf("<a href=\"%s\">%s</a>", link.Href, link.Rel)))
					}
				}
				buffer.Write([]byte("</div>"))
			}
		}

		buffer.Write([]byte("</body></html>"))
		_, err = originalWriter.Write(buffer.Bytes())
		if err != nil {
			c.Itsy().Logger.Error("Original Writer Error", zap.Error(err))
			c.Response().Writer.WriteHeader(StatusInternalServerError)
			return err
		}

		wrapper.statusCode = StatusOK
		c.Itsy().Logger.Info("Response", zap.Int("status", wrapper.statusCode))

		return nil
	}
}
