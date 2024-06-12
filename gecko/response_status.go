package gecko

import (
	"fmt"
	"net/http"
)

// ================================================================ //
// ========== Respuestas satisfactorias (2xx) ===================== //

func (c *Context) StringOk(msg string) error {
	c.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Response().WriteHeader(200)
	c.Response().Writer.Write([]byte(msg))
	return nil
}

// String sends a string response with status code 200 OK.
func (c *Context) StatusOk(msg string) (err error) {
	return c.Blob(http.StatusOK, MIMETextPlainCharsetUTF8, []byte(msg))
}

// String sends a string response with status code 200 OK.
func (c *Context) StatusOkf(format string, a ...any) (err error) {
	return c.Blob(http.StatusOK, MIMETextPlainCharsetUTF8, []byte(fmt.Sprintf(format, a...)))
}

// Retorna un estatus 202 aceptado con el mensaje dado.
func (c *Context) StatusAccepted(msg string) error {
	return &Gkerror{codigo: http.StatusAccepted, mensaje: msg}
}

// ================================================================ //
// ========== Redirecciones (3xx) ================================= //

// Redirect the request to a provided URL with status code.
func (c *Context) Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		return ErrInvalidRedirectCode
	}
	c.response.Header().Set(HeaderLocation, url)
	c.response.WriteHeader(code)
	return nil
}

// Redirige a la URL usando fmt.Sprintf con código 303 TemporaryRedirect.
func (c *Context) Redir(format string, a ...any) error {
	c.response.Header().Set(HeaderLocation, fmt.Sprintf(format, a...))
	c.response.WriteHeader(303)
	return nil
}

// ================================================================ //
// ========== Errores del cliente (4xx) =========================== //

// Retorna un error 400 Bad Request.
func (c *Context) ErrBadRequest(err error) *Gkerror {
	return &Gkerror{codigo: http.StatusBadRequest, err: err}
}

// Retorna un error 404 Not Found.
func (c *Context) ErrNotFound(err error) *Gkerror {
	return &Gkerror{codigo: http.StatusNotFound, err: err}
}

// ================================================================ //

// Retorna un error 400 Bad Request.
func (c *Context) StatusBadRequest(msg string) *Gkerror {
	if msg == "" {
		msg = "Solicitud no aceptada"
	}
	return &Gkerror{codigo: http.StatusBadRequest, mensaje: msg}
}

// Retorna un error 401 Unauthorized, para usuario no autenticado.
func (c *Context) StatusUnauthorized(msg string) *Gkerror {
	if msg == "" {
		msg = "No autorizado"
	}
	return &Gkerror{codigo: http.StatusUnauthorized, mensaje: msg}
}

// Retorna un error 402 Payment Required.
func (c *Context) StatusPaymentRequired(msg string) *Gkerror {
	if msg == "" {
		msg = "Pago requerido"
	}
	return &Gkerror{codigo: http.StatusPaymentRequired, mensaje: msg}
}

// Retorna un error 403 Forbidden, para privilegios insuficientes.
func (c *Context) StatusForbidden(msg string) *Gkerror {
	if msg == "" {
		msg = "No permitido"
	}
	return &Gkerror{codigo: http.StatusForbidden, mensaje: msg}
}

// Retorna un error 404 Not Found.
func (c *Context) StatusNotFound(msg string) *Gkerror {
	if msg == "" {
		msg = "Recurso no encontrado"
	}
	return &Gkerror{codigo: http.StatusNotFound, mensaje: msg}
}

// Retorna un error 409 Conflict, para already exists u otros conflictos.
func (c *Context) StatusConflict(msg string) *Gkerror {
	if msg == "" {
		msg = "Conflicto con recurso existente"
	}
	return &Gkerror{codigo: http.StatusConflict, mensaje: msg}
}

// Retorna un error 415 Unsupported Media Type.
func (c *Context) StatusUnsupportedMedia(msg string) *Gkerror {
	if msg == "" {
		msg = "Tipo de media no soportado"
	}
	return &Gkerror{codigo: http.StatusUnsupportedMediaType, mensaje: msg}
}

// Retorna un error 429.
func (c *Context) StatusTooManyRequests(msg string) *Gkerror {
	if msg == "" {
		msg = "Demasiadas solicitudes"
	}
	return &Gkerror{codigo: http.StatusTooManyRequests, mensaje: msg}
}

// ================================================================ //
// ========== Errores del servidor (5xx) =========================== //

// Retorna un error 500.
func (c *Context) ServerError(err error) *Gkerror {
	return &Gkerror{codigo: http.StatusInternalServerError, err: err}
}

// ================================================================ //

// Retorna un error 500 Internal Server Error.
func (c *Context) StatusServerError(msg string) *Gkerror {
	if msg == "" {
		msg = "Error en servidor"
	}
	return &Gkerror{codigo: http.StatusInternalServerError, mensaje: msg}
}
