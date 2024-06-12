// Insert{{ .Tabla.NombreItem }} valida el registro y lo inserta en la base de datos.
func (s *Repositorio) Insert{{ .Tabla.NombreItem }}({{ .Tabla.NombreAbrev }} {{ .Tabla.Paquete.Nombre }}.{{ .Tabla.NombreItem }}) error {
	const op string = "mysql{{ .Tabla.Paquete.Nombre }}.Insert{{ .Tabla.NombreItem }}"
	{{ range .Tabla.CamposRequeridosOrPK -}}
		{{ if .PrimaryKey }}{{ .IfZeroReturnErr "pk_indefinida" $.Tabla.NombreAbrev -}}
		{{ else }}{{ .IfZeroReturnErr "required_sin_valor" $.Tabla.NombreAbrev }}{{ end -}}
	{{ end -}}
	err := {{ .Tabla.NombreAbrev }}.Validar()
	if err != nil {
		return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec("INSERT INTO {{ .Tabla.NombreRepo }} "+
		"({{ .Tabla.CamposEditablesAsSnakeList ", " }}) "+
		"VALUES ({{ .Tabla.CamposEditablesAsPlaceholders }}) ",
		{{ .Tabla.CamposEditablesAsArguments .Tabla.NombreAbrev }},
	)
	if err != nil {
		if strings.HasPrefix(err.Error(),"Error 1062 (23000)"){
			return gecko.NewErr(http.StatusConflict).Err(err).Op(op)
		} else if strings.HasPrefix(err.Error(), "Error 1452 (23000)") {
			return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg("No se puede insertar la información porque el registro asociado no existe")
		} else {
			return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}
	}
	return nil
}