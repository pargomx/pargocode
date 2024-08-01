package gk{{$.Tabla.Nombre.Clave}}

import (

)

{{/* NUEVO y EDITAR: solo para el modelo principal editable */}}
// ================================================================ //
// ====== N U E V O =============================================== //

// Nuevo muestra el formulario para registrar un nuev{{$.Tabla.Nombre.LetraGen}} {{$.Tabla.Nombre.Humano}}.
func (s *handlersGecko) Nuevo(c *gecko.Context) error {
	data := map[string]any{
		"Sesion":       c.Sesion,
		"Titulo":       "{{$.Tabla.Nombre.Humano}} - nuev{{$.Tabla.Nombre.LetraGen}}",
		"LinkCerrar": "/inicio",

		"{{$.Tabla.Nombre.Camel}}":     {{$.Tabla.Nombre.Clave}}.{{$.Tabla.NombreItem}}{},
		{{ range $.Tabla.CamposEspeciales -}}
		"Lista{{.Nombre.Camel}}": {{$.Tabla.Nombre.Clave}}.Lista{{.Nombre.Camel}},
		{{ end }}
	}
	return c.RenderOk("{{$.Tabla.Nombre.ClavePlural}}/nuevo", data)
}

// ================================================================ //
// ====== E D I T A R ============================================= //

func (s *handlersGecko) Editar(c *gecko.Context) error {
	{{/* Obtener entidades foráneas si se necesitan sus PK para formar la PK de esta entidad */}}
	{{- range $.Tabla.PrimaryKeys }} {{- $eFK := (index $.TablasFK .Nombre.Camel) -}}
	{{ if len $eFK.Modelos -}}
	{{ $mFK := $eFK.PrimerModelo -}}
	{{ $cFK := index $mFK.UniqueKeys 0 }}
	{{$eFK.Nombre.Camel}}, err := s.{{$eFK.Nombre.ClavePlural}}.FetchItemBy{{$cFK.Nombre.Camel}}(c.PathVal("{{$eFK.Nombre.Clave}}"))
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{- end }}{{/* range $.Tabla.PrimaryKeys */}}
	{{- end }}{{/* if len $eFK.Modelos */}}

	{{/* Usar su clave única si la tiene */}}
	{{ if len $.Tabla.UniqueKeys }}
	{{ $cClave := index $.Tabla.UniqueKeys 0 -}}
	{{$.Tabla.Nombre.Camel}}, err := s.lector.Fetch{{$.Tabla.NombreItem}}By{{$cClave.Nombre.Camel}}(c.PathVal("{{$.Tabla.Nombre.Clave}}"))
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{/* O si solo tiene una sola PK y es numérica ... */}}
	{{ else if and (eq (len $.Tabla.PrimaryKeys) 1) ((index $.Tabla.PrimaryKeys 0).EsNumero) }}
	{{$.ParamsIDs}}, err := c.PathIntMust("{{$.Tabla.Nombre.Clave}}")
	if err != nil {
		return c.ErrBadRequest(err)
	}
	{{$.Tabla.Nombre.Camel}}, err := s.lector.Fetch{{$.Tabla.NombreItem}}ByID({{$.ParamsIDs}})
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{/* O si solo tiene una sola PK, probablemente string... */}}
	{{ else if eq (len $.Tabla.PrimaryKeys) 1 }}
	{{$.Tabla.Nombre.Camel}}, err := s.lector.Fetch{{$.Tabla.NombreItem}}ByID(c.PathVal("{{$.Tabla.Nombre.Clave}}"))
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{/* O si sabemos qué usar, dejar parámetros fáciles de popular. */}}
	{{ else }}
	{{$.Tabla.Nombre.Camel}}, err := s.lector.Fetch{{$.Tabla.NombreItem}}ByID{{if gt (len $.Tabla.PrimaryKeys) 1}}s{{end}}({{$.ParamsIDs}})
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{ end }}{{/* termina FetchItemBy... */}}

	{{/* FKs que no son PK */}}
	{{ range $.Tabla.ForeignKeys }}{{ if not .PrimaryKey }}
	{{ $eFK := (index $.TablasFK .Nombre.Camel) -}}
	{{ if len $eFK.Modelos -}}
	{{ $mFK := $eFK.PrimerModelo -}}
	{{ $cFK := index $mFK.PrimaryKeys 0 -}}
	{{$eFK.Nombre.Camel}}, err := s.{{$eFK.Nombre.ClavePlural}}.FetchItemByID({{$.Tabla.Nombre.Camel}}.{{$cFK.Nombre.Camel}})
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{- end }}{{- end }}
	{{- end }}

	data := map[string]any{
		"Sesion": c.Sesion,
		"Titulo": "{{$.Tabla.Nombre.Humano}} - editar",
		"LinkCerrar": "/inicio",

		"{{$.Tabla.Nombre.Camel}}": {{$.Tabla.Nombre.Camel}},
		{{ range $.TablasFK -}}
		"{{.Nombre.Camel}}": {{.Nombre.Camel}},
		{{end}}

		{{ range $.Tabla.CamposEspeciales -}}
		"Lista{{.Nombre.Camel}}": {{$.Tabla.Nombre.Clave}}.Lista{{.Nombre.Camel}},
		{{ end }}
	}
	return c.RenderOk("{{$.Tabla.Nombre.ClavePlural}}/editar", data)
}

// ================================================================ //
// ========== D E T A L L E S ===================================== //

func (s *handlersGecko) Detalles(c *gecko.Context) error {
	{{/* Obtener entidades foráneas si se necesitan sus PK para formar la PK de esta entidad */}}
	{{- range $.Tabla.PrimaryKeys }} {{- $eFK := (index $.TablasFK .Nombre.Camel) -}}
	{{ if len $eFK.Modelos -}}
	{{ $mFK := $eFK.PrimerModelo -}}
	{{ $cFK := index $mFK.UniqueKeys 0 }}
	{{$eFK.Nombre.Camel}}, err := s.{{$eFK.Nombre.ClavePlural}}.FetchItemBy{{$cFK.Nombre.Camel}}(c.PathVal("{{$eFK.Nombre.Clave}}"))
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{- end }}{{/* range $.Tabla.PrimaryKeys */}}
	{{- end }}{{/* if len $eFK.Modelos */}}

	{{/* Usar su clave única si la tiene */}}
	{{ if len $.Tabla.UniqueKeys }}
	{{ $cClave := index $.Tabla.UniqueKeys 0 -}}
	{{$.Tabla.Nombre.Camel}}, err := s.lector.Fetch{{$.Tabla.NombreItem}}By{{$cClave.Nombre.Camel}}(c.PathVal("{{$.Tabla.Nombre.Clave}}"))
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{/* O si solo tiene una sola PK y es numérica ... */}}
	{{ else if and (eq (len $.Tabla.PrimaryKeys) 1) ((index $.Tabla.PrimaryKeys 0).EsNumero) }}
	{{$.ParamsIDs}}, err := c.PathIntMust("{{$.Tabla.Nombre.Clave}}")
	if err != nil {
		return c.ErrBadRequest(err)
	}
	{{$.Tabla.Nombre.Camel}}, err := s.lector.Fetch{{$.Tabla.NombreItem}}ByID({{$.ParamsIDs}})
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{/* O si solo tiene una sola PK, probablemente string... */}}
	{{ else if eq (len $.Tabla.PrimaryKeys) 1 }}
	{{$.Tabla.Nombre.Camel}}, err := s.lector.Fetch{{$.Tabla.NombreItem}}ByID(c.PathVal("{{$.Tabla.Nombre.Clave}}"))
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{/* O si sabemos qué usar, dejar parámetros fáciles de popular. */}}
	{{ else }}
	{{$.Tabla.Nombre.Camel}}, err := s.lector.Fetch{{$.Tabla.NombreItem}}ByID{{if gt (len $.Tabla.PrimaryKeys) 1}}s{{end}}({{$.ParamsIDs}})
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{ end }}{{/* termina FetchItemBy... */}}

	{{/* FKs que no son PK */}}
	{{ range $.Tabla.ForeignKeys }}{{ if not .PrimaryKey }}
	{{ $eFK := (index $.TablasFK .Nombre.Camel) -}}
	{{ if len $eFK.Modelos -}}
	{{ $mFK := $eFK.PrimerModelo -}}
	{{ $cFK := index $mFK.PrimaryKeys 0 -}}
	{{$eFK.Nombre.Camel}}, err := s.{{$eFK.Nombre.ClavePlural}}.FetchItemByID({{$.Tabla.Nombre.Camel}}.{{$cFK.Nombre.Camel}})
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{- end }}{{- end }}
	{{- end }}

	data := map[string]any{
		"Sesion": c.Sesion,
		"Titulo": "{{$.Tabla.Nombre.Humano}} - detalles",
		"LinkCerrar": "/inicio",

		"{{$.Tabla.Nombre.Camel}}": {{$.Tabla.Nombre.Camel}},

		{{ range $.TablasFK -}}
		"{{.Nombre.Camel}}": {{.Nombre.Camel}},
		{{end -}}
	}
	return c.RenderOk("{{$.Tabla.Nombre.ClavePlural}}/detalles", data)
}

// ================================================================ //
// ========== L I S T A =========================================== //

{{ if not (len $.Tabla.CamposFiltro) }}
// Lista devuelve la lista de todos los registros existentes en el repositorio.
func (s *handlersGecko) Lista(c *gecko.Context) error {

	{{$.Tabla.Nombre.CamelPlural}}, err := s.lector.ListItemsAll()
	if err != nil {
		return c.ServerError(err)
	}

	data := map[string]any{
		"Sesion":     c.Sesion,
		"Titulo":     "{{$.Tabla.Nombre.HumanoPlural}} - tod{{$.Tabla.Nombre.LetraGen}}s",
		"LinkCerrar": "/inicio",

		"{{$.Tabla.Nombre.CamelPlural}}": {{$.Tabla.Nombre.CamelPlural}},
	}
	return c.RenderOk("{{$.Tabla.Nombre.ClavePlural}}/lista", data)
}
{{ end }}

{{ if (len $.Tabla.CamposFiltro) }}{{/* ListByFiltros */}}
// Listar aplica filtros y paginación, y si no se especificó nada
// entonces devuelve todos los registros.
func (s *handlersGecko) Listar(c *gecko.Context) error {

	filtros := {{$.Tabla.Nombre.Clave}}.FiltrosItem{}
	filtros.Limit = c.FormInt("ver")
	filtros.Offset = c.FormInt("pag") * filtros.Limit
	{{ range $.Tabla.CamposFiltro -}}
		{{ if .Especial -}}
		for _, v := range  {
			filtros.{{.Nombre.Camel}} = append(filtros.{{.Nombre.Camel}}, {{$.Tabla.Nombre.Clave}}.Set{{.Nombre.Camel}}Filtro(v))
		}
		{{ else -}}
		filtros.{{.Nombre.Camel}} = append(filtros.{{.Nombre.Camel}}, c.MultiQuery{{if eq .Tipo "int"}}Int{{else}}Val{{end}}("{{.NombreColumna}}")...)
		{{ end -}}
	{{ end }}

	{{$.Tabla.Nombre.CamelPlural}}, err := s.lector.ListItems(&filtros)
	if err != nil {
		if !strings.Contains(err.Error(), "filtros indefinidos") {
			return c.ServerError(err)
		}
		{{$.Tabla.Nombre.CamelPlural}}, err = s.lector.ListItemsAll()
		if err != nil {
			return c.ServerError(err)
		}
	}

	data := map[string]any{
		"Sesion":     c.Sesion,
		"Titulo":     "{{$.Tabla.Nombre.HumanoPlural}}",
		"LinkCerrar": "/inicio",

		"{{$.Tabla.Nombre.CamelPlural}}": {{$.Tabla.Nombre.CamelPlural}},
		
		{{ range $.Tabla.CamposEspeciales -}}
		"Lista{{.Nombre.Camel}}": {{$.Tabla.Nombre.Clave}}.ListaFiltro{{.Nombre.Camel}},
		"Filtro{{.Nombre.Camel}}": filtros.Primer{{.Nombre.Camel}}(),
		{{ end -}}
	}
	return c.RenderOk("{{$.Tabla.Nombre.ClavePlural}}/lista", data)
}
{{ end }}{{/* ListByFiltros */}}

{{/* --------- List By FK --------------- */}}
{{- range $.Tabla.ForeignKeys }}
{{ $eFK := (index $.TablasFK .Nombre.Camel) -}}
{{ if len $eFK.Modelos -}}
{{ $mFK := $eFK.PrimerModelo -}}
{{ $cFK := index $mFK.PrimaryKeys 0 -}}{{ if (len $mFK.UniqueKeys) -}}{{ $cFK = index $mFK.UniqueKeys 0 -}}{{ end }}
// ================================================================ //{{br}}
// ListaBy{{$eFK.Nombre.Camel}} retorna todos los registros pertenecientes
// {{$eFK.Nombre.Acusativo}} especificado como parámetro URL: "/:{{$eFK.Nombre.Clave}}/"
func (s *handlersGecko) ListaBy{{$eFK.Nombre.Camel}}(c *gecko.Context) error {
	filtros := {{$.Tabla.Nombre.Clave}}.FiltrosItem{}
	{{ range $.Tabla.CamposFiltro -}}
	for _, v := range c.MultiQuery{{if eq .Tipo "int"}}Int{{else}}Val{{end}}("{{.NombreColumna}}") {
		{{ if .Especial -}}
		filtros.{{.Nombre.Camel}} = append(filtros.{{.Nombre.Camel}}, {{$.Tabla.Nombre.Clave}}.Set{{.Nombre.Camel}}Filtro(v))
		{{ else -}}
		filtros.{{.Nombre.Camel}} = append(filtros.{{.Nombre.Camel}}, v)
		{{ end -}}
	}
	{{ end }}

	{{$eFK.Nombre.Camel}}, err := s.{{$eFK.Nombre.ClavePlural}}.FetchItemBy{{$cFK.Nombre.Camel}}(c.PathVal("{{$eFK.Nombre.Clave}}"))
	if err != nil {
		return c.ErrNotFound(err)
	}

	{{$.Tabla.Nombre.CamelPlural}}, err := s.lector.ListItemsBy{{$mFK.PrimerCampo.Nombre.Camel}}({{$eFK.Nombre.Camel}}.{{$mFK.PrimerCampo.Nombre.Camel}}, &filtros)
	if err != nil {
		return c.ServerError(err)
	}

	data := map[string]any{
		"Sesion":     c.Sesion,
		"Titulo":     "{{$.Tabla.Nombre.HumanoPlural}}",
		"LinkCerrar": "/inicio",

		"{{$eFK.Nombre.Camel}}": {{$eFK.Nombre.Camel}},
		"{{$.Tabla.Nombre.CamelPlural}}": {{$.Tabla.Nombre.CamelPlural}},
		
		{{ range $.Tabla.CamposEspeciales -}}
		"Lista{{.Nombre.Camel}}": {{$.Tabla.Nombre.Clave}}.ListaFiltro{{.Nombre.Camel}},
		"Filtro{{.Nombre.Camel}}": filtros.Primer{{.Nombre.Camel}}(),
		{{ end -}}
	}
	return c.RenderOk("{{$.Tabla.Nombre.ClavePlural}}/lista", data)
}
{{ end }}
{{ end }}


