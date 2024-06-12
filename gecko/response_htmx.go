package gecko

import "fmt"

// ================================================================ //
// ========== Request HTMX ======================================== //

// Si la solicitud viene de HTMX significa que tiene el header HX-Request = true.
// Cuando es HX-History-Restore-Request se necesita enviar la página entera.
func (c *Context) EsHTMX() bool {
	return c.Request().Header.Get("HX-Request") == "true" &&
		c.Request().Header.Get("HX-History-Restore-Request") != "true"
}

// ================================================================ //
// ========== Responder HTMX ====================================== //

// Devuelve un estatus "204 No Content" e instruye a HTMX para que
// vuelva a cargar la página entera con el header "HX-Refresh".
//
// Conveniente como respuesta a una solicitud PUT.
func (c *Context) RefreshHTMX() error {
	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(204)
}

// Instruye al cliente para redirigir a la nueva ubicación.
//
// Responde un 200 OK: Redirigiendo a... y el header HX-Redirect.
//
// Utiliza fmt.Sprintf para construir el path.
func (c *Context) RedirectHTMX(path string, a ...any) error {
	c.Response().Header().Set("HX-Redirect", fmt.Sprintf(path, a...))
	return c.StatusOk("Redirigiendo a " + fmt.Sprintf(path, a...))
}

// Agrega un evento al HX-Trigger
func (c *Context) TriggerEventoHTMX(evento string) {
	c.Response().Header().Set("HX-Trigger", evento)
}
