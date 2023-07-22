package itsy

import "net/http"

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}