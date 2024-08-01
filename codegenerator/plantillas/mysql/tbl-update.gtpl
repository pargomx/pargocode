// Update{{ .Tabla.NombreItem }} valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) Update{{ .Tabla.NombreItem }}({{ .Tabla.NombreAbrev }} {{ .Tabla.Paquete.Nombre }}.{{ .Tabla.NombreItem }}) error {
	const op string = "mysql{{ .Tabla.Paquete.Nombre }}.Update{{ .Tabla.NombreItem }}"
	{{ range .Tabla.CamposRequeridosOrPK -}}
		{{ if .PrimaryKey }}{{ .IfZeroReturnErr "pk_indefinida" $.Tabla.NombreAbrev -}}
		{{ else }}{{ .IfZeroReturnErr "required_sin_valor" $.Tabla.NombreAbrev }}{{ end -}}
	{{ end -}}
	err := {{ .Tabla.NombreAbrev }}.Validar()
	if err != nil {
		return gko.ErrDatoInvalido().Err(err).Op(op).Msg(err.Error())
	}
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
	return nil
}