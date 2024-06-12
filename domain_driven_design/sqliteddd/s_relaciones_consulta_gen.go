package sqliteddd

import (
	"database/sql"
	"net/http"
	"strings"

	"monorepo/domain_driven_design/ddd"
	"monorepo/gecko"
)

//  ================================================================  //
//  ========== MYSQL/CONSTANTES ====================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
const columnasConsultaRelacion string = "consulta_id, posicion, tipo_join, join_tabla_id, join_as, join_on, from_tabla_id"

// Origen de los datos de ddd.ConsultaRelacion
//
// FROM consulta_relaciones
const fromConsultaRelacion string = "FROM consulta_relaciones "

//  ================================================================  //
//  ========== MYSQL/TBL-INSERT ====================================  //

// InsertConsultaRelacion valida el registro y lo inserta en la base de datos.
func (s *Repositorio) InsertConsultaRelacion(rel ddd.ConsultaRelacion) error {
	const op string = "mysqlddd.InsertConsultaRelacion"
	if rel.ConsultaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("ConsultaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if rel.Posicion == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("Posicion sin especificar").Ctx(op, "pk_indefinida")
	}
	err := rel.Validar()
	if err != nil {
		return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec("INSERT INTO consulta_relaciones "+
		"(consulta_id, posicion, tipo_join, join_tabla_id, join_as, join_on, from_tabla_id) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?) ",
		rel.ConsultaID, rel.Posicion, rel.TipoJoin.String, rel.JoinTablaID, rel.JoinAs, rel.JoinOn, rel.FromTablaID,
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
//  ========== MYSQL/TBL-UPDATE ====================================  //

// UpdateConsultaRelacion valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateConsultaRelacion(rel ddd.ConsultaRelacion) error {
	const op string = "mysqlddd.UpdateConsultaRelacion"
	if rel.ConsultaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("ConsultaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if rel.Posicion == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("Posicion sin especificar").Ctx(op, "pk_indefinida")
	}
	err := rel.Validar()
	if err != nil {
		return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec(
		"UPDATE consulta_relaciones SET "+
			"consulta_id=?, posicion=?, tipo_join=?, join_tabla_id=?, join_as=?, join_on=?, from_tabla_id=? "+
			"WHERE consulta_id = ? AND posicion = ?",
		rel.ConsultaID, rel.Posicion, rel.TipoJoin.String, rel.JoinTablaID, rel.JoinAs, rel.JoinOn, rel.FromTablaID,
		rel.ConsultaID, rel.Posicion,
	)
	if err != nil {
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROWS =====================================  //

// scanRowsConsultaRelacion escanea cada row en la struct ConsultaRelacion
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsConsultaRelacion(rows *sql.Rows, op string) ([]ddd.ConsultaRelacion, error) {
	defer rows.Close()
	items := []ddd.ConsultaRelacion{}
	for rows.Next() {
		rel := ddd.ConsultaRelacion{}
		var tipoJoin string
		err := rows.Scan(
			&rel.ConsultaID, &rel.Posicion, &tipoJoin, &rel.JoinTablaID, &rel.JoinAs, &rel.JoinOn, &rel.FromTablaID,
		)
		if err != nil {
			return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}
		rel.TipoJoin = ddd.SetTipoJoinDB(tipoJoin)
		items = append(items, rel)
	}
	return items, nil
}

//  ================================================================  //
//  ========== MYSQL/LIST_BY =======================================  //

// ListConsultaRelacionesByConsultaID retorna los registros a partir de ConsultaID.
func (s *Repositorio) ListConsultaRelacionesByConsultaID(ConsultaID int) ([]ddd.ConsultaRelacion, error) {
	const op string = "mysqlddd.ListConsultaRelacionesByConsultaID"
	if ConsultaID == 0 {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("ConsultaID sin especificar").Ctx(op, "param_indefinido")
	}
	where := "WHERE consulta_id = ?"
	argumentos := []any{}
	argumentos = append(argumentos, ConsultaID)

	rows, err := s.db.Query(
		"SELECT "+columnasConsultaRelacion+" "+fromConsultaRelacion+
			where,
		argumentos...,
	)
	if err != nil {
		return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return s.scanRowsConsultaRelacion(rows, op)
}
