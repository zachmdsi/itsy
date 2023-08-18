package itsy

import (
	"errors"
	"io"
	"net/http"
	"text/template"

	"go.uber.org/zap"
)

type (
	// Context describes the context of a request.
	Context interface {
		Request() *http.Request                        // The HTTP request.
		Response() *Response                           // The HTTP response.
		Resource() Resource                            // The resource.
		SetResource(res Resource)                      // Set the resource.
		AddParam(name, value string)                   // Set a parameter.
		GetParamValue(name string) string              // Get a parameter.
		GetParams() []Param                            // The parameters.
		Path() string                                  // The path of the request.
		Itsy() *Itsy                                   // The main framework instance.
		WriteString(s string) error                    // Write a string to the response.
		WriteHTML() error                              // Write the response as HTML.
		SetTemplateRenderer(renderer TemplateRenderer) // Set the template renderer.
		GetTemplateRenderer() TemplateRenderer         // Get the template renderer.
	}
	// TemplateRenderer is the interface that describes a template renderer.
	TemplateRenderer interface {
		RenderTemplate(w io.Writer, c Context) error            // Render a template.
		RenderLinks(c Context, w io.Writer, links []Link) error // Render links.
	}
	// Param is a parameter.
	Param struct {
		Name  string
		Value string
	}
	// defaultTemplateRenderer is the default template renderer.
	defaultTemplateRenderer struct{}
	// baseContext is the base implementation of the Context interface.
	baseContext struct {
		req              *http.Request
		res              *Response
		resource         Resource
		params           []Param
		path             string
		itsy             *Itsy
		templateRenderer TemplateRenderer
	}
)

// newBaseContext creates a new base context.
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

// Context interface implementation.
func (c *baseContext) Request() *http.Request   { return c.req }
func (c *baseContext) Response() *Response      { return c.res }
func (c *baseContext) Resource() Resource       { return c.resource }
func (c *baseContext) SetResource(res Resource) { c.resource = res }
func (c *baseContext) Path() string             { return c.path }
func (c *baseContext) Itsy() *Itsy              { return c.itsy }

func (c *baseContext) SetTemplateRenderer(renderer TemplateRenderer) {
	c.templateRenderer = renderer
}

func (c *baseContext) GetTemplateRenderer() TemplateRenderer {
	return c.templateRenderer
}

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

	renderer := c.GetTemplateRenderer()
	if renderer == nil {
		renderer = &defaultTemplateRenderer{}
	}

	// Render the main HTML template.
	if err := renderer.RenderTemplate(originalWriter, c); err != nil {
		return err
	}

	// Render the links using the standard link template.
	links := c.Resource().Links()
	if len(links) > 0 {
		// Render the links.
		if err := renderer.RenderLinks(c, originalWriter, links); err != nil {
			return err
		}
	}

	return nil
}

// Default template renderer implementation.
func (r *defaultTemplateRenderer) RenderTemplate(w io.Writer, c Context) error {
	return nil
}

func (r *defaultTemplateRenderer) RenderLinks(c Context, w io.Writer, links []Link) error {
	// Modify the href attributes to replace the placeholders with the corresponding parameter values.
	for i, link := range links {
		href := link.Href
		href = link.re.ReplaceAllStringFunc(href, func(s string) string {
			paramName := s[1:]
			return c.GetParamValue(paramName)
		})
		links[i].Href = href
	}

	// Parse the standard link template.
	t, err := template.New("links").Parse(linkTemplate)
	if err != nil {
		return err
	}

	// Execute the template with the links and write the result to the response
	if err := t.Execute(w, links); err != nil {
		return err
	}

	return nil
}
