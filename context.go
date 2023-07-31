package itsy

import "net/http"

type (
	Context interface {
		Request() *http.Request
		Response() http.ResponseWriter
		SetParam(name, value string)
		GetParam(name string) string
	}
	baseContext struct {
		req    *http.Request
		res    *http.ResponseWriter
		params map[string]string
	}
)

func (c *baseContext) Request() *http.Request        { return c.req }
func (c *baseContext) Response() http.ResponseWriter { return *c.res }
func (c *baseContext) SetParam(name, value string)   { c.params[name] = value }
func (c *baseContext) GetParam(name string) string   { return c.params[name] }
