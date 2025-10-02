package main

import (
	"strings"

	"github.com/pargomx/pargocode/appdominio"
	"github.com/pargomx/pargocode/ddd"
	"github.com/pargomx/pargocode/textutils"

	"github.com/pargomx/gecko"
)

func (s *servidor) formNuevaConsulta(c *gecko.Context) error {
	if c.EsHTMX() {
		return c.RedirFullf("/consultas/nueva?nombre=%v&paquete_id=%v", c.Request().Header.Get("HX-Prompt"), c.QueryInt("paquete_id"))
	}
	paquetes, err := s.ddd.ListPaquetes()
	if err != nil {
		return err
	}
	tablas, err := s.ddd.ListTablas()
	if err != nil {
		return err
	}
	con := ddd.Consulta{
		PaqueteID: c.QueryInt("paquete_id"),
	}
	// Facilitar con sugerencias los derivados del nombre de la consulta.
	nuevoNombre := c.QueryVal("nombre")
	if nuevoNombre != "" {
		con.NombreItem = textutils.PrimeraMayusc(textutils.QuitarAcentos(nuevoNombre))
		con.NombreItems = textutils.PrimeraMayusc(textutils.QuitarAcentos(textutils.DeducirNombrePlural(nuevoNombre)))
		con.NombreItem = textutils.KebabToCamel(strings.ReplaceAll(con.NombreItem, " ", "-"))
		con.NombreItems = textutils.KebabToCamel(strings.ReplaceAll(con.NombreItems, " ", "-"))
		if len(con.NombreItem) >= 3 {
			con.Abrev = strings.ToLower(con.NombreItem[:3])
		}
	}
	data := map[string]any{
		"Titulo":   "Nueva consulta",
		"Consulta": con,
		"Paquetes": paquetes,
		"Tablas":   tablas,
	}
	return c.RenderOk("ddd/consulta-nueva", data)
}

func (s *servidor) getConsulta(c *gecko.Context) error {
	agregadoConsulta, err := appdominio.GetConsulta(c.PathInt("consulta_id"), s.ddd)
	if err != nil {
		return err
	}
	paquetes, err := s.ddd.ListPaquetes()
	if err != nil {
		return err
	}
	tablas, err := s.ddd.ListTablas()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":           agregadoConsulta.Consulta.NombreItem + " (qry)",
		"AgregadoConsulta": agregadoConsulta,
		"Paquetes":         paquetes,
		"Tablas":           tablas,
	}
	return c.RenderOk("ddd/consulta", data)
}
