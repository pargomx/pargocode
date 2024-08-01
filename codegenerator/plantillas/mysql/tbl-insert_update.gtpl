// InsertUpdate{{ .Tabla.NombreItem }} valida e inserta o sobreescribe el registro en la base de datos.
func (s *Repositorio) InsertUpdate{{ .Tabla.NombreItem }}({{ .Tabla.NombreAbrev }} {{ .Tabla.Paquete.Nombre }}.{{ .Tabla.NombreItem }}) error {
	const op string = "InsertUpdate{{ .Tabla.NombreItem }}"
	{{ range .Tabla.CamposRequeridosOrPK -}}
		{{ if .PrimaryKey }}{{ .IfZeroReturnErr "pk_indefinida" $.Tabla.NombreAbrev -}}
		{{ else }}{{ .IfZeroReturnErr "required_sin_valor" $.Tabla.NombreAbrev }}{{ end -}}
	{{ end -}}
	_, err = s.db.Exec("INSERT INTO {{ .Tabla.NombreRepo }} "+
		"({{ .Tabla.CamposEditablesAsSnakeList ", " }}) "+
		"VALUES ({{ .Tabla.CamposEditablesAsPlaceholders }}) ",
		{{ .Tabla.CamposEditablesAsArguments .Tabla.NombreAbrev }},
	)
	if err != nil {
		if strings.HasPrefix(err.Error(),"Error 1062 (23000)"){
			_, err = s.db.Exec(
				"UPDATE {{ .Tabla.NombreRepo }} SET " +
				"{{ .Tabla.CamposEditablesAsSnakeEqPlaceholder }} " +
				"{{ .Tabla.PrimaryKeysAsSqlWhere }}",
				{{ .Tabla.CamposEditablesAsArguments .Tabla.NombreAbrev }},
				{{ .Tabla.PrimaryKeysAsArguments .Tabla.NombreAbrev }},
			)
			if err != nil {
				return gko.ErrInesperado().Err(err).Op(op)
			}
		}else if strings.HasPrefix(err.Error(), "Error 1452 (23000)") {
			return gko.ErrDatoInvalido().Err(err).Op(op).Msg("No se puede insertar la información porque el registro asociado no existe")
		}else {
			return gko.ErrInesperado().Err(err).Op(op)
		}
	}
	return nil
}