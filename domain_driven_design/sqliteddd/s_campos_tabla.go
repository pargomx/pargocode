package sqliteddd

import (
	"monorepo/domain_driven_design/ddd"

	"github.com/pargomx/gecko/gko"
)

// ListCamposByTablaID retorna los registros a partir de TablaID.
func (s *Repositorio) ListCamposByTablaID(TablaID int) ([]ddd.Campo, error) {
	const op string = "mysqlddd.ListCamposByTablaID"
	if TablaID == 0 {
		return nil, gko.ErrDatoInvalido().Msg("TablaID sin especificar").Ctx(op, "param_indefinido")
	}
	rows, err := s.db.Query(
		"SELECT "+columnasCampo+" "+fromCampo+"WHERE tabla_id = ? ORDER BY posicion",
		TablaID,
	)
	if err != nil {
		return nil, gko.ErrInesperado().Err(err).Op(op)
	}
	return s.scanRowsCampo(rows, op)
}

func (s *Repositorio) GetCampoPrimaryKey(nombreColumna string) (*ddd.Campo, error) {
	const op string = "mysqlddd.GetCampoPrimaryKey"
	if nombreColumna == "" {
		return nil, gko.ErrDatoInvalido().Msg("nombreColumna sin especificar").Ctx(op, "pk_indefinida")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasCampo+" "+fromCampo+
			"WHERE primary_key = 1 AND foreign_key = 0 AND nombre_columna = ?",
		nombreColumna,
	)
	cam := &ddd.Campo{}
	return cam, s.scanRowCampo(row, cam, op)
}

func (s *Repositorio) GetCampoByNombre(nombre string) (*ddd.Campo, error) {
	op := gko.Op("mysqlddd.GetCampoByNombre")
	if nombre == "" {
		return nil, op.Msg("nombre sin especificar")
	}
	row := s.db.QueryRow(
		"SELECT "+columnasCampo+" "+fromCampo+
			"WHERE nombre_columna = ? OR nombre_campo = ? OR nombre_humano = ? ORDER BY primary_key DESC, foreign_key ASC LIMIT 1",
		nombre, nombre, nombre,
	)
	cam := &ddd.Campo{}
	err := s.scanRowCampo(row, cam, "mysqlddd.GetCampoByNombre")
	if err != nil {
		return nil, op.Err(err)
	}
	return cam, nil
}

// ================================================================ //
// ================================================================ //

func (s *Repositorio) ReordenarCampo(cam *ddd.Campo, newPosicion int) error {
	const op string = "mysqlddd.ReordenarCampo"
	if cam.CampoID == 0 {
		return gko.ErrDatoInvalido().Msg("CampoID sin especificar").Ctx(op, "pk_indefinida")
	}
	if cam.TablaID == 0 {
		return gko.ErrDatoInvalido().Msg("TablaID sin especificar").Ctx(op, "fk_indefinida")
	}
	if cam.Posicion == newPosicion {
		return nil
	}

	// Validar nueva posición
	var hermanos int
	err := s.db.QueryRow("SELECT COUNT(campo_id) FROM campos WHERE tabla_id = ?", cam.TablaID).Scan(&hermanos)
	if err != nil {
		return gko.ErrInesperado().Err(err).Op(op).Op("contar_hermanos")
	}
	if newPosicion < 1 || newPosicion > hermanos {
		return gko.ErrDatoInvalido().Msg("Posición fuera de rango").Op(op).Ctx("posicion_invalida", newPosicion).Ctx("hermanos", hermanos)
	}

	// Cambiar posición de hermanos
	if cam.Posicion < newPosicion {
		_, err = s.db.Exec(
			"UPDATE campos SET posicion = (posicion - 1) WHERE posicion > ? AND posicion <= ?",
			cam.Posicion, newPosicion,
		)
	} else {
		_, err = s.db.Exec(
			"UPDATE campos SET posicion = (posicion + 1) WHERE posicion >= ? AND posicion <= ?",
			newPosicion, cam.Posicion,
		)
	}
	if err != nil {
		return gko.ErrInesperado().Err(err).Op(op).Op("posicion_hermanos")
	}

	// Cambiar posición del campo
	_, err = s.db.Exec(
		"UPDATE campos SET posicion = ? WHERE campo_id = ?",
		newPosicion, cam.CampoID,
	)
	if err != nil {
		return gko.ErrInesperado().Err(err).Op(op).Op("posicion_nueva")
	}
	return nil
}
