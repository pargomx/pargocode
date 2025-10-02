package sqliteddd

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/pargomx/pargocode/ddd"

	"github.com/pargomx/gecko/gko"
)

//  ================================================================  //
//  ========== MYSQL/CONSTANTES ====================================  //

// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
const columnasPaquete string = "paquete_id, go_module, directorio, nombre, descripcion"

// Origen de los datos de ddd.Paquete
//
// FROM paquetes
const fromPaquete string = "FROM paquetes "

//  ================================================================  //
//  ========== MYSQL/TBL-INSERT ====================================  //

// InsertPaquete valida el registro y lo inserta en la base de datos.
func (s *Repositorio) InsertPaquete(paq ddd.Paquete) error {
	const op string = "mysqlddd.InsertPaquete"
	if paq.PaqueteID == 0 {
		return gko.ErrDatoInvalido.Msg("PaqueteID sin especificar").Ctx(op, "pk_indefinida")
	}
	err := paq.Validar()
	if err != nil {
		return gko.ErrDatoInvalido.Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec("INSERT INTO paquetes "+
		"(paquete_id, go_module, directorio, nombre, descripcion) "+
		"VALUES (?, ?, ?, ?, ?) ",
		paq.PaqueteID, paq.GoModule, paq.Directorio, paq.Nombre, paq.Descripcion,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062 (23000)") {
			return gko.ErrYaExiste.Err(err).Op(op)
		} else if strings.HasPrefix(err.Error(), "Error 1452 (23000)") {
			return gko.ErrDatoInvalido.Err(err).Op(op).Msg("No se puede insertar la información porque el registro asociado no existe")
		} else {
			return gko.ErrInesperado.Err(err).Op(op)
		}
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/TBL-UPDATE ====================================  //

// UpdatePaquete valida y sobreescribe el registro en la base de datos.
func (s *Repositorio) UpdatePaquete(paq ddd.Paquete) error {
	const op string = "mysqlddd.UpdatePaquete"
	if paq.PaqueteID == 0 {
		return gko.ErrDatoInvalido.Msg("PaqueteID sin especificar").Ctx(op, "pk_indefinida")
	}
	err := paq.Validar()
	if err != nil {
		return gko.ErrDatoInvalido.Err(err).Op(op).Msg(err.Error())
	}
	_, err = s.db.Exec(
		"UPDATE paquetes SET "+
			"paquete_id=?, go_module=?, directorio=?, nombre=?, descripcion=? "+
			"WHERE paquete_id = ?",
		paq.PaqueteID, paq.GoModule, paq.Directorio, paq.Nombre, paq.Descripcion,
		paq.PaqueteID,
	)
	if err != nil {
		return gko.ErrInesperado.Err(err).Op(op)
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/TBL-DELETE ====================================  //

// DeletePaquete elimina permanentemente un registro del paquete de la base de datos.
// Error si el registro no existe o si no se da la clave primaria.
func (s *Repositorio) DeletePaquete(PaqueteID int) error {
	const op string = "mysqlddd.DeletePaquete"
	if PaqueteID == 0 {
		return gko.ErrDatoInvalido.Msg("PaqueteID sin especificar").Ctx(op, "pk_indefinida")
	}
	// Verificar que solo se borre un registro.
	var num int
	err := s.db.QueryRow("SELECT COUNT(paquete_id) FROM paquetes WHERE paquete_id = ?",
		PaqueteID,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gko.ErrNoEncontrado.Err(ddd.ErrPaqueteNotFound).Op(op)
		}
		return gko.ErrInesperado.Err(err).Op(op)
	}
	if num > 1 {
		return gko.ErrInesperado.Err(nil).Op(op).Msgf("abortado porque serían borrados %v registros", num)
	} else if num == 0 {
		return gko.ErrNoEncontrado.Err(ddd.ErrPaqueteNotFound).Op(op).Msg("cero resultados")
	}
	// Eliminar registro
	_, err = s.db.Exec(
		"DELETE FROM paquetes WHERE paquete_id = ?",
		PaqueteID,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1451 (23000)") {
			return gko.ErrYaExiste.Err(err).Op(op).Msg("Este registro es referenciado por otros y no se puede eliminar")
		} else {
			return gko.ErrInesperado.Err(err).Op(op)
		}
	}
	return nil
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROW ======================================  //

// Utilizar luego de un sql.QueryRow(). No es necesario hacer row.Close()
func (s *Repositorio) scanRowPaquete(row *sql.Row, paq *ddd.Paquete, op string) error {

	err := row.Scan(
		&paq.PaqueteID, &paq.GoModule, &paq.Directorio, &paq.Nombre, &paq.Descripcion,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gko.ErrNoEncontrado.Msg("el paquete no se encuentra").Op(op)
		}
		return gko.ErrInesperado.Err(err).Op(op)
	}

	return nil
}

//  ================================================================  //
//  ========== MYSQL/GET ===========================================  //

// GetPaquete devuelve un Paquete de la DB por su clave primaria.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetPaquete(PaqueteID int) (*ddd.Paquete, error) {
	const op string = "mysqlddd.GetPaquete"
	if PaqueteID == 0 {
		return nil, gko.ErrDatoInvalido.Msg("PaqueteID sin especificar").Ctx(op, "pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasPaquete+" "+fromPaquete+
			"WHERE paquete_id = ?",
		PaqueteID,
	)
	paq := &ddd.Paquete{}
	return paq, s.scanRowPaquete(row, paq, op)
}

//  ================================================================  //
//  ========== MYSQL/GET_BY ========================================  //

// GetPaqueteByNombre devuelve un Paquete de la DB.
// Error si no encuentra ninguno, o si encuentra más de uno.
func (s *Repositorio) GetPaqueteByNombre(Nombre string) (*ddd.Paquete, error) {
	const op string = "mysqlddd.GetPaqueteByNombre"
	if Nombre == "" {
		return nil, gko.ErrDatoInvalido.Msg("Nombre sin especificar").Ctx(op, "param_indefinido")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasPaquete+" "+fromPaquete+
			"WHERE nombre = ?",
		Nombre,
	)
	paq := &ddd.Paquete{}
	return paq, s.scanRowPaquete(row, paq, op)
}

//  ================================================================  //
//  ========== MYSQL/SCAN-ROWS =====================================  //

// scanRowsPaquete escanea cada row en la struct Paquete
// y devuelve un slice con todos los items.
// Siempre se encarga de llamar rows.Close()
func (s *Repositorio) scanRowsPaquete(rows *sql.Rows, op string) ([]ddd.Paquete, error) {
	defer rows.Close()
	items := []ddd.Paquete{}
	for rows.Next() {
		paq := ddd.Paquete{}

		err := rows.Scan(
			&paq.PaqueteID, &paq.GoModule, &paq.Directorio, &paq.Nombre, &paq.Descripcion,
		)
		if err != nil {
			return nil, gko.ErrInesperado.Err(err).Op(op)
		}

		items = append(items, paq)
	}
	return items, nil
}

//  ================================================================  //
//  ========== MYSQL/LIST ==========================================  //

// ListPaquetes devuelve todos los registros de los paquetes
func (s *Repositorio) ListPaquetes() ([]ddd.Paquete, error) {
	const op string = "mysqlddd.ListPaquetes"
	rows, err := s.db.Query(
		"SELECT " + columnasPaquete + " " + fromPaquete,
	)
	if err != nil {
		return nil, gko.ErrInesperado.Err(err).Op(op)
	}
	return s.scanRowsPaquete(rows, op)
}
