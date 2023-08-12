package itsy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

type (
	// Context describes the context of a request.
	Context interface {
		Request() *http.Request                    // The HTTP request.
		Response() *Response                       // The HTTP response.
		Resource() Resource                        // The resource.
		SetResource(res Resource)                  // Set the resource.
		SetParam(name, value string)               // Set a parameter.
		GetParam(name string) string               // Get a parameter.
		Params() map[string]string                 // The parameters.
		Path() string                              // The path of the request.
		Itsy() *Itsy                               // The main framework instance.
		Mutex() *sync.RWMutex                      // The mutex.
		WriteString(s string) error                // Write a string to the response.
		WriteHTML() error                          // Write the response as HTML.
		WriteJSON() error                          // Write the response as JSON.
		CreateField(key string, value interface{}) // Create a field.
	}
	baseContext struct {
		mu       sync.RWMutex
		req      *http.Request
		res      *Response
		resource Resource
		params   map[string]string
		path     string
		itsy     *Itsy
	}
)

func newBaseContext(req *http.Request, res *Response, resource Resource, path string, itsy *Itsy) *baseContext {
	return &baseContext{
		mu:       sync.RWMutex{},
		req:      req,
		res:      res,
		resource: resource,
		params:   make(map[string]string),
		path:     path,
		itsy:     itsy,
	}
}

func (c *baseContext) Request() *http.Request      { return c.req }
func (c *baseContext) Response() *Response         { return c.res }
func (c *baseContext) Resource() Resource          { return c.resource }
func (c *baseContext) SetResource(res Resource)    { c.resource = res }
func (c *baseContext) SetParam(name, value string) { c.params[name] = value }
func (c *baseContext) GetParam(name string) string { return c.params[name] }
func (c *baseContext) Params() map[string]string   { return c.params }
func (c *baseContext) Path() string                { return c.path }
func (c *baseContext) Itsy() *Itsy                 { return c.itsy }
func (c *baseContext) Mutex() *sync.RWMutex        { return &c.mu }

func (c *baseContext) CreateField(key string, value interface{}) {
	c.Resource().Hypermedia().Fields[key] = value
}

// WriteString writes a string to the response.
func (c *baseContext) WriteString(s string) error {
	r := c.Response()
	data := []byte(s)
	written, err := r.Write(data)
	if err != nil {
		return err
	}

	if written != len(s) {
		r.itsy.sendHTTPError(StatusInternalServerError, "Response length mismatch", r.Writer, r.itsy.Logger)
		return errors.New("Response length mismatch")
	}

	return nil
}

// WriteHTML writes the response as HTML.
func (c *baseContext) WriteHTML() error {
	originalWriter := c.Response().Writer
	if originalWriter == nil {
		return errors.New("Response writer is nil")
	}

	wrapper := &responseWriterWrapper{writer: originalWriter, original: originalWriter}
	c.Response().Writer = wrapper

	// Write the initial HTML to the response.
	wrapper.Write([]byte("<html>\n  <body>\n"))

	// Write the fields to the response.
	if fields := c.Resource().Hypermedia().Fields; len(fields) > 0 {
		wrapper.Write([]byte("    <div>Fields:\n"))
		for key, value := range fields {
			wrapper.Write([]byte(fmt.Sprintf("      <div>%s: %v</div>\n", key, value)))
		}
		wrapper.Write([]byte("    </div>\n"))
	}

	// Write the hypermedia controls to the response.
	if err := writeHypermediaControls(c, wrapper); err != nil {
		return err
	}

	// Write the final HTML to the response.
	wrapper.Write([]byte("  </body>\n</html>\n"))

	wrapper.statusCode = StatusOK

	return nil
}

// writeHypermediaControls writes the hypermedia controls to the response.
func writeHypermediaControls(c Context, writer io.Writer) error {
	if resource := c.Resource(); resource != nil {
		links := resource.Links()
		if links != nil {
			writer.Write([]byte("    <div>Links:\n"))
			for _, link := range links {
				writer.Write([]byte(link.Render(c)))
			}
			writer.Write([]byte("    </div>\n"))
		}
	}
	return nil
}

// WriteJSON writes the response as JSON.
func (c *baseContext) WriteJSON() error {
	originalWriter := c.Response().Writer
	if originalWriter == nil {
		return errors.New("Response writer is nil")
	}

	// Create a map to hold the response data.
	responseData := make(map[string]interface{})

	// Write the fields and links to the response data.
	writeFields(c, responseData)
	writeLinks(c, responseData)

	// Serialize the response data to JSON.
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		return err
	}

	// Set the content type to JSON.
	originalWriter.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response.
	originalWriter.Write(jsonData)

	return nil
}

// getFullURL constructs the full URL from the request and link.
func getFullURL(c Context, link *Link) string {
	scheme := "http"
	if c.Request().TLS != nil {
		scheme = "https"
	}
	host := c.Request().Host
	return fmt.Sprintf("%s://%s%s", scheme, host, link.Href)
}

// writeFields writes the fields to the response data.
func writeFields(c Context, responseData map[string]interface{}) {
	if fields := c.Resource().Hypermedia().Fields; fields != nil {
		responseData["fields"] = fields
	}
}

// writeLinks writes the hypermedia controls to the response data.
func writeLinks(c Context, responseData map[string]interface{}) {
	links := make(map[string]string)
	if resource := c.Resource(); resource != nil {
		for rel, link := range resource.Links() {
			if params := c.Resource().GetParams(); params != nil {
				for param, value := range params {
					placeholder := fmt.Sprintf(":%s", param)
					link.SetHref(strings.Replace(link.Href, placeholder, value, -1))
				}
			}
			fullURL := getFullURL(c, link)
			links[rel] = fullURL
		}
	}
	responseData["links"] = links
}
