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
	return c.RedirFullf("/consultas/%v", con.ConsultaID)
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
	return c.RedirFullf("/paquetes")
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
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	repoTx := sqliteddd.NuevoRepositorio(tx)
	err = appdominio.EliminarCampoConsulta(c.PathInt("consulta_id"), c.PathInt("posicion"), repoTx)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
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

func (s *servidor) campoConsModifExpresion(c *gecko.Context) error {
	campo, err := appdominio.CampoConsultaModifExpresion(appdominio.CampoConsultaModif{
		ConsultaID: c.PathInt("consulta_id"),
		Posicion:   c.PathInt("posicion"),
		Valor:      c.FormVal("expresion"),
	}, s.ddd, s.txt)
	if err != nil {
		return err
	}
	return c.Render(200, "consultas/campo-tr", campo)
}
func (s *servidor) campoConsModifAlias(c *gecko.Context) error {
	campo, err := appdominio.CampoConsultaModifAlias(appdominio.CampoConsultaModif{
		ConsultaID: c.PathInt("consulta_id"),
		Posicion:   c.PathInt("posicion"),
		Valor:      c.FormVal("alias_sql"),
	}, s.ddd, s.txt)
	if err != nil {
		return err
	}
	return c.Render(200, "consultas/campo-tr", campo)
}
func (s *servidor) campoConsModifNombre(c *gecko.Context) error {
	campo, err := appdominio.CampoConsultaModifNombre(appdominio.CampoConsultaModif{
		ConsultaID: c.PathInt("consulta_id"),
		Posicion:   c.PathInt("posicion"),
		Valor:      c.FormVal("nombre_campo"),
	}, s.ddd, s.txt)
	if err != nil {
		return err
	}
	return c.Render(200, "consultas/campo-tr", campo)
}
func (s *servidor) campoConsModifTipo(c *gecko.Context) error {
	campo, err := appdominio.CampoConsultaModifTipo(appdominio.CampoConsultaModif{
		ConsultaID: c.PathInt("consulta_id"),
		Posicion:   c.PathInt("posicion"),
		Valor:      c.FormVal("tipo_go"),
	}, s.ddd, s.txt)
	if err != nil {
		return err
	}
	return c.Render(200, "consultas/campo-tr", campo)
}
func (s *servidor) campoConsModifPK(c *gecko.Context) error {
	campo, err := appdominio.CampoConsultaModifPK(appdominio.CampoConsultaModifBool{
		ConsultaID: c.PathInt("consulta_id"),
		Posicion:   c.PathInt("posicion"),
		Valor:      c.FormBool("pk"),
	}, s.ddd, s.txt)
	if err != nil {
		return err
	}
	return c.Render(200, "consultas/campo-tr", campo)
}
func (s *servidor) campoConsModifFiltro(c *gecko.Context) error {
	campo, err := appdominio.CampoConsultaModifFiltro(appdominio.CampoConsultaModifBool{
		ConsultaID: c.PathInt("consulta_id"),
		Posicion:   c.PathInt("posicion"),
		Valor:      c.FormBool("filtro"),
	}, s.ddd, s.txt)
	if err != nil {
		return err
	}
	return c.Render(200, "consultas/campo-tr", campo)
}
func (s *servidor) campoConsModifGroup(c *gecko.Context) error {
	campo, err := appdominio.CampoConsultaModifGroup(appdominio.CampoConsultaModifBool{
		ConsultaID: c.PathInt("consulta_id"),
		Posicion:   c.PathInt("posicion"),
		Valor:      c.FormBool("group_by"),
	}, s.ddd, s.txt)
	if err != nil {
		return err
	}
	return c.Render(200, "consultas/campo-tr", campo)
}
func (s *servidor) campoConsModifDesc(c *gecko.Context) error {
	campo, err := appdominio.CampoConsultaModifDesc(appdominio.CampoConsultaModif{
		ConsultaID: c.PathInt("consulta_id"),
		Posicion:   c.PathInt("posicion"),
		Valor:      c.FormVal("descripcion"),
	}, s.ddd, s.txt)
	if err != nil {
		return err
	}
	return c.Render(200, "consultas/campo-tr", campo)
}
