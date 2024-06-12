// ================================================================ //
// ========== I N S E R T / U P D A T E============================ //

// Insertar el recurso o actualizarlo si ya existe.
func (s *handlersGecko) InsertUpdate(c *gecko.Context) error {
	{{/* Colocar nuevo objeto vacío */ -}}
	{{ $.Tabla.Nombre.Camel }} := {{ $.Tabla.Nombre.Clave }}.{{ $.Tabla.NombreItem }}{
		{{ range $.Tabla.PrimaryKeys -}}
		{{ .Nombre.Camel }}: c.Path{{ .TipoGeckoGet }}("{{ .NombreColumna }}"),
		{{ end }}
	}
	{{/* Obtener info existente si la hay */ -}}
	err := s.{{$.Tabla.Nombre.ClavePlural}}.Get{{$.Tabla.NombreItem}}ByID({{ range $.Tabla.PrimaryKeys }}{{ $.Tabla.Nombre.Camel }}.{{ .Nombre.Camel }},{{ end }} &{{ $.Tabla.Nombre.Camel }})
	if err != nil && !gecko.EsErrNotFound(err) {
		return gecko.Err(err).Op("gk{{ $.Tabla.Nombre.Clave }}.Put")
	}

	{{/* Poner nuevos datos */ -}}
{{ range $.Tabla.CamposEditables -}}
	
	{{ if eq .Tipo "string" -}}
	{{$.Tabla.Nombre.Camel}}.{{.Nombre.Camel}} = c.FormVal("{{.NombreColumna}}")
	{{/* */ -}}
	
	{{ else if eq .Tipo "int" -}}
	{{$.Tabla.Nombre.Camel}}.{{.Nombre.Camel}}, err = c.FormIntMust("{{.NombreColumna}}")
	if err != nil {
		return err
	}
	{{/* */ -}}

	{{ else if eq .Tipo "bool" -}}
	{{$.Tabla.Nombre.Camel}}.{{.Nombre.Camel}} = c.FormBool("{{.NombreColumna}}")
	{{/* */ -}}

	{{ else if eq .Tipo "*time.Time" -}}
	{{$.Tabla.Nombre.Camel}}.{{.Nombre.Camel}}, err = c.FormTime{{if .Null}}Nullable{{end}}("{{.NombreColumna}}", {{if eq .TimeTipo "date"}}"2006-01-02"{{else}}"????"{{end}})
	if err != nil {
		return err
	}
	{{/* */ -}}

	{{ else if .Especial -}}
	{{$.Tabla.Nombre.Camel}}.{{.Nombre.Camel}} = {{$.Tabla.Nombre.Clave}}.Set{{.Nombre.Camel}}DB(c.FormVal("{{.NombreColumna}}"))
	{{/* */ -}}

	{{ else -}}
	// {{$.Tabla.Nombre.Camel}}.{{.Nombre.Camel}} = c.FormVal("{{.NombreColumna}}") ???
	{{/* */ -}}

	{{ end -}}
{{ end }}{{/* range $.Tabla.CamposEditables */}}

	{{/* Insertar o actualizar item en repositorio */ -}}
	err = s.{{$.Tabla.Nombre.ClavePlural}}.InsertUpdateItem({{$.Tabla.Nombre.Camel}})
	if err != nil {
		return gecko.Err(err).Op("gk{{ $.Tabla.Nombre.Clave }}.Put")
	}
	return c.Redir("/{{$.Tabla.Nombre.ClavePlural}}{{ range $.Tabla.PrimaryKeys }}/%v{{ end }}/{{ $.Tabla.Nombre.Kebab }}", {{ range $.Tabla.PrimaryKeys }}{{ $.Tabla.Nombre.Camel }}.{{ .Nombre.Camel }},{{ end }})
}

