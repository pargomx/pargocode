package main

import (
	"fmt"
	"html"
	"monorepo/domain_driven_design/ddd"
	"monorepo/domain_driven_design/dpaquete"
	"monorepo/textutils"
	"time"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *servidor) postTablaNueva(c *gecko.Context) error {
	tbl := ddd.Tabla{
		TablaID:      ddd.NewTablaID(),
		PaqueteID:    c.FormInt("paquete_id"),
		Humano:       textutils.PrimeraMayusc(c.FormVal("etiqueta")),
		HumanoPlural: textutils.PrimeraMayusc(c.FormVal("etiqueta_plural")),
		NombreRepo:   c.FormVal("nombre_tabla"),
		NombreItem:   c.FormVal("nombre_item"),
		NombreItems:  c.FormVal("nombre_items"),
		Kebab:        c.FormVal("kebab"),
		Abrev:        c.FormVal("abrev"),
		EsFemenino:   c.FormBool("nombre_femenino"),
		Descripcion:  c.FormVal("descripcion"),
		Directrices:  c.FormValue("directrices"),
	}
	err := dpaquete.AgregarTabla(tbl, s.ddd)
	if err != nil {
		return err
	}
	gko.LogInfof("Tabla nueva '%s'", tbl.NombreRepo)
	return c.RedirectHTMX("/tablas/%d", tbl.TablaID)
}

// ================================================================ //
// ========== ACTUALIZAR ========================================== //

func (s *servidor) putTabla(c *gecko.Context) error {
	tbl := ddd.Tabla{
		TablaID:      c.PathInt("tabla_id"),
		PaqueteID:    c.FormInt("paquete_id"),
		Humano:       textutils.PrimeraMayusc(c.FormVal("etiqueta")),
		HumanoPlural: textutils.PrimeraMayusc(c.FormVal("etiqueta_plural")),
		NombreRepo:   c.FormVal("nombre_tabla"),
		NombreItem:   c.FormVal("nombre_item"),
		NombreItems:  c.FormVal("nombre_items"),
		Kebab:        c.FormVal("kebab"),
		Abrev:        c.FormVal("abrev"),
		EsFemenino:   c.FormBool("nombre_femenino"),
		Descripcion:  c.FormVal("descripcion"),
		Directrices:  c.FormValue("directrices"),
	}
	err := dpaquete.ActualizarTabla(tbl.TablaID, tbl, s.ddd)
	if err != nil {
		return err
	}
	gko.LogInfof("Tabla actualizada '%s'", tbl.NombreRepo)
	return c.StatusOk(fmt.Sprintf("Guardado %v", time.Now().Format("03:04:05pm")))
}

func (s *servidor) eliminarTabla(c *gecko.Context) error {
	tbl, err := s.ddd.GetTabla(c.PathInt("tabla_id"))
	if err != nil {
		return err
	}
	if c.PromptVal() != "ok" {
		return gko.ErrDatoInvalido().Msg("Para eliminarlo escribe el 'ok' en el campo de confirmación")
	}
	err = s.ddd.DeleteTabla(tbl.TablaID)
	if err != nil {
		return err
	}
	gko.LogInfof("Tabla '%s' eliminada", tbl.NombreRepo)
	return c.Redir("/paquetes")
}

// ================================================================ //
// ========== GENERAR ============================================= //

func (s *servidor) generarDeTabla(c *gecko.Context) error {
	tbl, err := dpaquete.GetTabla(c.PathInt("tabla_id"), s.ddd)
	if err != nil {
		return err
	}
	codigo, err := codeGenerator.GenerarDeTablaString(tbl, c.QueryVal("tipo"))
	if err != nil {
		return err
	}
	if c.EsHTMX() {
		return c.StatusOk(html.EscapeString(codigo))
	}
	return c.StatusOk(codigo)
}

func (s *servidor) generarDeTablaArchivos(c *gecko.Context) error {
	tbl, err := dpaquete.GetTabla(c.PathInt("tabla_id"), s.ddd)
	if err != nil {
		return err
	}
	call := codeGenerator.TblGenerarArchivos(tbl, c.PathVal("tipo"))
	err = call.Generar()
	if err != nil {
		return err
	}
	return c.StatusOk("Código generado en " + call.Destino())
}
