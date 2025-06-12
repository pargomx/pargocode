{{ with $.TablaOrConsulta -}}
func (s *Repositorio) List{{ .NombreItems }}({{ if .CamposFiltro }}filtros *Filtros{{ .NombreItem }}{{ end }}) ([]{{ .Paquete.Nombre }}.{{ .NombreItem }}, error) {
	const op string = "List{{ .NombreItems }}"
	{{ if .CamposFiltro -}}
	argumentos := []any{}
	where := ""
	if filtros != nil{		
		{{ range .CamposFiltro }}{{ template "filtro" . }}{{ end }}
		if len(where) > 4 {
			where = "WHERE"+where[4:] //Reemplazar " AND" por "WHERE "
		}
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
	{{ end -}}

	rows, err := s.db.Query(
		"SELECT " + columnas{{ .NombreItem }} + " " + from{{ .NombreItem }}
		{{- if .CamposFiltro }} + 
		where {{ if .SqlGroupClause " " }}+ " " + group{{ .NombreItem }}{{ end }} + limit,
		argumentos...
		{{- end }},
	)
	if err != nil {
		return nil, gko.ErrInesperado.Err(err).Op(op)
	}
	return s.scanRows{{ .NombreItem }}(rows, op)
}
{{- end }}