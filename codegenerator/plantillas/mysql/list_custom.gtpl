{{ with $.TablaOrConsulta -}}
func (s *Repositorio) List{{ .NombreItems }}{{ .CustomList.CompFunc }}({{ .CustomList.ArgsFunc }}) ([]{{ .Paquete.Nombre }}.{{ .NombreItem }}, error) {
	const op string = "List{{ .NombreItems }}{{ .CustomList.CompFunc }}"
	rows, err := s.db.Query(
		"SELECT " + columnas{{ .NombreItem }} + " " + from{{ .NombreItem }}+
			"{{ .CustomList.CompSQL }}",
		{{ if .CustomList.ArgsSQL }}{{ .CustomList.ArgsSQL }},{{ end }}
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRows{{ .NombreItem }}(rows, op)
}
{{- end }}