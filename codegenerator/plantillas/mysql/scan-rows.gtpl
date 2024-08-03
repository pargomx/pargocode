{{ with $.TablaOrConsulta -}}
// scanRows{{ .NombreItem }} escanea cada row en la struct {{ .NombreItem }}
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRows{{ .NombreItem }}(rows *sql.Rows, op string) ([]{{ .Paquete.Nombre }}.{{ .NombreItem }}, error) {
	defer rows.Close()
	items := []{{ .Paquete.Nombre }}.{{ .NombreItem }}{}
	for rows.Next() {
		{{ .NombreAbrev }} := {{ .Paquete.Nombre }}.{{ .NombreItem }}{}
		{{- .ScanTempVars }}
		err := rows.Scan(
			{{ .ScanArgs }},
		)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}
		{{- .ScanSetters }}
		items = append(items, {{ .NombreAbrev }})
	}
	return items, nil
}
{{- end }}
