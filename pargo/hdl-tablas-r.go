package main

import (
	"monorepo/ddd"
	"monorepo/dpaquete"
	"monorepo/textutils"
	"strings"

	"github.com/pargomx/gecko"
)

func (s *servidor) getTabla(c *gecko.Context) error {
	Agregado, err := dpaquete.GetTabla(c.PathInt("tabla_id"), s.ddd)
	if err != nil {
		return err
	}
	// Agregar un campo vacío para mostrar formulario y poder agregar campos
	// Agregado.Campos = append(Agregado.Campos, dpaquete.CampoTabla{TablaID: Agregado.Tabla.TablaID})

	// Para poder cambiar de paquete
	paquetes, err := s.ddd.ListPaquetes()
	if err != nil {
		return err
	}

	data := map[string]any{
		"Titulo":   Agregado.Tabla.Humano,
		"Agregado": Agregado,
		"Paquetes": paquetes,
	}
	return c.RenderOk("ddd/tabla", data)
}

// ================================================================ //
// ========== NUEVA TABLA ========================================= //

func (s *servidor) getTablaNueva(c *gecko.Context) error {
	if c.EsHTMX() {
		return c.RedirectHTMX("/tablas/nueva?etiqueta=%v&paquete_id=%v", c.Request().Header.Get("HX-Prompt"), c.QueryInt("paquete_id"))
	}
	tbl := ddd.Tabla{
		PaqueteID: c.QueryInt("paquete_id"),
	}
	nuevoNombre := c.QueryVal("etiqueta") // Nuevo nombre humano
	if nuevoNombre != "" {
		tbl.Humano = textutils.PrimeraMayusc(nuevoNombre)
		tbl.HumanoPlural = textutils.PrimeraMayusc(textutils.DeducirNombrePlural(nuevoNombre))
		tbl.NombreRepo = textutils.KebabToSnake(textutils.QuitarAcentos(strings.ReplaceAll(textutils.DeducirNombrePlural(nuevoNombre), " ", "-")))
		tbl.NombreItem = textutils.KebabToCamel(textutils.QuitarAcentos(strings.ToLower(strings.ReplaceAll(nuevoNombre, " ", "-"))))
		tbl.NombreItems = textutils.KebabToCamel(textutils.QuitarAcentos(strings.ToLower(strings.ReplaceAll(textutils.DeducirNombrePlural(nuevoNombre), " ", "-"))))
		tbl.Kebab = textutils.QuitarAcentos(strings.ToLower(strings.ReplaceAll(nuevoNombre, " ", "-")))
		if len(tbl.NombreRepo) >= 3 {
			tbl.Abrev = tbl.NombreRepo[:3]
		}
		// tbl.Clave = textutils.QuitarAcentos(strings.ToLower(strings.ReplaceAll(nuevoNombre, " ", "")))
		// tbl.ClavePlural = textutils.DeducirNombrePlural(tbl.Clave)
		// tbl.KebabPlural = textutils.DeducirNombrePlural(tbl.Kebab)
	}
	paquetes, err := s.ddd.ListPaquetes()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":   "Tabla nueva",
		"Tabla":    tbl,
		"Paquetes": paquetes,
	}
	return c.RenderOk("ddd/tabla-nueva", data)
}
