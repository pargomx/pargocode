func (s *Repositorio) Delete{{ .Tabla.NombreItem }}({{ .Tabla.PrimaryKeysAsFuncParams }}) error {
	const op string = "mysql{{ .Tabla.Paquete.Nombre }}.Delete{{ .Tabla.NombreItem }}"
	{{ range .Tabla.PrimaryKeys }}{{ .IfZeroReturnErr "pk_indefinida" "" }}{{ end -}}
	// Verificar que solo se borre un registro.
	var num int
	err := s.db.QueryRow("SELECT COUNT({{ .Tabla.PrimerCampo.NombreColumna }}) FROM {{ .Tabla.NombreRepo }} {{ .Tabla.PrimaryKeysAsSqlWhere }}",
		{{ .Tabla.PrimaryKeysAsArguments "" }},
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gecko.NewErr(http.StatusNotFound).Err({{ .Tabla.Paquete.Nombre }}.Err{{ .Tabla.NombreItem }}NotFound).Op(op)
		}
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	if num > 1 {
		return gecko.NewErr(http.StatusInternalServerError).Err(nil).Op(op).Msgf("abortado porque serían borrados %v registros", num)
	} else if num == 0 {
		return gecko.NewErr(http.StatusNotFound).Err({{ .Tabla.Paquete.Nombre }}.Err{{ .Tabla.NombreItem }}NotFound).Op(op).Msg("cero resultados")
	}
	// Eliminar registro
	_, err = s.db.Exec(
		"DELETE FROM {{ .Tabla.NombreRepo }} {{ .Tabla.PrimaryKeysAsSqlWhere }}",
		{{ .Tabla.PrimaryKeysAsArguments "" }},
	)
	if err != nil {
		if strings.HasPrefix(err.Error(),"Error 1451 (23000)"){
			return gecko.NewErr(http.StatusConflict).Err(err).Op(op).Msg("Este registro es referenciado por otros y no se puede eliminar")
		} else {
			return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}
	}
	return nil
}