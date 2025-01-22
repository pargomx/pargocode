package main

import (
	"fmt"
	"monorepo/appdominio"
	"monorepo/ddd"
	"monorepo/sqliteddd"
	"time"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
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
	err := appdominio.CrearConsulta(con, s.ddd)
	if err != nil {
		return err
	}
	gko.LogOkeyf("Consulta %v creada %v", con.ConsultaID, time.Now().Format("03:04:05pm"))
	return c.Redirf("/consultas/%v", con.ConsultaID)
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
	err := appdominio.ActualizarConsulta(con.ConsultaID, con, s.ddd)
	if err != nil {
		return err
	}
	return c.StatusOk(fmt.Sprintf("Consulta %v actualizada %v", con.ConsultaID, time.Now().Format("03:04:05pm")))
}

func (s *servidor) deleteConsulta(c *gecko.Context) error {
	err := appdominio.EliminarConsulta(c.PathInt("consulta_id"), s.ddd)
	if err != nil {
		return err
	}
	gko.LogEvento("consulta eliminada")
	return c.Redir("/paquetes")
}

// ================================================================ //
// ========== RELACIONES ========================================== //

func (s *servidor) postRelacionConsulta(c *gecko.Context) error {
	err := appdominio.AgregarRelacionConsulta(
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
	err := appdominio.EliminarRelacionConsulta(c.PathInt("consulta_id"), c.PathInt("posicion"), s.ddd)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) actualizarRelacionConsulta(c *gecko.Context) error {
	err := appdominio.ActualizarRelacionConsulta(
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
	err := appdominio.AgregarCampoConsulta(c.PathInt("consulta_id"), c.FormVal("from_abrev"), c.FormVal("expresion"), s.ddd)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *servidor) eliminarCampoConsulta(c *gecko.Context) error {
	err := appdominio.EliminarCampoConsulta(c.PathInt("consulta_id"), c.PathInt("posicion"), s.ddd)
	if err != nil {
		return err
	}
	return s.getConsulta(c)
}

func (s *servidor) reordenarCampoConsulta(c *gecko.Context) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	err = appdominio.ReordenarCampoConsulta(c.PathInt("consulta_id"), c.FormInt("oldPosicion"), c.FormInt("newPosicion"), sqliteddd.NuevoRepositorio(tx))
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
	err := appdominio.ActualizarCampoConsulta(cam, s.ddd)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}
