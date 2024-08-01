{{ with $.TablaOrConsulta -}}
// Get{{ .NombreItem }} devuelve un {{ .NombreItem }} de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) Get{{ .NombreItem }}({{ .PrimaryKeysAsFuncParams }}) (*{{ .Paquete.Nombre }}.{{ .NombreItem }}, error) {
	const op string = "Get{{ .NombreItem }}"
	{{ range .PrimaryKeys }}{{ .IfZeroReturnNilAndErr "pk_indefinida" "" }}{{ end -}}
	row := s.db.QueryRow(
		"SELECT " + columnas{{ .NombreItem }} + " " + from{{ .NombreItem }} +
		"{{ .PrimaryKeysAsSqlWhere }}" {{ if .SqlGroupClause " " }}+ " " + group{{ .NombreItem }}{{ end }},
		{{ .PrimaryKeysAsArguments "" }},
	)
	{{ .NombreAbrev }} := &{{ .Paquete.Nombre }}.{{ .NombreItem }}{}
	err := s.scanRow{{ .NombreItem }}(row, {{ .NombreAbrev }})
	if err != nil {
		return nil, err
	}
	return {{ .NombreAbrev }}, nil
}
{{- end }}