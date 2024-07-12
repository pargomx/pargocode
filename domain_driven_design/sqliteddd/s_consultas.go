package sqliteddd

import (
	"database/sql"
	"monorepo/domain_driven_design/ddd"
	"net/http"
	"strings"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== RELACIONES ========================================== //

// DeleteRelacionConsulta elimina permanentemente un registro de la consulta relación de la base de datos.
// Error si el registro no existe o si no se da la clave primaria.
func (s *Repositorio) DeleteRelacionConsulta(ConsultaID int, Posicion int) error {
	const op string = "mysqlddd.DeleteRelacionConsulta"
	if ConsultaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("ConsultaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if Posicion == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("Posicion sin especificar").Ctx(op, "pk_indefinida")
	}
	// Verificar que solo se borre un registro.
	var num int
	err := s.db.QueryRow("SELECT COUNT(consulta_id) FROM consulta_relaciones WHERE consulta_id = ? AND posicion = ?",
		ConsultaID, Posicion,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gecko.NewErr(http.StatusNotFound).Err(ddd.ErrConsultaRelacionNotFound).Op(op)
		}
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	if num > 1 {
		return gecko.NewErr(http.StatusInternalServerError).Err(nil).Op(op).Msgf("abortado porque serían borrados %v registros", num)
	} else if num == 0 {
		return gecko.NewErr(http.StatusNotFound).Err(ddd.ErrConsultaRelacionNotFound).Op(op).Msg("cero resultados")
	}
	// Eliminar registro
	_, err = s.db.Exec(
		"DELETE FROM consulta_relaciones WHERE consulta_id = ? AND posicion = ?",
		ConsultaID, Posicion,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1451 (23000)") {
			return gecko.NewErr(http.StatusConflict).Err(err).Op(op).Msg("Este registro es referenciado por otros y no se puede eliminar")
		} else {
			return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}
	}
	// Actualizar hermanos
	_, err = s.db.Exec(
		"UPDATE consulta_relaciones SET posicion = posicion - 1 WHERE consulta_id = ? AND posicion > ?",
		ConsultaID, Posicion,
	)
	if err != nil {
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op).Op("actualizar_hermanos")
	}
	return nil
}

//  ================================================================  //
//  ========== CAMPOS DE CONSULTA ==================================  //

// DeleteConsultaCampo elimina permanentemente un registro del campo de consulta de la base de datos.
// Error si el registro no existe o si no se da la clave primaria.
func (s *Repositorio) DeleteConsultaCampo(ConsultaID int, Posicion int) error {
	const op string = "mysqlddd.DeleteConsultaCampo"
	if ConsultaID == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("ConsultaID sin especificar").Ctx(op, "pk_indefinida")
	}
	if Posicion == 0 {
		return gecko.NewErr(http.StatusBadRequest).Msg("Posicion sin especificar").Ctx(op, "pk_indefinida")
	}
	// Verificar que solo se borre un registro.
	var num int
	err := s.db.QueryRow("SELECT COUNT(consulta_id) FROM consulta_campos WHERE consulta_id = ? AND posicion = ?",
		ConsultaID, Posicion,
	).Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return gecko.NewErr(http.StatusNotFound).Err(ddd.ErrConsultaCampoNotFound).Op(op)
		}
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
	}
	if num > 1 {
		return gecko.NewErr(http.StatusInternalServerError).Err(nil).Op(op).Msgf("abortado porque serían borrados %v registros", num)
	} else if num == 0 {
		return gecko.NewErr(http.StatusNotFound).Err(ddd.ErrConsultaCampoNotFound).Op(op).Msg("cero resultados")
	}
	// Eliminar registro
	_, err = s.db.Exec(
		"DELETE FROM consulta_campos WHERE consulta_id = ? AND posicion = ?",
		ConsultaID, Posicion,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1451 (23000)") {
			return gecko.NewErr(http.StatusConflict).Err(err).Op(op).Msg("Este registro es referenciado por otros y no se puede eliminar")
		} else {
			return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op)
		}
	}
	// Actualizar hermanos
	_, err = s.db.Exec(
		"UPDATE consulta_campos SET posicion = posicion - 1 WHERE consulta_id = ? AND posicion > ?",
		ConsultaID, Posicion,
	)
	if err != nil {
		return gecko.NewErr(http.StatusInternalServerError).Err(err).Op(op).Op("actualizar_hermanos")
	}
	return nil
}

// ================================================================ //
// ========== REORDENAR CAMPO ===================================== //

func (s *Repositorio) ReordenarCampoConsulta(consultaID int, oldPosicion int, newPosicion int) error {
	var op = gecko.NewOp("mysqlddd.ReordenarCampoConsulta")
	if consultaID == 0 {
		return op.Msg("consultaID sin especificar")
	}
	if oldPosicion < 1 {
		return op.Msg("oldPosicion inválida")
	}
	if newPosicion < 1 {
		return op.Msg("newPosicion inválida")
	}
	if oldPosicion == newPosicion {
		return nil
	}
	var err error

	// Cambiar posición del campo a negativa para distinguir de los hermanos.
	_, err = s.db.Exec(
		"UPDATE consulta_campos SET posicion = -(?) WHERE consulta_id = ? AND posicion = ?",
		newPosicion, consultaID, oldPosicion,
	)
	if err != nil {
		return op.Err(err).Op("posicion_nueva_negativa")
	}
	// Cambiar posición de hermanos
	if oldPosicion < newPosicion {
		_, err = s.db.Exec(
			"UPDATE consulta_campos SET posicion = -(posicion - 1) WHERE consulta_id = ? AND posicion > ? AND posicion <= ?",
			consultaID, oldPosicion, newPosicion,
		)
	} else {
		_, err = s.db.Exec(
			"UPDATE consulta_campos SET posicion = -(posicion + 1) WHERE consulta_id = ? AND posicion >= ? AND posicion <= ?",
			consultaID, newPosicion, oldPosicion,
		)
	}
	if err != nil {
		return op.Err(err).Op("posicion_hermanos")
	}
	// Volver la posición a positiva.
	_, err = s.db.Exec(
		"UPDATE consulta_campos SET posicion = -(posicion) WHERE consulta_id = ? AND posicion < 0",
		consultaID,
	)
	if err != nil {
		return op.Err(err).Op("posicion_nueva")
	}
	return nil
}
