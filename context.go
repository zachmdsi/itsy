package itsy

import "net/http"

type (
	// Context describes the context of a request.
	Context interface {
		Request() *http.Request      // The HTTP request.
		Response() *Response         // The HTTP response.
		Resource() Resource          // The resource.
		SetResource(res Resource) 	 // Set the resource.
		SetParam(name, value string) // Set a parameter.
		GetParam(name string) string // Get a parameter.
		Params() map[string]string   // The parameters.
		Path() string                // The path of the request.
		Itsy() *Itsy                 // The main framework instance.
	}
	baseContext struct {
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
