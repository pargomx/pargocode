// HandlersGecko para {{ lower $.Tabla.Nombre.HumanoPlural }}.
//
//  {{ $.Tabla.Nombre.ClavePlural }} := gk{{ $.Tabla.Nombre.Clave }}.NewHandlers({{ $.Tabla.Nombre.ClavePlural }}DB, {{ range .EntidadesFK }}{{ .Nombre.ClavePlural }}DB{{- end -}})

package gk{{ $.Tabla.Nombre.Clave }}

type handlersGecko struct {
	{{ $.Tabla.Nombre.ClavePlural }} *mysql{{$.Tabla.Nombre.Clave}}.Repositorio
	{{ range .EntidadesFK -}}
	{{ .Nombre.ClavePlural }} *mysql{{.Nombre.Clave}}.Repositorio
	{{ end }}
}

// HandlersGecko para {{ lower $.Tabla.Nombre.HumanoPlural }}.
//
//  {{ $.Tabla.Nombre.ClavePlural }} := gk{{ $.Tabla.Nombre.Clave }}.NewHandlers({{ $.Tabla.Nombre.ClavePlural }}DB, {{ range .EntidadesFK }}{{ .Nombre.ClavePlural }}DB{{- end -}})
func NewHandlers(
	{{ $.Tabla.Nombre.ClavePlural }} *mysql{{$.Tabla.Nombre.Clave}}.Repositorio,
	{{ range .EntidadesFK }}{{ .Nombre.ClavePlural }} *mysql{{.Nombre.Clave}}.Repositorio,
	{{ end }}
) *handlersGecko {
	if {{ $.Tabla.Nombre.ClavePlural }} == nil { {{/* comprobar repositorios nil */}}
		gecko.FatalFmt("gk{{ $.Tabla.Nombre.Clave }}.NuevoServicio", errors.New("{{$.Tabla.Nombre.Clave}}.Repositorio es nil"))
	}
	{{ range .EntidadesFK -}}
	if {{ .Nombre.ClavePlural }} == nil {
		gecko.FatalFmt("gk{{ $.Tabla.Nombre.Clave }}.NuevoServicio", errors.New("{{ .Nombre.ClavePlural }}.Lector es nil"))
	}
	{{ end -}}
	return &handlersGecko{
		{{ $.Tabla.Nombre.ClavePlural }}: {{ $.Tabla.Nombre.ClavePlural }},
		{{ range .EntidadesFK -}}
		{{ .Nombre.ClavePlural }}: {{ .Nombre.ClavePlural }},
		{{ end }}
	}
}
