package sqliteddd

import (
	"database/sql"
	"errors"
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
const columnasConsultaCampo string = "consulta_id, posicion, expresion, alias_sql, nombre_campo, tipo_go, campo_id, pk, filtro, group_by, descripcion"

// Origen de los datos de ddd.ConsultaCampo
//
// FROM consulta_campos
const fromConsultaCampo string = "FROM consulta_campos "

//  ================================================================  //
//  ========== MYSQL/TBL-INSERT ====================================  //

// InsertConsultaCampo valida el registro y lo inserta en la base de datos.
func (s *Repositorio) InsertConsultaCampo(campcons ddd.ConsultaCampo) error {
	const op string = "mysqlddd.InsertConsultaCampo"
	if campcons.ConsultaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("ConsultaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if campcons.Posicion == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("Posicion sin especificar").Ctx(op, "pk_indefinida")
	}
	if campcons.NombreCampo == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("NombreCampo sin especificar").Ctx(op, "required_sin_valor")
	}
	if campcons.TipoGo == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("TipoGo sin especificar").Ctx(op, "required_sin_valor")
	}
	err := campcons.Validar()
	if err != nil {
		return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec("INSERT INTO consulta_campos "+
		"(consulta_id, posicion, expresion, alias_sql, nombre_campo, tipo_go, campo_id, pk, filtro, group_by, descripcion) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ",
		campcons.ConsultaID, campcons.Posicion, campcons.Expresion, campcons.AliasSql, campcons.NombreCampo, campcons.TipoGo, campcons.CampoID, campcons.Pk, campcons.Filtro, campcons.GroupBy, campcons.Descripcion,
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

// UpdateConsultaCampo valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateConsultaCampo(campcons ddd.ConsultaCampo) error {
	const op string = "mysqlddd.UpdateConsultaCampo"
	if campcons.ConsultaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("ConsultaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if campcons.Posicion == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("Posicion sin especificar").Ctx(op, "pk_indefinida")
	}
	if campcons.NombreCampo == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("NombreCampo sin especificar").Ctx(op, "required_sin_valor")
	}
	if campcons.TipoGo == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("TipoGo sin especificar").Ctx(op, "required_sin_valor")
	}
	err := campcons.Validar()
	if err != nil {
		return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec(
		"UPDATE consulta_campos SET "+
			"consulta_id=?, posicion=?, expresion=?, alias_sql=?, nombre_campo=?, tipo_go=?, campo_id=?, pk=?, filtro=?, group_by=?, descripcion=? "+
			"WHERE consulta_id = ? AND posicion = ?",
		campcons.ConsultaID, campcons.Posicion, campcons.Expresion, campcons.AliasSql, campcons.NombreCampo, campcons.TipoGo, campcons.CampoID, campcons.Pk, campcons.Filtro, campcons.GroupBy, campcons.Descripcion,
		campcons.ConsultaID, campcons.Posicion,
	)
	if err != nil {
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROW ======================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowConsultaCampo(row *sql.Row, campcons *ddd.ConsultaCampo, op string) error {
	var campoID sql.NullInt64
	err := row.Scan(
		&campcons.ConsultaID, &campcons.Posicion, &campcons.Expresion, &campcons.AliasSql, &campcons.NombreCampo, &campcons.TipoGo, &campoID, &campcons.Pk, &campcons.Filtro, &campcons.GroupBy, &campcons.Descripcion,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gecko.NewErr(http.StatusNotFound).Msg("el campo de consulta no se encuentra").Op(op)
		}
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}

	if campoID.Valid {
		num := int(campoID.Int64)
		campcons.CampoID = &num
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/GET ===========================================  //

// GetConsultaCampo devuelve un ConsultaCampo de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetConsultaCampo(ConsultaID int, Posicion int) (*ddd.ConsultaCampo, error) {
	const op string = "mysqlddd.GetConsultaCampo"
	if ConsultaID == 0 {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("ConsultaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if Posicion == 0 {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("Posicion sin especificar").Ctx(op, "pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasConsultaCampo+" "+fromConsultaCampo+
			"WHERE consulta_id = ? AND posicion = ?",
		ConsultaID, Posicion,
	)
	campcons := &ddd.ConsultaCampo{}
	return campcons, s.scanRowConsultaCampo(row, campcons, op)
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROWS =====================================  //

// scanRowsConsultaCampo escanea cada row en la struct ConsultaCampo
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsConsultaCampo(rows *sql.Rows, op string) ([]ddd.ConsultaCampo, error) {
	defer rows.Close()
	items := []ddd.ConsultaCampo{}
	for rows.Next() {
		campcons := ddd.ConsultaCampo{}
		var campoID sql.NullInt64
		err := rows.Scan(
			&campcons.ConsultaID, &campcons.Posicion, &campcons.Expresion, &campcons.AliasSql, &campcons.NombreCampo, &campcons.TipoGo, &campoID, &campcons.Pk, &campcons.Filtro, &campcons.GroupBy, &campcons.Descripcion,
		)
		if err != nil {
			return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}

		if campoID.Valid {
			num := int(campoID.Int64)
			campcons.CampoID = &num
		}
		items = append(items, campcons)
	}
	return items, nil
}

//  ================================================================  //
//  ========== MYSQL/LIST_BY =======================================  //

// ListConsultaCamposByConsultaID retorna los registros a partir de ConsultaID.
func (s *Repositorio) ListConsultaCamposByConsultaID(ConsultaID int) ([]ddd.ConsultaCampo, error) {
	const op string = "mysqlddd.ListConsultaCamposByConsultaID"
	if ConsultaID == 0 {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("ConsultaID sin especificar").Ctx(op, "param_indefinido")
	}
	where := "WHERE consulta_id = ?"
	argumentos := []any{}
	argumentos = append(argumentos, ConsultaID)

	rows, err := s.db.Query(
		"SELECT "+columnasConsultaCampo+" "+fromConsultaCampo+
			where,
		argumentos...,
	)
	if err != nil {
		return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return s.scanRowsConsultaCampo(rows, op)
}
