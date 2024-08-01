{{ with $.TablaOrConsulta -}}
{{ if .CamposFiltro -}}
func (s *Repositorio) List{{ .NombreItems }}By{{ range .CamposSeleccionados }}{{ .NombreCampo }}{{ end }}({{ .CamposSeleccionadosAsFuncParams }}, filtros *Filtros{{ .NombreItem }}) ([]{{ .Paquete.Nombre }}.{{ .NombreItem }}, error) {
	const op string = "mysql{{ .Paquete.Nombre }}.List{{ .NombreItems }}By{{ range .CamposSeleccionados }}{{ .NombreCampo }}{{ end }}"
	{{ range .CamposSeleccionados }}{{ .IfZeroReturnNilAndErr "param_indefinido" "" }}{{ end -}}
	where := "{{ .CamposSeleccionadosAsSqlWhere }}"
	argumentos := []any{}
	argumentos = append(argumentos, {{ .CamposSeleccionadosAsArguments "" }})
	if filtros != nil{		
		{{ range .CamposFiltro }}{{ template "filtro" . }}{{ end }}
	}

	limit := ""
	if filtros != nil{
		if filtros.Limit > 0  && filtros.Offset == 0{
			limit = " LIMIT ?"
			argumentos = append(argumentos, filtros.Limit)
		} else if filtros.Limit > 0  && filtros.Offset > 0{
			limit = " LIMIT ? OFFSET ?"
			argumentos = append(argumentos, filtros.Limit, filtros.Offset)
		}
	}

	rows, err := s.db.Query(
		"SELECT " + columnas{{ .NombreItem }} + " " + from{{ .NombreItem }} +
		where {{ if .SqlGroupClause " " }}+ " " + group{{ .NombreItem }}{{ end }} + limit,
		argumentos...,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRows{{ .NombreItem }}(rows, op)
}
{{ else -}}
func (s *Repositorio) List{{ .NombreItems }}By{{ range .CamposSeleccionados }}{{ .NombreCampo }}{{ end }}({{ .CamposSeleccionadosAsFuncParams }}) ([]{{ .Paquete.Nombre }}.{{ .NombreItem }}, error) {
	const op string = "mysql{{ .Paquete.Nombre }}.List{{ .NombreItems }}By{{ range .CamposSeleccionados }}{{ .NombreCampo }}{{ end }}"
	{{ range .CamposSeleccionados }}{{ .IfZeroReturnNilAndErr "param_indefinido" "" }}{{ end -}}
	rows, err := s.db.Query(
		"SELECT " + columnas{{ .NombreItem }} + " " + from{{ .NombreItem }} +
		"{{ .CamposSeleccionadosAsSqlWhere }}" {{ if .SqlGroupClause " " }}+ " " + group{{ .NombreItem }}{{ end }},
		{{ .CamposSeleccionadosAsArguments "" }},
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRows{{ .NombreItem }}(rows, op)
}
{{ end -}}
{{ end -}}
