package main

import (
	"fmt"
	"monorepo/domain_driven_design/ddd"
	"monorepo/domain_driven_design/dpaquete"
	"monorepo/domain_driven_design/sqliteddd"
	"time"

	"github.com/pargomx/gecko"
)

func (s *servidor) crearConsulta(c *gecko.Context) error {
	con := ddd.Consulta{
		ConsultaID:  ddd.NewConsultaID(),
		PaqueteID:   c.FormInt("paquete_id"),
		TablaID:     c.FormInt("tabla_id"),
		NombreItem:  c.FormVal("nombre_item"),
		NombreItems: c.FormVal("nombre_items"),
		Abrev:       c.FormVal("abrev"),
		EsFemenino:  c.FormBool("es_femenino"),
		Descripcion: c.FormVal("descripcion"),
		Directrices: c.FormValue("directrices"),
	}
	err := dpaquete.CrearConsulta(con, s.ddd)
	if err != nil {
		return err
	}
	c.LogOkeyf("Consulta %v creada %v", con.ConsultaID, time.Now().Format("03:04:05pm"))
	return c.RedirectHTMX("/consultas/%v", con.ConsultaID)
}

func (s *servidor) actualizarConsulta(c *gecko.Context) error {
	con := ddd.Consulta{
		ConsultaID:  c.PathInt("consulta_id"),
		PaqueteID:   c.FormInt("paquete_id"),
		TablaID:     c.FormInt("tabla_id"),
		NombreItem:  c.FormVal("nombre_item"),
		NombreItems: c.FormVal("nombre_items"),
		Abrev:       c.FormVal("abrev"),
		EsFemenino:  c.FormBool("es_femenino"),
		Descripcion: c.FormVal("descripcion"),
		Directrices: c.FormValue("directrices"),
	}
	err := dpaquete.ActualizarConsulta(con.ConsultaID, con, s.ddd)
	if err != nil {
		return err
	}
	return c.StatusOk(fmt.Sprintf("Consulta %v actualizada %v", con.ConsultaID, time.Now().Format("03:04:05pm")))
}

func (s *servidor) deleteConsulta(c *gecko.Context) error {
	err := dpaquete.EliminarConsulta(c.PathInt("consulta_id"), s.ddd)
	if err != nil {
		return err
	}
	c.LogEvento("consulta eliminada")
	return c.Redir("/paquetes")
}

// ================================================================ //
// ========== RELACIONES ========================================== //

func (s *servidor) postRelacionConsulta(c *gecko.Context) error {
	err := dpaquete.AgregarRelacionConsulta(
		c.PathInt("consulta_id"),
		c.FormVal("tipo_join"),
		c.FormInt("join_tabla_id"),
		c.FormVal("from_abrev"),
		s.ddd)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) eliminarRelacionConsulta(c *gecko.Context) error {
	err := dpaquete.EliminarRelacionConsulta(c.PathInt("consulta_id"), c.PathInt("posicion"), s.ddd)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) actualizarRelacionConsulta(c *gecko.Context) error {
	err := dpaquete.ActualizarRelacionConsulta(
		c.PathInt("consulta_id"),
		c.PathInt("posicion"),
		c.FormVal("tipo_join"),
		c.FormVal("join_as"),
		c.FormVal("join_on"),
		s.ddd)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

// ================================================================ //
// ========== CAMPOS ============================================== //

func (s *servidor) postCampoConsulta(c *gecko.Context) error {
	err := dpaquete.AgregarCampoConsulta(c.PathInt("consulta_id"), c.FormVal("from_abrev"), c.FormVal("expresion"), s.ddd)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) eliminarCampoConsulta(c *gecko.Context) error {
	err := dpaquete.EliminarCampoConsulta(c.PathInt("consulta_id"), c.PathInt("posicion"), s.ddd)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) reordenarCampoConsulta(c *gecko.Context) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	err = dpaquete.ReordenarCampoConsulta(c.PathInt("consulta_id"), c.FormInt("oldPosicion"), c.FormInt("newPosicion"), sqliteddd.NuevoRepositorio(tx))
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return c.StatusOk("Campo reordenado")
}

func (s *servidor) actualizarCampoConsulta(c *gecko.Context) error {
	cam := ddd.ConsultaCampo{
		ConsultaID:  c.PathInt("consulta_id"),
		Posicion:    c.PathInt("posicion"),
		Expresion:   c.FormVal("expresion"),
		AliasSql:    c.FormVal("alias_sql"),
		NombreCampo: c.FormVal("nombre_campo"),
		TipoGo:      c.FormVal("tipo_go"),
		Pk:          c.FormBool("pk"),
		Filtro:      c.FormBool("filtro"),
		GroupBy:     c.FormBool("group_by"),
		Descripcion: c.FormVal("descripcion"),
	}
	err := dpaquete.ActualizarCampoConsulta(cam, s.ddd)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}
