// Filtros se utiliza para listar {{ .Consulta.NombreItem }} e indica
// los parámetros que deben tener los resultados de la consulta.
type Filtros{{ .Consulta.NombreItem }} struct {
	{{- range .Consulta.CamposFiltro }}
	{{ .NombreCampo }} {{ if .EsBool }}*{{ else }}[]{{ end -}}
	{{- if .Especial }}{{ .Tabla.Paquete.Nombre }}.{{ end -}}
	{{- .Tipo }}
	{{- end }}

	Limit  int // Se limita el número de registros devueltos si es mayor a 0.
	Offset int // Utilizado para paginación. Limit debe ser mayor a 0.
}

{{ range .Consulta.CamposFiltro }}{{ if .Especial }}
// Retorna el primer filtro especificado, o bien Todos si no hay ninguno.
func (f Filtros{{ $.Consulta.NombreItem }}) Primer{{ .NombreCampo }}() {{ if .Especial }}{{ .Tabla.Paquete.Nombre }}.{{ end }}{{ .Tipo }} {
	if len(f.{{ .NombreCampo }}) == 0 {
		return {{ if .Especial }}{{ .Tabla.Paquete.Nombre }}.{{ end }}{{ .Tipo }}Todos
	} else {
		return f.{{ .NombreCampo }}[0]
	}
}
{{- end }}{{ end }}