package plantillas

import (
	"io"
	"monorepo/gecko"
	"strings"
)

// Render satisface la interfaz gecko.Renderer
//
// Es lo último que se debe llamar en un handler.
//
// Ejecuta una plantilla previamente instanciada al crear el servicio.
//
// Si la plantilla no existe, responde con el error definido en NuevoServicio.
func (s *TemplateResponder) Render(w io.Writer, nombre string, data interface{}, c *gecko.Context) error {
	if s.reparse {
		s.ReParse()
	}
	if strings.HasSuffix(nombre, ".html") {
		c.LogWarnf("plantilla.Render: no es necesario poner .html a '%v'", nombre)
		nombre = strings.TrimSuffix(nombre, ".html")
	}
	return s.t.ExecuteTemplate(w, nombre, data)
}
