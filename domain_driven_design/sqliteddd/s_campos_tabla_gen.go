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
const columnasCampo string = "campo_id, tabla_id, nombre_campo, nombre_columna, nombre_humano, tipo_go, tipo_sql, setter, importado, primary_key, foreign_key, uq, req, ro, filtro, nullable, max_lenght, uns, default_sql, especial, referencia_campo, expresion, es_femenino, descripcion, posicion"

// Origen de los datos de ddd.Campo
//
// FROM campos
const fromCampo string = "FROM campos "

//  ================================================================  //
//  ========== MYSQL/TBL-INSERT ====================================  //

// InsertCampo valida el registro y lo inserta en la base de datos.
func (s *Repositorio) InsertCampo(cam ddd.Campo) error {
	const op string = "mysqlddd.InsertCampo"
	if cam.CampoID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("CampoID sin especificar").Ctx(op, "pk_indefinida")
	}
	if cam.NombreCampo == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("NombreCampo sin especificar").Ctx(op, "required_sin_valor")
	}
	if cam.NombreColumna == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("NombreColumna sin especificar").Ctx(op, "required_sin_valor")
	}
	if cam.NombreHumano == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("NombreHumano sin especificar").Ctx(op, "required_sin_valor")
	}
	err := cam.Validar()
	if err != nil {
		return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec("INSERT INTO campos "+
		"(campo_id, tabla_id, nombre_campo, nombre_columna, nombre_humano, tipo_go, tipo_sql, setter, importado, primary_key, foreign_key, uq, req, ro, filtro, nullable, max_lenght, uns, default_sql, especial, referencia_campo, expresion, es_femenino, descripcion, posicion) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, (SELECT count(campo_id) + 1 FROM campos WHERE tabla_id = ?)) ",
		cam.CampoID, cam.TablaID, cam.NombreCampo, cam.NombreColumna, cam.NombreHumano, cam.TipoGo, cam.TipoSql, cam.Setter, cam.Importado, cam.PrimaryKey, cam.ForeignKey, cam.Uq, cam.Req, cam.Ro, cam.Filtro, cam.Nullable, cam.MaxLenght, cam.Uns, cam.DefaultSql, cam.Especial, cam.ReferenciaCampo, cam.Expresion, cam.EsFemenino, cam.Descripcion, cam.TablaID,
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

// UpdateCampo valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateCampo(cam ddd.Campo) error {
	const op string = "mysqlddd.UpdateCampo"
	if cam.CampoID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("CampoID sin especificar").Ctx(op, "pk_indefinida")
	}
	if cam.NombreCampo == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("NombreCampo sin especificar").Ctx(op, "required_sin_valor")
	}
	if cam.NombreColumna == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("NombreColumna sin especificar").Ctx(op, "required_sin_valor")
	}
	if cam.NombreHumano == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("NombreHumano sin especificar").Ctx(op, "required_sin_valor")
	}
	err := cam.Validar()
	if err != nil {
		return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec(
		"UPDATE campos SET "+
			"campo_id=?, tabla_id=?, nombre_campo=?, nombre_columna=?, nombre_humano=?, tipo_go=?, tipo_sql=?, setter=?, importado=?, primary_key=?, foreign_key=?, uq=?, req=?, ro=?, filtro=?, nullable=?, max_lenght=?, uns=?, default_sql=?, especial=?, referencia_campo=?, expresion=?, es_femenino=?, descripcion=?, posicion=? "+
			"WHERE campo_id = ?",
		cam.CampoID, cam.TablaID, cam.NombreCampo, cam.NombreColumna, cam.NombreHumano, cam.TipoGo, cam.TipoSql, cam.Setter, cam.Importado, cam.PrimaryKey, cam.ForeignKey, cam.Uq, cam.Req, cam.Ro, cam.Filtro, cam.Nullable, cam.MaxLenght, cam.Uns, cam.DefaultSql, cam.Especial, cam.ReferenciaCampo, cam.Expresion, cam.EsFemenino, cam.Descripcion, cam.Posicion,
		cam.CampoID,
	)
	if err != nil {
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/TBL-DELETE ====================================  //

// DeleteCampo elimina permanentemente un registro del campo de la base de datos.
// Error si el registro no existe o si no se da la clave primaria.
func (s *Repositorio) DeleteCampo(CampoID int) error {
	const op string = "mysqlddd.DeleteCampo"
	if CampoID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("CampoID sin especificar").Ctx(op, "pk_indefinida")
	}
	// Verificar que solo se borre un registro.
	var num int
	err := s.db.QueryRow("SELECT COUNT(campo_id) FROM campos WHERE campo_id = ?",
		CampoID,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gecko.NewErr(http.StatusNotFound).Err(ddd.ErrCampoNotFound).Op(op)
		}
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	if num > 1 {
		return gecko.NewErr(http.StatusInternalServerError).Err(nil).Op(op).Msgf("abortado porque serían borrados %v registros", num)
	} else if num == 0 {
		return gecko.NewErr(http.StatusNotFound).Err(ddd.ErrCampoNotFound).Op(op).Msg("cero resultados")
	}
	// Eliminar registro
	_, err = s.db.Exec(
		"DELETE FROM campos WHERE campo_id = ?",
		CampoID,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1451 (23000)") {
			return gecko.NewErr(http.StatusConflict).Err(err).Op(op).Msg("Este registro es referenciado por otros y no se puede eliminar")
		} else {
			return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROW ======================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowCampo(row *sql.Row, cam *ddd.Campo, op string) error {
	var referenciaCampo sql.NullInt64
	err := row.Scan(
		&cam.CampoID, &cam.TablaID, &cam.NombreCampo, &cam.NombreColumna, &cam.NombreHumano, &cam.TipoGo, &cam.TipoSql, &cam.Setter, &cam.Importado, &cam.PrimaryKey, &cam.ForeignKey, &cam.Uq, &cam.Req, &cam.Ro, &cam.Filtro, &cam.Nullable, &cam.MaxLenght, &cam.Uns, &cam.DefaultSql, &cam.Especial, &referenciaCampo, &cam.Expresion, &cam.EsFemenino, &cam.Descripcion, &cam.Posicion,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gecko.NewErr(http.StatusNotFound).Msg("el campo no se encuentra").Op(op)
		}
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}

	if referenciaCampo.Valid {
		num := int(referenciaCampo.Int64)
		cam.ReferenciaCampo = &num
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/GET ===========================================  //

// GetCampo devuelve un Campo de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetCampo(CampoID int) (*ddd.Campo, error) {
	const op string = "mysqlddd.GetCampo"
	if CampoID == 0 {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("CampoID sin especificar").Ctx(op, "pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasCampo+" "+fromCampo+
			"WHERE campo_id = ?",
		CampoID,
	)
	cam := &ddd.Campo{}
	return cam, s.scanRowCampo(row, cam, op)
}

//  ================================================================  //
//  ========== MYSQL/GET_BY ========================================  //

// GetCampoByTablaIDNombreCampo devuelve un Campo de la DB.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetCampoByTablaIDNombreCampo(TablaID int, NombreCampo string) (*ddd.Campo, error) {
	const op string = "mysqlddd.GetCampoByTablaIDNombreCampo"
	if TablaID == 0 {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("TablaID sin especificar").Ctx(op, "param_indefinido")
	}
	if NombreCampo == "" {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("NombreCampo sin especificar").Ctx(op, "param_indefinido")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasCampo+" "+fromCampo+
			"WHERE tabla_id = ? AND nombre_campo = ?",
		TablaID, NombreCampo,
	)
	cam := &ddd.Campo{}
	return cam, s.scanRowCampo(row, cam, op)
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROWS =====================================  //

// scanRowsCampo escanea cada row en la struct Campo
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsCampo(rows *sql.Rows, op string) ([]ddd.Campo, error) {
	defer rows.Close()
	items := []ddd.Campo{}
	for rows.Next() {
		cam := ddd.Campo{}
		var referenciaCampo sql.NullInt64
		err := rows.Scan(
			&cam.CampoID, &cam.TablaID, &cam.NombreCampo, &cam.NombreColumna, &cam.NombreHumano, &cam.TipoGo, &cam.TipoSql, &cam.Setter, &cam.Importado, &cam.PrimaryKey, &cam.ForeignKey, &cam.Uq, &cam.Req, &cam.Ro, &cam.Filtro, &cam.Nullable, &cam.MaxLenght, &cam.Uns, &cam.DefaultSql, &cam.Especial, &referenciaCampo, &cam.Expresion, &cam.EsFemenino, &cam.Descripcion, &cam.Posicion,
		)
		if err != nil {
			return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}

		if referenciaCampo.Valid {
			num := int(referenciaCampo.Int64)
			cam.ReferenciaCampo = &num
		}
		items = append(items, cam)
	}
	return items, nil
}
