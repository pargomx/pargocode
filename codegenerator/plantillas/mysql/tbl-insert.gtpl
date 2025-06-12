func (s *Repositorio) Insert{{ .Tabla.NombreItem }}({{ .Tabla.NombreAbrev }} {{ .Tabla.Paquete.Nombre }}.{{ .Tabla.NombreItem }}) error {
	const op string = "Insert{{ .Tabla.NombreItem }}"
	{{ range .Tabla.CamposRequeridosOrPK -}}
		{{ if .PrimaryKey }}{{ .IfZeroReturnErr "pk_indefinida" $.Tabla.NombreAbrev -}}
		{{ else }}{{ .IfZeroReturnErr "required_sin_valor" $.Tabla.NombreAbrev }}{{ end -}}
	{{ end -}}
	_, err := s.db.Exec("INSERT INTO {{ .Tabla.NombreRepo }} "+
		"({{ .Tabla.CamposEditablesAsSnakeList ", " }}) "+
		"VALUES ({{ .Tabla.CamposEditablesAsPlaceholders }}) ",
		{{ .Tabla.CamposEditablesAsArguments .Tabla.NombreAbrev }},
	)
	if err != nil {
		{{ if .Tabla.MySQL -}}
		if strings.HasPrefix(err.Error(),"Error 1062 (23000)"){
			return gko.ErrYaExiste.Err(err).Op(op)
		} else if strings.HasPrefix(err.Error(), "Error 1452 (23000)") {
			return gko.ErrDatoInvalido.Err(err).Op(op).Msg("No se puede insertar la información porque el registro asociado no existe")
		} else {
			return gko.ErrInesperado.Err(err).Op(op)
		}
		{{ else -}}
		return gko.ErrAlEscribir.Err(err).Op(op)
		{{- end }}
	}
	return nil
}