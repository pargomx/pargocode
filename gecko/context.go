package gecko

import (
	"net/http"
	"net/url"
)

// Representa la solicitud HTTP actual y
// ofrece los medios para responderla.
type Context struct {
	request  *http.Request
	response *Response
	path     string
	query    url.Values
	gecko    *Gecko
	Sesion   Sesion // Sesión del usuario autenticado o anónimo.
}

func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) Response() *Response {
	return c.response
}
