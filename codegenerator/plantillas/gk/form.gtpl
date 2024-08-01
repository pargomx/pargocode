// ================================================================ //
// ========== F O R M U L A R I O ================================= //

func (s *handlersGecko) Formulario(c *gecko.Context) error {
	{{/* Obtener item por PKs */ -}}
	{{$.Tabla.Nombre.Camel}}, err := s.{{$.Tabla.Nombre.ClavePlural}}.Fetch{{$.Tabla.NombreItem}}ByID({{ range $.Tabla.PrimaryKeys }}c.Path{{ .TipoGeckoGet }}("{{ .NombreColumna }}"),{{ end }})
	if gecko.EsErrNotFound(err) {
		{{ $.Tabla.Nombre.Camel }} = &{{ $.Tabla.Nombre.Clave }}.{{ $.Tabla.NombreItem }}{
			{{ range $.Tabla.PrimaryKeys -}}
			{{ .Nombre.Camel }}: c.Path{{ .TipoGeckoGet }}("{{ .NombreColumna }}"),
			{{ end }}
		}
	} else if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo": "{{$.Tabla.Nombre.Humano}} - form",
		"{{$.Tabla.Nombre.Camel}}": {{$.Tabla.Nombre.Camel}},
		{{ range $.Tabla.CamposEspeciales -}}
		"Lista{{.Nombre.Camel}}": {{$.Tabla.Nombre.Clave}}.Lista{{.Nombre.Camel}},
		{{ end }}
	}
	return c.RenderOk("{{$.Tabla.Nombre.ClavePlural}}/form", data)
}

// ================================================================ //
// ========== F O R M U L A R I O ================================= //

func (s *handlersGecko) Formulario(c *gecko.Context) error {
	{{/* Obtener item por PKs */ -}}
	{{$.Tabla.Nombre.Camel}}, err := s.{{$.Tabla.Nombre.ClavePlural}}.Fetch{{$.Tabla.NombreItem}}ByID({{ range $.Tabla.PrimaryKeys }}c.Path{{ .TipoGeckoGet }}("{{ .NombreColumna }}"),{{ end }})
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo": "{{$.Tabla.Nombre.Humano}} - form",
		"{{$.Tabla.Nombre.Camel}}": {{$.Tabla.Nombre.Camel}},
		{{ range $.Tabla.CamposEspeciales -}}
		"Lista{{.Nombre.Camel}}": {{$.Tabla.Nombre.Clave}}.Lista{{.Nombre.Camel}},
		{{ end }}
	}
	return c.RenderOk("{{$.Tabla.Nombre.ClavePlural}}/form", data)
}
