package itsy

import (
	"errors"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type (
	// Context describes the context of a request.
	Context interface {
		Request() *http.Request           // The HTTP request.
		Response() *Response              // The HTTP response.
		Resource() Resource               // The resource.
		SetResource(res Resource)         // Set the resource.
		AddParam(name, value string)      // Set a parameter.
		GetParamValue(name string) string // Get a parameter.
		GetParams() []Param               // The parameters.
		Path() string                     // The path of the request.
		Itsy() *Itsy                      // The main framework instance.
		WriteString(s string) error       // Write a string to the response.
		WriteHTML() error                 // Write the response as HTML.
	}
	Param struct {
		Name  string
		Value string
	}
	baseContext struct {
		req      *http.Request
		res      *Response
		resource Resource
		params   []Param
		path     string
		itsy     *Itsy
	}
)

func newBaseContext(req *http.Request, res *Response, resource Resource, path string, itsy *Itsy) *baseContext {
	return &baseContext{
		req:      req,
		res:      res,
		resource: resource,
		params:   make([]Param, 0),
		path:     path,
		itsy:     itsy,
	}
}

func (c *baseContext) Request() *http.Request   { return c.req }
func (c *baseContext) Response() *Response      { return c.res }
func (c *baseContext) Resource() Resource       { return c.resource }
func (c *baseContext) SetResource(res Resource) { c.resource = res }
func (c *baseContext) Path() string             { return c.path }
func (c *baseContext) Itsy() *Itsy              { return c.itsy }

func (c *baseContext) GetParams() []Param {
	return c.params
}

func (c *baseContext) AddParam(name, value string) {
	c.params = append(c.params, Param{Name: name, Value: value})
}

func (c *baseContext) GetParamValue(name string) string {
	for _, param := range c.params {
		if param.Name == name {
			return param.Value
		}
	}
	c.itsy.Logger.Error("Parameter not found", zap.String("name", name))
	return ""
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
