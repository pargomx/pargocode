// Filtros se utiliza para listar {{ .Consulta.NombreItem }} e indica
// los parámetros que deben tener los resultados de la consulta.
type Filtros{{ .Consulta.NombreItem }} struct {
	{{- range .Consulta.CamposFiltro }}
	{{ .NombreCampo }} {{ if .EsBool }}*{{ else }}[]{{ end -}}
	{{- .TipoGo }}
	{{- end }}

	Limit  int // Se limita el número de registros devueltos si es mayor a 0.
	Offset int // Utilizado para paginación. Limit debe ser mayor a 0.
}
