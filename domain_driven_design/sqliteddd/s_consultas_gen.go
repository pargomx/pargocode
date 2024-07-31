package sqliteddd

import (
	"database/sql"
	"errors"
	"strings"

	"monorepo/domain_driven_design/ddd"

	"github.com/pargomx/gecko/gko"
)

//  ================================================================  //
//  ========== MYSQL/CONSTANTES ====================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
const columnasConsulta string = "consulta_id, paquete_id, tabla_id, nombre_item, nombre_items, abrev, es_femenino, descripcion, directrices"

// Origen de los datos de ddd.Consulta
//
// FROM consultas
const fromConsulta string = "FROM consultas "

//  ================================================================  //
//  ========== MYSQL/TBL-INSERT ====================================  //

// InsertConsulta valida el registro y lo inserta en la base de datos.
func (s *Repositorio) InsertConsulta(con ddd.Consulta) error {
	const op string = "mysqlddd.InsertConsulta"
	if con.ConsultaID == 0 {
		return gko.ErrDatoInvalido().Msg("ConsultaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if con.NombreItem == "" {
		return gko.ErrDatoInvalido().Msg("NombreItem sin especificar").Ctx(op, "required_sin_valor")
	}
	err := con.Validar()
	if err != nil {
		return gko.ErrDatoInvalido().Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec("INSERT INTO consultas "+
		"(consulta_id, paquete_id, tabla_id, nombre_item, nombre_items, abrev, es_femenino, descripcion, directrices) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ",
		con.ConsultaID, con.PaqueteID, con.TablaID, con.NombreItem, con.NombreItems, con.Abrev, con.EsFemenino, con.Descripcion, con.Directrices,
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
//  ========== MYSQL/TBL-UPDATE ====================================  //

// UpdateConsulta valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateConsulta(con ddd.Consulta) error {
	const op string = "mysqlddd.UpdateConsulta"
	if con.ConsultaID == 0 {
		return gko.ErrDatoInvalido().Msg("ConsultaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if con.NombreItem == "" {
		return gko.ErrDatoInvalido().Msg("NombreItem sin especificar").Ctx(op, "required_sin_valor")
	}
	err := con.Validar()
	if err != nil {
		return gko.ErrDatoInvalido().Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec(
		"UPDATE consultas SET "+
			"consulta_id=?, paquete_id=?, tabla_id=?, nombre_item=?, nombre_items=?, abrev=?, es_femenino=?, descripcion=?, directrices=? "+
			"WHERE consulta_id = ?",
		con.ConsultaID, con.PaqueteID, con.TablaID, con.NombreItem, con.NombreItems, con.Abrev, con.EsFemenino, con.Descripcion, con.Directrices,
		con.ConsultaID,
	)
	if err != nil {
		return gko.ErrInesperado().Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/TBL-DELETE ====================================  //

// DeleteConsulta elimina permanentemente un registro de la consulta de la base de datos.
// Error si el registro no existe o si no se da la clave primaria.
func (s *Repositorio) DeleteConsulta(ConsultaID int) error {
	const op string = "mysqlddd.DeleteConsulta"
	if ConsultaID == 0 {
		return gko.ErrDatoInvalido().Msg("ConsultaID sin especificar").Ctx(op, "pk_indefinida")
	}
	// Verificar que solo se borre un registro.
	var num int
	err := s.db.QueryRow("SELECT COUNT(consulta_id) FROM consultas WHERE consulta_id = ?",
		ConsultaID,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gko.ErrNoEncontrado().Err(ddd.ErrConsultaNotFound).Op(op)
		}
		return gko.ErrInesperado().Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado().Err(nil).Op(op).Msgf("abortado porque serían borrados %v registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado().Err(ddd.ErrConsultaNotFound).Op(op).Msg("cero resultados")
	}
	// Eliminar registro
	_, err = s.db.Exec(
		"DELETE FROM consulta_campos WHERE consulta_id = ?",
		ConsultaID,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1451 (23000)") {
			return gko.ErrYaExiste().Err(err).Op(op).Msg("Este registro es referenciado por otros y no se puede eliminar")
		} else {
			return gko.ErrInesperado().Err(err).Op(op)
		}
	}
	_, err = s.db.Exec(
		"DELETE FROM consulta_relaciones WHERE consulta_id = ?",
		ConsultaID,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1451 (23000)") {
			return gko.ErrYaExiste().Err(err).Op(op).Msg("Este registro es referenciado por otros y no se puede eliminar")
		} else {
			return gko.ErrInesperado().Err(err).Op(op)
		}
	}
	_, err = s.db.Exec(
		"DELETE FROM consultas WHERE consulta_id = ?",
		ConsultaID,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1451 (23000)") {
			return gko.ErrYaExiste().Err(err).Op(op).Msg("Este registro es referenciado por otros y no se puede eliminar")
		} else {
			return gko.ErrInesperado().Err(err).Op(op)
		}
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROW ======================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowConsulta(row *sql.Row, con *ddd.Consulta, op string) error {

	err := row.Scan(
		&con.ConsultaID, &con.PaqueteID, &con.TablaID, &con.NombreItem, &con.NombreItems, &con.Abrev, &con.EsFemenino, &con.Descripcion, &con.Directrices,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado().Msg("la consulta no se encuentra").Op(op)
		}
		return gko.ErrInesperado().Err(err).Op(op)
	}

	return nil
}

//  ================================================================  //
//  ========== MYSQL/GET ===========================================  //

// GetConsulta devuelve un Consulta de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetConsulta(ConsultaID int) (*ddd.Consulta, error) {
	const op string = "mysqlddd.GetConsulta"
	if ConsultaID == 0 {
		return nil, gko.ErrDatoInvalido().Msg("ConsultaID sin especificar").Ctx(op, "pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasConsulta+" "+fromConsulta+
			"WHERE consulta_id = ?",
		ConsultaID,
	)
	con := &ddd.Consulta{}
	return con, s.scanRowConsulta(row, con, op)
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROWS =====================================  //

// scanRowsConsulta escanea cada row en la struct Consulta
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsConsulta(rows *sql.Rows, op string) ([]ddd.Consulta, error) {
	defer rows.Close()
	items := []ddd.Consulta{}
	for rows.Next() {
		con := ddd.Consulta{}

		err := rows.Scan(
			&con.ConsultaID, &con.PaqueteID, &con.TablaID, &con.NombreItem, &con.NombreItems, &con.Abrev, &con.EsFemenino, &con.Descripcion, &con.Directrices,
		)
		if err != nil {
			return nil, gko.ErrInesperado().Err(err).Op(op)
		}

		items = append(items, con)
	}
	return items, nil
}

//  ================================================================  //
//  ========== MYSQL/LIST_BY =======================================  //

// ListConsultasByPaqueteID retorna los registros a partir de PaqueteID.
func (s *Repositorio) ListConsultasByPaqueteID(PaqueteID int) ([]ddd.Consulta, error) {
	const op string = "mysqlddd.ListConsultasByPaqueteID"
	if PaqueteID == 0 {
		return nil, gko.ErrDatoInvalido().Msg("PaqueteID sin especificar").Ctx(op, "param_indefinido")
	}
	where := "WHERE paquete_id = ?"
	argumentos := []any{}
	argumentos = append(argumentos, PaqueteID)

	rows, err := s.db.Query(
		"SELECT "+columnasConsulta+" "+fromConsulta+
			where,
		argumentos...,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsConsulta(rows, op)
}
