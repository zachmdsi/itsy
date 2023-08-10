package itsy

import (
	"errors"
	"fmt"
	"io"
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
)

// HypermediaMiddleware is a middleware that processes a handler and adds hypermedia controls to the response.
func HypermediaMiddleware(next HandlerFunc) HandlerFunc {
	return func(c Context) error {
		originalWriter := c.Response().Writer
		if originalWriter == nil {
			return errors.New("Response writer is nil")
		}

		wrapper := &responseWriterWrapper{writer: originalWriter, original: originalWriter}
		c.Response().Writer = wrapper

		// Write the initial HTML to the response.
		wrapper.Write([]byte("<html><body>"))

		// Process the handler.
		err := next(c)
		if err != nil {
			c.Itsy().Logger.Error("Handler Error", zap.Error(err))
			c.Response().Writer.WriteHeader(StatusInternalServerError)
			return err
		}

		// Write the hypermedia controls to the response.
		if err := writeHypermediaControls(c, wrapper); err != nil {
			return err
		}

		// Write the final HTML to the response.
		wrapper.Write([]byte("</body></html>"))
		wrapper.statusCode = StatusOK

		c.Itsy().Logger.Info("Response", zap.Int("status", wrapper.statusCode))
		return nil
	}
}

// writeHypermediaControls writes the hypermedia controls to the response.
func writeHypermediaControls(c Context, writer io.Writer) error {
	if resource := c.Resource(); resource != nil {
		hypermedia := resource.Hypermedia()
		if hypermedia != nil && len(hypermedia.Controls) > 0 {
			writer.Write([]byte("<div>Links:"))
			for _, control := range hypermedia.Controls {
				if err := writeLink(c, control, writer); err != nil {
					return err
				}
			}
			writer.Write([]byte("</div>"))
		}
	}
	return nil
}

// writeLink writes a link to the response.
func writeLink(c Context, control HypermediaControl, writer io.Writer) error {
	if link, ok := control.(*Link); ok {
		if params := c.Resource().GetParams(); params != nil {
			for param, value := range params {
				placeholder := fmt.Sprintf(":%s", param)
				link.Href = strings.Replace(link.Href, placeholder, value, -1)
			}
		}
		writer.Write([]byte(fmt.Sprintf("<a href=\"%s\">%s</a>", link.Href, link.Rel)))
	}
	return nil
}
