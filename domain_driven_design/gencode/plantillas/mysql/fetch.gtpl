{{ with $.TablaOrConsulta -}}
// Fetch{{ .NombreItem }} popula el {{ .NombreItem }} con el registro de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) Fetch{{ .NombreItem }}({{ .PrimaryKeysAsFuncParams }}, {{ .NombreAbrev }} *{{ .Paquete.Nombre }}.{{ .NombreItem }}) error {
	const op string = "mysql{{ .Paquete.Nombre }}.Fetch{{ .NombreItem }}"
	{{ range .PrimaryKeys }}{{ .IfZeroReturnErr "pk_indefinida" "" }}{{ end -}}
	row := s.db.QueryRow(
		"SELECT " + columnas{{ .NombreItem }} + " " + from{{ .NombreItem }} +
		"{{ .PrimaryKeysAsSqlWhere }}" {{ if .SqlGroupClause " " }}+ " " + group{{ .NombreItem }}{{ end }},
		{{ .PrimaryKeysAsArguments "" }},
	)
	return s.scanRow{{ .NombreItem }}(row, {{ .NombreAbrev }}, op)
}
{{- end }}