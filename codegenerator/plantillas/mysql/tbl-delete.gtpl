func (s *Repositorio) Delete{{ .Tabla.NombreItem }}({{ .Tabla.PrimaryKeysAsFuncParams }}) error {
	const op string = "Delete{{ .Tabla.NombreItem }}"
	{{ range .Tabla.PrimaryKeys }}{{ .IfZeroReturnErr "pk_indefinida" "" }}{{ end -}}
	err := s.Existe{{ .Tabla.NombreItem }}({{ .Tabla.PrimaryKeysAsArguments "" }})
	if err != nil {
		return gko.Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"DELETE FROM {{ .Tabla.NombreRepo }} {{ .Tabla.PrimaryKeysAsSqlWhere }}",
		{{ .Tabla.PrimaryKeysAsArguments "" }},
	)
	if err != nil {
		{{ if .Tabla.MySQL -}}
		if strings.HasPrefix(err.Error(),"Error 1451 (23000)"){
			return gko.ErrYaExiste().Err(err).Op(op).Msg("Este registro es referenciado por otros y no se puede eliminar")
		} else {
			return gko.ErrInesperado().Err(err).Op(op)
		}
		{{ else -}}
		return gko.ErrAlEscribir().Err(err).Op(op)
		{{- end }}
	}
	return nil
}