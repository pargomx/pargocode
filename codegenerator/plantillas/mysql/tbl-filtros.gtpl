// Filtros se utiliza para listar {{ .Tabla.NombreItem }} e indica
// los parámetros que deben tener los resultados de la consulta.
type Filtros{{ .Tabla.NombreItem }} struct {
	{{- range .Tabla.CamposFiltro }}
	{{ .NombreCampo }} {{ if .EsBool }}*{{ else }}[]{{ end }}{{ if .Especial }}{{ .Paquete.Nombre }}.{{ end }}{{ .TipoGo }}
	{{- end }}

	Limit  int // Se limita el número de registros devueltos si es mayor a 0.
	Offset int // Utilizado para paginación. Limit debe ser mayor a 0.
}

{{ range .Tabla.CamposFiltro }}{{ if .Especial }}
// Retorna el primer filtro especificado, o bien Todos si no hay ninguno.
func (f Filtros{{ $.Tabla.NombreItem }}) Primer{{ .NombreCampo }}() {{ if .Especial }}{{ .Paquete.Nombre }}.{{ end }}{{ .TipoGo }} {
	if len(f.{{ .NombreCampo }}) == 0 {
		return {{ if .Especial }}{{ .Paquete.Nombre }}.{{ end }}{{ .TipoGo }}Todos
	} else {
		return f.{{ .NombreCampo }}[0]
	}
}
{{- end }}{{ end }}