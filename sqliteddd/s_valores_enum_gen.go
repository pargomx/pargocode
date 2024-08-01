package sqliteddd

import (
	"database/sql"
	"strings"

	"monorepo/ddd"

	"github.com/pargomx/gecko/gko"
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
		return gko.ErrDatoInvalido().Msg("CampoID sin especificar").Ctx(op, "pk_indefinida")
	}
	if val.Clave == "" {
		return gko.ErrDatoInvalido().Msg("Clave sin especificar").Ctx(op, "pk_indefinida")
	}
	if val.Etiqueta == "" {
		return gko.ErrDatoInvalido().Msg("Etiqueta sin especificar").Ctx(op, "required_sin_valor")
	}
	err := val.Validar()
	if err != nil {
		return gko.ErrDatoInvalido().Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec("INSERT INTO valores_enum "+
		"(campo_id, numero, clave, etiqueta, descripcion) "+
		"VALUES (?, ?, ?, ?, ?) ",
		val.CampoID, val.Numero, val.Clave, val.Etiqueta, val.Descripcion,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062 (23000)") {
			return gko.ErrYaExiste().Err(err).Op(op)
		} else if strings.HasPrefix(err.Error(), "Error 1452 (23000)") {
			return gko.ErrDatoInvalido().Err(err).Op(op).Msg("No se puede insertar la información porque el registro asociado no existe")
		} else {
			return gko.ErrInesperado().Err(err).Op(op)
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
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}

		items = append(items, val)
	}
	return items, nil
}
