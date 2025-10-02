package sqliteddd

import (
	"database/sql"
	"errors"

	"github.com/pargomx/gecko/gko"

	"github.com/pargomx/pargocode/ddd"
)

//  ================================================================  //
//  ========== INSERT ==============================================  //

func (s *Repositorio) InsertCampo(cam ddd.Campo) error {
	const op string = "InsertCampo"
	if cam.CampoID == 0 {
		return gko.ErrDatoIndef.Msg("CampoID sin especificar").Str("pk_indefinida").Op(op)
	}
	if cam.NombreCampo == "" {
		return gko.ErrDatoIndef.Msg("NombreCampo sin especificar").Str("required_sin_valor").Op(op)
	}
	if cam.NombreColumna == "" {
		return gko.ErrDatoIndef.Msg("NombreColumna sin especificar").Str("required_sin_valor").Op(op)
	}
	if cam.NombreHumano == "" {
		return gko.ErrDatoIndef.Msg("NombreHumano sin especificar").Str("required_sin_valor").Op(op)
	}
	_, err := s.db.Exec("INSERT INTO campos "+
		"(campo_id, tabla_id, nombre_campo, nombre_columna, nombre_humano, tipo_go, tipo_sql, setter, importado, primary_key, foreign_key, uq, req, ro, filtro, nullable, max_lenght, uns, default_sql, especial, referencia_campo, expresion, es_femenino, descripcion, posicion, zero_is_null) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ",
		cam.CampoID, cam.TablaID, cam.NombreCampo, cam.NombreColumna, cam.NombreHumano, cam.TipoGo, cam.TipoSql, cam.Setter, cam.Importado, cam.PrimaryKey, cam.ForeignKey, cam.Uq, cam.Req, cam.Ro, cam.Filtro, cam.Nullable, cam.MaxLenght, cam.Uns, cam.DefaultSql, cam.Especial, cam.ReferenciaCampo, cam.Expresion, cam.EsFemenino, cam.Descripcion, cam.Posicion, cam.ZeroIsNull,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== UPDATE ==============================================  //

// UpdateCampo valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateCampo(cam ddd.Campo) error {
	const op string = "UpdateCampo"
	if cam.CampoID == 0 {
		return gko.ErrDatoIndef.Msg("CampoID sin especificar").Str("pk_indefinida").Op(op)
	}
	if cam.NombreCampo == "" {
		return gko.ErrDatoIndef.Msg("NombreCampo sin especificar").Str("required_sin_valor").Op(op)
	}
	if cam.NombreColumna == "" {
		return gko.ErrDatoIndef.Msg("NombreColumna sin especificar").Str("required_sin_valor").Op(op)
	}
	if cam.NombreHumano == "" {
		return gko.ErrDatoIndef.Msg("NombreHumano sin especificar").Str("required_sin_valor").Op(op)
	}
	_, err := s.db.Exec(
		"UPDATE campos SET "+
			"campo_id=?, tabla_id=?, nombre_campo=?, nombre_columna=?, nombre_humano=?, tipo_go=?, tipo_sql=?, setter=?, importado=?, primary_key=?, foreign_key=?, uq=?, req=?, ro=?, filtro=?, nullable=?, max_lenght=?, uns=?, default_sql=?, especial=?, referencia_campo=?, expresion=?, es_femenino=?, descripcion=?, posicion=?, zero_is_null=? "+
			"WHERE campo_id = ?",
		cam.CampoID, cam.TablaID, cam.NombreCampo, cam.NombreColumna, cam.NombreHumano, cam.TipoGo, cam.TipoSql, cam.Setter, cam.Importado, cam.PrimaryKey, cam.ForeignKey, cam.Uq, cam.Req, cam.Ro, cam.Filtro, cam.Nullable, cam.MaxLenght, cam.Uns, cam.DefaultSql, cam.Especial, cam.ReferenciaCampo, cam.Expresion, cam.EsFemenino, cam.Descripcion, cam.Posicion, cam.ZeroIsNull,
		cam.CampoID,
	)
	if err != nil {
		return gko.ErrInesperado.Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== EXISTE ==============================================  //

// Retorna error nil si existe solo un registro con esta clave primaria.
func (s *Repositorio) ExisteCampo(CampoID int) error {
	const op string = "ExisteCampo"
	var num int
	err := s.db.QueryRow("SELECT COUNT(campo_id) FROM campos WHERE campo_id = ?",
		CampoID,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gko.ErrNoEncontrado.Msg("Campo no encontrado").Op(op)
		}
		return gko.ErrInesperado.Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado.Err(nil).Op(op).Str("existen más de un registro para la pk").Ctx("registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado.Msg("Campo no encontrado").Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== DELETE ==============================================  //

func (s *Repositorio) DeleteCampo(CampoID int) error {
	const op string = "DeleteCampo"
	if CampoID == 0 {
		return gko.ErrDatoIndef.Msg("CampoID sin especificar").Str("pk_indefinida").Op(op)
	}
	err := s.ExisteCampo(CampoID)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	_, err = s.db.Exec(
		"DELETE FROM campos WHERE campo_id = ?",
		CampoID,
	)
	if err != nil {
		return gko.ErrAlEscribir.Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== CONSTANTES ==========================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
//
//	campo_id,
//	tabla_id,
//	nombre_campo,
//	nombre_columna,
//	nombre_humano,
//	tipo_go,
//	tipo_sql,
//	setter,
//	importado,
//	primary_key,
//	foreign_key,
//	uq,
//	req,
//	ro,
//	filtro,
//	nullable,
//	max_lenght,
//	uns,
//	default_sql,
//	especial,
//	referencia_campo,
//	expresion,
//	es_femenino,
//	descripcion,
//	posicion,
//	zero_is_null
const columnasCampo string = "campo_id, tabla_id, nombre_campo, nombre_columna, nombre_humano, tipo_go, tipo_sql, setter, importado, primary_key, foreign_key, uq, req, ro, filtro, nullable, max_lenght, uns, default_sql, especial, referencia_campo, expresion, es_femenino, descripcion, posicion, zero_is_null"

// Origen de los datos de ddd.Campo
//
//	FROM campos
const fromCampo string = "FROM campos "

//  ================================================================  //
//  ========== SCAN ================================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowCampo(row *sql.Row, cam *ddd.Campo) error {
	var referenciaCampo sql.NullInt64
	err := row.Scan(
		&cam.CampoID, &cam.TablaID, &cam.NombreCampo, &cam.NombreColumna, &cam.NombreHumano, &cam.TipoGo, &cam.TipoSql, &cam.Setter, &cam.Importado, &cam.PrimaryKey, &cam.ForeignKey, &cam.Uq, &cam.Req, &cam.Ro, &cam.Filtro, &cam.Nullable, &cam.MaxLenght, &cam.Uns, &cam.DefaultSql, &cam.Especial, &referenciaCampo, &cam.Expresion, &cam.EsFemenino, &cam.Descripcion, &cam.Posicion, &cam.ZeroIsNull,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado.Msg("Campo no encontrado")
		}
		return gko.ErrInesperado.Err(err)
	}

	if referenciaCampo.Valid {
		num := int(referenciaCampo.Int64)
		cam.ReferenciaCampo = &num
	}
	return nil
}

//  ================================================================  //
//  ========== GET =================================================  //

// GetCampo devuelve un Campo de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetCampo(CampoID int) (*ddd.Campo, error) {
	const op string = "GetCampo"
	if CampoID == 0 {
		return nil, gko.ErrDatoIndef.Msg("CampoID sin especificar").Str("pk_indefinida").Op(op)
	}
	row := s.db.QueryRow(
		"SELECT "+columnasCampo+" "+fromCampo+
			"WHERE campo_id = ?",
		CampoID,
	)
	cam := &ddd.Campo{}
	err := s.scanRowCampo(row, cam)
	if err != nil {
		return nil, err
	}
	return cam, nil
}

//  ================================================================  //
//  ========== GET_BY TABLA_ID NOMBRE_CAMPO ========================  //

// GetCampoByTablaIDNombreCampo devuelve un Campo de la DB.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetCampoByTablaIDNombreCampo(TablaID int, NombreCampo string) (*ddd.Campo, error) {
	const op string = "GetCampoByTablaIDNombreCampo"
	if TablaID == 0 {
		return nil, gko.ErrDatoIndef.Msg("TablaID sin especificar").Str("param_indefinido").Op(op)
	}
	if NombreCampo == "" {
		return nil, gko.ErrDatoIndef.Msg("NombreCampo sin especificar").Str("param_indefinido").Op(op)
	}
	row := s.db.QueryRow(
		"SELECT "+columnasCampo+" "+fromCampo+
			"WHERE tabla_id = ? AND nombre_campo = ?",
		TablaID, NombreCampo,
	)
	cam := &ddd.Campo{}
	err := s.scanRowCampo(row, cam)
	if err != nil {
		return nil, gko.Err(err).Op(op)
	}
	return cam, nil
}

//  ================================================================  //
//  ========== SCAN ================================================  //

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
			&cam.CampoID, &cam.TablaID, &cam.NombreCampo, &cam.NombreColumna, &cam.NombreHumano, &cam.TipoGo, &cam.TipoSql, &cam.Setter, &cam.Importado, &cam.PrimaryKey, &cam.ForeignKey, &cam.Uq, &cam.Req, &cam.Ro, &cam.Filtro, &cam.Nullable, &cam.MaxLenght, &cam.Uns, &cam.DefaultSql, &cam.Especial, &referenciaCampo, &cam.Expresion, &cam.EsFemenino, &cam.Descripcion, &cam.Posicion, &cam.ZeroIsNull,
		)
		if err != nil {
			return nil, gko.ErrInesperado.Err(err).Op(op)
		}

		if referenciaCampo.Valid {
			num := int(referenciaCampo.Int64)
			cam.ReferenciaCampo = &num
		}
		items = append(items, cam)
	}
	return items, nil
}
