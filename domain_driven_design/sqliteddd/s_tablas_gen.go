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
const columnasTabla string = "tabla_id, paquete_id, nombre_repo, nombre_item, nombre_items, abrev, humano, humano_plural, kebab, es_femenino, descripcion, directrices"

// Origen de los datos de ddd.Tabla
//
// FROM tablas
const fromTabla string = "FROM tablas "

//  ================================================================  //
//  ========== MYSQL/TBL-INSERT ====================================  //

// InsertTabla valida el registro y lo inserta en la base de datos.
func (s *Repositorio) InsertTabla(tab ddd.Tabla) error {
	const op string = "mysqlddd.InsertTabla"
	if tab.TablaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("TablaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if tab.NombreRepo == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("NombreRepo sin especificar").Ctx(op, "required_sin_valor")
	}
	if tab.NombreItem == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("NombreItem sin especificar").Ctx(op, "required_sin_valor")
	}
	if tab.NombreItems == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("NombreItems sin especificar").Ctx(op, "required_sin_valor")
	}
	if tab.Abrev == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("Abrev sin especificar").Ctx(op, "required_sin_valor")
	}
	if tab.Humano == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("Humano sin especificar").Ctx(op, "required_sin_valor")
	}
	if tab.Kebab == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("Kebab sin especificar").Ctx(op, "required_sin_valor")
	}
	err := tab.Validar()
	if err != nil {
		return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec("INSERT INTO tablas "+
		"(tabla_id, paquete_id, nombre_repo, nombre_item, nombre_items, abrev, humano, humano_plural, kebab, es_femenino, descripcion, directrices) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ",
		tab.TablaID, tab.PaqueteID, tab.NombreRepo, tab.NombreItem, tab.NombreItems, tab.Abrev, tab.Humano, tab.HumanoPlural, tab.Kebab, tab.EsFemenino, tab.Descripcion, tab.Directrices,
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

// UpdateTabla valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdateTabla(tab ddd.Tabla) error {
	const op string = "mysqlddd.UpdateTabla"
	if tab.TablaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("TablaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if tab.NombreRepo == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("NombreRepo sin especificar").Ctx(op, "required_sin_valor")
	}
	if tab.NombreItem == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("NombreItem sin especificar").Ctx(op, "required_sin_valor")
	}
	if tab.NombreItems == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("NombreItems sin especificar").Ctx(op, "required_sin_valor")
	}
	if tab.Abrev == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("Abrev sin especificar").Ctx(op, "required_sin_valor")
	}
	if tab.Humano == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("Humano sin especificar").Ctx(op, "required_sin_valor")
	}
	if tab.Kebab == "" {
		return gecko.NewErr(http.StatusBadRequest).Msg("Kebab sin especificar").Ctx(op, "required_sin_valor")
	}
	err := tab.Validar()
	if err != nil {
		return gecko.NewErr(http.StatusBadRequest).Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec(
		"UPDATE tablas SET "+
			"tabla_id=?, paquete_id=?, nombre_repo=?, nombre_item=?, nombre_items=?, abrev=?, humano=?, humano_plural=?, kebab=?, es_femenino=?, descripcion=?, directrices=? "+
			"WHERE tabla_id = ?",
		tab.TablaID, tab.PaqueteID, tab.NombreRepo, tab.NombreItem, tab.NombreItems, tab.Abrev, tab.Humano, tab.HumanoPlural, tab.Kebab, tab.EsFemenino, tab.Descripcion, tab.Directrices,
		tab.TablaID,
	)
	if err != nil {
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/TBL-DELETE ====================================  //

func (s *Repositorio) DeleteTabla(TablaID int) error {
	const op string = "mysqlddd.DeleteTabla"
	if TablaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("TablaID sin especificar").Ctx(op, "pk_indefinida")
	}
	// Verificar que solo se borre un registro.
	var num int
	err := s.db.QueryRow("SELECT COUNT(tabla_id) FROM tablas WHERE tabla_id = ?",
		TablaID,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gecko.NewErr(http.StatusNotFound).Err(ddd.ErrTablaNotFound).Op(op)
		}
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	if num > 1 {
		return gecko.NewErr(http.StatusInternalServerError).Err(nil).Op(op).Msgf("abortado porque serían borrados %v registros", num)
	} else if num == 0 {
		return gecko.NewErr(http.StatusNotFound).Err(ddd.ErrTablaNotFound).Op(op).Msg("cero resultados")
	}
	// Eliminar todo lo relacionado a la tabla.
	_, err = s.db.Exec(
		"DELETE FROM valores_enum WHERE campo_id IN (SELECT v.campo_id FROM valores_enum v INNER JOIN campos c ON v.campo_id = c.campo_id INNER JOIN tablas t ON c.tabla_id = t.tabla_id WHERE t.tabla_id = ? GROUP BY v.campo_id)",
		TablaID,
	)
	if err != nil {
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op).Op("delete_valores_enum")
	}
	_, err = s.db.Exec(
		"DELETE FROM campos WHERE tabla_id = ?",
		TablaID,
	)
	if err != nil {
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op).Op("delete_campos")
	}
	_, err = s.db.Exec(
		"DELETE FROM tablas WHERE tabla_id = ?",
		TablaID,
	)
	if err != nil {
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op).Op("delete_tabla")
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROW ======================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowTabla(row *sql.Row, tab *ddd.Tabla, op string) error {

	err := row.Scan(
		&tab.TablaID, &tab.PaqueteID, &tab.NombreRepo, &tab.NombreItem, &tab.NombreItems, &tab.Abrev, &tab.Humano, &tab.HumanoPlural, &tab.Kebab, &tab.EsFemenino, &tab.Descripcion, &tab.Directrices,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gecko.NewErr(http.StatusNotFound).Msg("Tabla no se encuentra").Op(op)
		}
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}

	return nil
}

//  ================================================================  //
//  ========== MYSQL/GET ===========================================  //

// GetTabla devuelve un Tabla de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetTabla(TablaID int) (*ddd.Tabla, error) {
	const op string = "mysqlddd.GetTabla"
	if TablaID == 0 {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("TablaID sin especificar").Ctx(op, "pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasTabla+" "+fromTabla+
			"WHERE tabla_id = ?",
		TablaID,
	)
	tab := &ddd.Tabla{}
	return tab, s.scanRowTabla(row, tab, op)
}

//  ================================================================  //
//  ========== MYSQL/GET_BY ========================================  //

// GetTablaByNombreRepo devuelve un Tabla de la DB.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetTablaByNombreRepo(NombreRepo string) (*ddd.Tabla, error) {
	const op string = "mysqlddd.GetTablaByNombreRepo"
	if NombreRepo == "" {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("NombreRepo sin especificar").Ctx(op, "param_indefinido")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasTabla+" "+fromTabla+
			"WHERE nombre_repo LIKE ?",
		NombreRepo+"%",
	)
	tab := &ddd.Tabla{}
	return tab, s.scanRowTabla(row, tab, op)
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROWS =====================================  //

// scanRowsTabla escanea cada row en la struct Tabla
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsTabla(rows *sql.Rows, op string) ([]ddd.Tabla, error) {
	defer rows.Close()
	items := []ddd.Tabla{}
	for rows.Next() {
		tab := ddd.Tabla{}

		err := rows.Scan(
			&tab.TablaID, &tab.PaqueteID, &tab.NombreRepo, &tab.NombreItem, &tab.NombreItems, &tab.Abrev, &tab.Humano, &tab.HumanoPlural, &tab.Kebab, &tab.EsFemenino, &tab.Descripcion, &tab.Directrices,
		)
		if err != nil {
			return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}

		items = append(items, tab)
	}
	return items, nil
}

//  ================================================================  //
//  ========== MYSQL/LIST ==========================================  //

func (s *Repositorio) ListTablas() ([]ddd.Tabla, error) {
	const op string = "mysqlddd.ListTablas"
	rows, err := s.db.Query(
		"SELECT " + columnasTabla + " " + fromTabla,
	)
	if err != nil {
		return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return s.scanRowsTabla(rows, op)
}

//  ================================================================  //
//  ========== MYSQL/LIST_BY =======================================  //

// ListTablasByPaqueteID retorna los registros a partir de PaqueteID.
func (s *Repositorio) ListTablasByPaqueteID(PaqueteID int) ([]ddd.Tabla, error) {
	const op string = "mysqlddd.ListTablasByPaqueteID"
	if PaqueteID == 0 {
		return nil, gecko.NewErr(http.StatusBadRequest).Msg("PaqueteID sin especificar").Ctx(op, "param_indefinido")
	}
	where := "WHERE paquete_id = ?"
	argumentos := []any{}
	argumentos = append(argumentos, PaqueteID)

	rows, err := s.db.Query(
		"SELECT "+columnasTabla+" "+fromTabla+
			where,
		argumentos...,
	)
	if err != nil {
		return nil, gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	return s.scanRowsTabla(rows, op)
}
