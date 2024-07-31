{{ with $.TablaOrConsulta -}}
// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRow{{ .NombreItem }}(row *sql.Row, {{ .NombreAbrev }} *{{ .Paquete.Nombre }}.{{ .NombreItem }}, op string) error {
	{{ .ScanTempVars }}
	err := row.Scan(
		{{ .ScanArgs }},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado().Msg("{{ .TablaOrigen.Humano }} no se encuentra").Op(op)
		}
		return gko.ErrInesperado().Err(err).Op(op)
	}
	{{ .ScanSetters }}
	return nil
}
{{- end }}
