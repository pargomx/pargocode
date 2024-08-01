{{ with $.TablaOrConsulta -}}
// Get{{ .NombreItem }}By{{ range .CamposSeleccionados }}{{ .NombreCampo }}{{ end }} devuelve un {{ .NombreItem }} de la DB.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) Get{{ .NombreItem }}By{{ range .CamposSeleccionados }}{{ .NombreCampo }}{{ end }}({{ .CamposSeleccionadosAsFuncParams }}) (*{{ .Paquete.Nombre }}.{{ .NombreItem }}, error) {
	const op string = "mysql{{ .Paquete.Nombre }}.Get{{ .NombreItem }}By{{ range .CamposSeleccionados }}{{ .NombreCampo }}{{ end }}"
	{{ range .CamposSeleccionados }}{{ .IfZeroReturnNilAndErr "param_indefinido" "" }}{{ end -}}
	row := s.db.QueryRow(
		"SELECT " + columnas{{ .NombreItem }} + " " + from{{ .NombreItem }} +
		"{{ .CamposSeleccionadosAsSqlWhere }}" {{ if .SqlGroupClause " " }}+ " " + group{{ .NombreItem }}{{ end }},
		{{ .CamposSeleccionadosAsArguments "" }},
	)
	{{ .NombreAbrev }} := &{{ .Paquete.Nombre }}.{{ .NombreItem }}{}
	return {{ .NombreAbrev }}, s.scanRow{{ .NombreItem }}(row, {{ .NombreAbrev }}, op)
}
{{- end }}