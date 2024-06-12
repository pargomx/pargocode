package gecko

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
)

// Renderer is the interface that wraps the Render function.
type Renderer interface {
	Render(io.Writer, string, interface{}, *Context) error
}

// Renderizar una plantilla registrada en gecko.Renderer bajo "name"
// y responder con MIME "text/html" y el status "code".
func (c *Context) Render(code int, name string, data interface{}) error {
	if c.gecko.Renderer == nil {
		return ErrRendererNotRegistered
	}
	buf := new(bytes.Buffer)
	err := c.gecko.Renderer.Render(buf, name, data, c)
	if err != nil {
		return err
	}
	return c.HTMLBlob(code, buf.Bytes())
}

// Renderizar una plantilla registrada en gecko.Renderer bajo "name"
// y responder con MIME "text/html" y el status "200 OK".
//
// Si la solicitud es HTMX manda la plantilla tal cual, sino le agrega
// el layout para ser una página HTML completa.
func (c *Context) RenderOk(name string, data map[string]any) error {
	if c.gecko.Renderer == nil {
		return ErrRendererNotRegistered
	}
	if data == nil {
		data = map[string]any{}
	}
	data["Sesion"] = c.Sesion

	if c.EsHTMX() { // Enviar solo parcial a HTMX

		data["EsHTMX"] = true
		buf := new(bytes.Buffer)
		err := c.gecko.Renderer.Render(buf, name, data, c)
		if err != nil {
			return err
		}
		c.response.Header().Add("Cache-Control", "no-store") // No guardar en ningún caché. HTMX se encarga con hx-push-url.
		return c.HTMLBlob(http.StatusOK, buf.Bytes())

	} else { // Enviar encapsulado en layout HTML a navegador.

		buf := new(bytes.Buffer)
		err := c.gecko.Renderer.Render(buf, name, data, c)
		if err != nil {
			return err
		}
		data["Contenido"] = template.HTML(buf.String())
		buf2 := new(bytes.Buffer)
		err = c.gecko.Renderer.Render(buf2, c.gecko.TmplBaseLayout, data, c)
		if err != nil {
			return err
		}
		return c.HTMLBlob(http.StatusOK, buf2.Bytes())
	}
}

// ================================================================ //
// ========== MAIN CONTENT ======================================== //

// Renderiza en #maincontent para htmx y en .Contenido par navegador.
func (c *Context) RenderContenido(name string, data map[string]any) error {
	if c.gecko.Renderer == nil {
		return ErrRendererNotRegistered
	}
	if data == nil {
		data = map[string]any{}
	}
	data["Sesion"] = c.Sesion

	if c.EsHTMX() { // Enviar solo parcial a HTMX

		data["EsHTMX"] = true
		buf := new(bytes.Buffer)
		err := c.gecko.Renderer.Render(buf, name, data, c)
		if err != nil {
			return err
		}
		c.response.Header().Add("Cache-Control", "no-store") // No guardar en ningún caché. HTMX se encarga con hx-push-url.
		c.response.Header().Add("HX-Retarget", "#contenido")
		return c.HTMLBlob(http.StatusOK, buf.Bytes())

	} else { // Enviar encapsulado en layout HTML a navegador.

		buf := new(bytes.Buffer)
		err := c.gecko.Renderer.Render(buf, name, data, c)
		if err != nil {
			return err
		}
		data["Contenido"] = template.HTML(buf.String())
		buf2 := new(bytes.Buffer)
		err = c.gecko.Renderer.Render(buf2, c.gecko.TmplBaseLayout, data, c)
		if err != nil {
			return err
		}
		return c.HTMLBlob(http.StatusOK, buf2.Bytes())
	}
}

// ================================================================ //
// ========== CARD ================================================ //

func (c *Context) RenderCard(name string, data map[string]any) error {
	if c.gecko.Renderer == nil {
		return ErrRendererNotRegistered
	}
	if data == nil {
		data = map[string]any{}
	}
	data["Sesion"] = c.Sesion
	c.response.Header().Add("Cache-Control", "no-store") // ningún caché

	if c.EsHTMX() { //* Enviar solo parcial a HTMX

		data["EsHTMX"] = true
		buf := new(bytes.Buffer)
		err := c.gecko.Renderer.Render(buf, name, data, c)
		if err != nil {
			return err
		}

		_, esWorkcard := data["MainCardURL"]
		esMaincard := !esWorkcard
		reqFromWorkcard := c.request.Header.Get("Hx-Target") == "workcard"

		// Si es workcard y viene de workard pero de otra maincard, solicitar nueva maincard.
		if c.request.Header.Get("HX-GetMaincard") == "true" && esWorkcard {
			// c.LogInfo("workcard con new maincard")
			data["CardBody"] = template.HTML(buf.String())
			buf2 := new(bytes.Buffer)
			err = c.gecko.Renderer.Render(buf2, "workcard-con-maincard", data, c)
			if err != nil {
				return err
			}
			c.response.Header().Add("HX-Retarget", "#contenido") // reemplazar ambas tarjetas
			return c.HTMLBlob(http.StatusOK, buf2.Bytes())
		}

		// Si es maincard solicitada desde workcard entonces limpiar workcard.
		if esMaincard && reqFromWorkcard {
			c.response.Header().Add("HX-Retarget", "#maincard")                              // reemplazar maincard
			buf.WriteString("<section id=\"workcard\" hx-swap-oob=\"innerHTML\"></section>") // limpiar workcard
		}

		// Si viene de contenido y el target es contenido, requiere wrap en card.
		reqFromContainer := c.request.Header.Get("Hx-Target") == "contenido"
		if reqFromContainer {
			// c.LogInfo("ReqFromContainer")
			data["CardBody"] = template.HTML(buf.String())
			buf2 := new(bytes.Buffer)
			err = c.gecko.Renderer.Render(buf2, "contenido", data, c)
			if err != nil {
				return err
			}
			return c.HTMLBlob(http.StatusOK, buf2.Bytes())
		}

		return c.HTMLBlob(http.StatusOK, buf.Bytes())

	} else { //* Enviar HTML completo al navegador.

		buf := new(bytes.Buffer)
		err := c.gecko.Renderer.Render(buf, name, data, c)
		if err != nil {
			return err
		}
		data["CardBody"] = template.HTML(buf.String())
		buf2 := new(bytes.Buffer)
		err = c.gecko.Renderer.Render(buf2, c.gecko.TmplBaseLayout, data, c)
		if err != nil {
			return err
		}
		return c.HTMLBlob(http.StatusOK, buf2.Bytes())
	}
}

// ================================================================ //
// ================================================================ //
