package main

import (
	"html"
	"monorepo/codegenerator"
	"strings"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== TABLA =============================================== //

func (s *servidor) generarDeTabla(c *gecko.Context) error {
	gen, err := codegenerator.NuevoDeTabla(s.ddd, c.PathInt("tabla_id"))
	if err != nil {
		return err
	}
	err = gen.PrepararJob(c.FormVal("tipo")).Generar()
	if err != nil {
		return err
	}
	if c.QueryVal("modo") == "archivo" {
		err = gen.ToFile()
		if err != nil {
			return err
		}
		return c.StatusOk(strings.Join(gen.GetHechos(), "\n"))

	} else if c.EsHTMX() {
		return c.StatusOk(html.EscapeString(gen.ToString()))

	} else {
		return c.StatusOk(gen.ToString())
	}
}

// ================================================================ //
// ========== CONSULTA ============================================ //

func (s *servidor) generarDeConsulta(c *gecko.Context) error {
	gen, err := codegenerator.NuevoDeConsulta(s.ddd, c.PathInt("consulta_id"))
	if err != nil {
		return err
	}
	err = gen.PrepararJob(c.FormVal("tipo")).Generar()
	if err != nil {
		return err
	}
	if c.QueryVal("modo") == "archivo" {
		err = gen.ToFile()
		if err != nil {
			return err
		}
		return c.StatusOk(strings.Join(gen.GetHechos(), "\n"))

	} else if c.EsHTMX() {
		return c.StatusOk(html.EscapeString(gen.ToString()))

	} else {
		return c.StatusOk(gen.ToString())
	}
}

// ================================================================ //
// ========== PAQUETE ============================================= //

func (s *servidor) generarDePaqueteArchivos(c *gecko.Context) error {
	generadores, err := codegenerator.NuevoDePaquete(s.ddd, c.PathInt("paquete_id"))
	if err != nil {
		return err
	}
	reporte := "ARCHIVOS GENERADOS:\n\n"
	errores := []error{}
	for _, gen := range generadores {
		err = gen.PrepararJob(c.FormVal("tipo")).Generar()
		if err != nil {
			errores = append(errores, err)
			continue
		}
		err := gen.ToFile()
		if err != nil {
			errores = append(errores, err)
		}
		reporte += strings.Join(gen.GetHechos(), "\n") + "\n"
	}
	if len(errores) > 0 {
		reporte += "\nERRORES:\n\n"
		for _, e := range errores {
			reporte += e.Error() + "\n\n"
		}
	}
	return c.StatusOk(reporte)
}
