package sqliteddd

import (
	"database/sql"
	"net/http"
	"strings"

	"monorepo/domain_driven_design/ddd"

	"github.com/pargomx/gecko"
)

//  ================================================================  //
//  ========== MYSQL/CONSTANTES ====================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
const columnasValorEnum string = "campo_id, numero, clave, etiqueta, descripcion"

//  ================================================================  //
//  ========== MYSQL/TBL-INSERT ====================================  //

// InsertValorEnum valida el registro y lo inserta en la base de datos.
func (s *Repositorio) InsertValorEnum(val ddd.ValorEnum) error {
	const op string = "mysqlddd.InsertValorEnum"
	if val.CampoID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("CampoID sin especificar").Ctx(op, "pk_indefinida")
	}
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
	_, err = s.db.Exec("INSERT INTO valores_enum "+
		"(campo_id, numero, clave, etiqueta, descripcion) "+
		"VALUES (?, ?, ?, ?, ?) ",
		val.CampoID, val.Numero, val.Clave, val.Etiqueta, val.Descripcion,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062 (23000)") {
			return gecko.NewErr(http.StatusConflict).Err(err).Op(op)
		} else if strings.HasPrefix(err.Error(), "Error 1452 (23000)") {
			return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg("No se puede insertar la información porque el registro asociado no existe")
		} else {
			return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROWS =====================================  //

// scanRowsValorEnum escanea cada row en la struct ValorEnum
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsValorEnum(rows *sql.Rows, op string) ([]ddd.ValorEnum, error) {
	defer rows.Close()
	items := []ddd.ValorEnum{}
	for rows.Next() {
		val := ddd.ValorEnum{}

		err := rows.Scan(
			&val.CampoID, &val.Numero, &val.Clave, &val.Etiqueta, &val.Descripcion,
		)
		if err != nil {
			return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}

		items = append(items, val)
	}
	return items, nil
}
