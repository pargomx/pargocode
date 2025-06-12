// Retorna error nil si existe solo un registro con esta clave primaria.
func (s *Repositorio) Existe{{ .Tabla.NombreItem }}({{ .Tabla.PrimaryKeysAsFuncParams }}) error {
	const op string = "Existe{{ .Tabla.NombreItem }}"
	var num int
	err := s.db.QueryRow("SELECT COUNT({{ .Tabla.PrimerCampo.NombreColumna }}) FROM {{ .Tabla.NombreRepo }} {{ .Tabla.PrimaryKeysAsSqlWhere }}",
		{{ .Tabla.PrimaryKeysAsArguments "" }},
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gko.ErrNoEncontrado.Msg("{{ .Tabla.Tabla.Humano }} no encontrado").Op(op)
		}
		return gko.ErrInesperado.Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado.Err(nil).Op(op).Str("existen más de un registro para la pk").Ctx("registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado.Msg("{{ .Tabla.Tabla.Humano }} no encontrado").Op(op)
	}
	return nil
}