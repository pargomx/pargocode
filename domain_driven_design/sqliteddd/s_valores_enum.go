package sqliteddd

import (
	"monorepo/domain_driven_design/ddd"
	"monorepo/gecko"
	"net/http"
)

// GetValoresEnum retorna los registros a partir de CampoID.
func (s *Repositorio) GetValoresEnum(CampoID int) ([]ddd.ValorEnum, error) {
	const op string = "mysqlddd.GetValoresEnum"
	if CampoID == 0 {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("CampoID sin especificar").Ctx(op, "CampoID_indefinido")
	}
	rows, err := s.db.Query(
		"SELECT "+columnasValorEnum+" FROM valores_enum WHERE campo_id = ? ORDER BY numero ASC",
		CampoID,
	)
	if err != nil {
		return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return s.scanRowsValorEnum(rows, op)
}

func (s *Repositorio) GuardarValoresEnum(campoID int, vals []ddd.ValorEnum) error {
	const op string = "sqliteddd.GuardarValoresEnum"
	if campoID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("CampoID sin especificar").Ctx(op, "pk_indefinida")
	}

	_, err := s.db.Exec("DELETE FROM valores_enum WHERE campo_id = ?", campoID)
	if err != nil {
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op).Op("borrar_valores_enum_anteriores")
	}

	sqlStr := "INSERT INTO valores_enum (campo_id, numero, clave, etiqueta, descripcion) VALUES "
	args := []any{}
	for _, val := range vals {
		if val.Clave == "" {
			return gecko.NewErr(http.StatusBadRequest).Msg("Clave sin especificar").Ctx(op, "pk_indefinida")
		}
		if val.Etiqueta == "" {
			return gecko.NewErr(http.StatusBadRequest).Msg("Etiqueta sin especificar").Ctx(op, "required_sin_valor")
		}
		err := val.Validar()
		if err != nil {
			return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
		}
		sqlStr += "(?, ?, ?, ?, ?), "
		args = append(args, campoID, val.Numero, val.Clave, val.Etiqueta, val.Descripcion)
	}

	sqlStr = sqlStr[0 : len(sqlStr)-2] // Quitar última coma
	_, err = s.db.Exec(sqlStr, args...)
	if err != nil {
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return nil
}
