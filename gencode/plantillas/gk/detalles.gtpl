// ================================================================ //
// ========== D E T A L L E S ===================================== //

func (s *handlersGecko) Detalles(c *gecko.Context) error {
	{{/* Obtener item por PKs */ -}}
	{{$.Tabla.Nombre.Camel}}, err := s.{{$.Tabla.Nombre.ClavePlural}}.Fetch{{$.Tabla.NombreItem}}ByID({{ range $.Tabla.PrimaryKeys }}c.Path{{ .TipoGeckoGet }}("{{ .NombreColumna }}"),{{ end }})
	if gecko.EsErrNotFound(err) {
		return c.Redir("/{{$.Tabla.Nombre.ClavePlural}}{{ range $.Tabla.PrimaryKeys }}/%v{{ end }}/{{ $.Tabla.Nombre.Kebab }}/form",{{ range $.Tabla.PrimaryKeys }} c.Path{{ .TipoGeckoGet }}("{{ .NombreColumna }}"),{{ end }})
	} else if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo": "{{$.Tabla.Nombre.Humano}} - detalles",
		"{{$.Tabla.Nombre.Camel}}": {{$.Tabla.Nombre.Camel}},
	}
	return c.RenderOk("{{$.Tabla.Nombre.ClavePlural}}/detalles", data)
}

// ================================================================ //
// ========== D E T A L L E S ===================================== //

func (s *handlersGecko) DetallesMust(c *gecko.Context) error {
	{{/* Obtener item por PKs */ -}}
	{{$.Tabla.Nombre.Camel}}, err := s.{{$.Tabla.Nombre.ClavePlural}}.Fetch{{$.Tabla.NombreItem}}ByID({{ range $.Tabla.PrimaryKeys }}c.Path{{ .TipoGeckoGet }}("{{ .NombreColumna }}"),{{ end }})
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo": "{{$.Tabla.Nombre.Humano}} - detalles",
		"{{$.Tabla.Nombre.Camel}}": {{$.Tabla.Nombre.Camel}},
	}
	return c.RenderOk("{{$.Tabla.Nombre.ClavePlural}}/detalles", data)
}

