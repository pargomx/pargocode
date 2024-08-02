package main

import (
	"html"
	"monorepo/codegenerator"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== TABLA =============================================== //

func (s *servidor) generarDeTabla(c *gecko.Context) error {
	tbl, err := s.codeGenerator.GetTabla(c.PathInt("tabla_id"))
	if err != nil {
		return err
	}
	codigo, err := tbl.GenerarDeTablaString(c.QueryVal("tipo"))
	if err != nil {
		return err
	}
	if c.EsHTMX() {
		return c.StatusOk(html.EscapeString(codigo))
	}
	return c.StatusOk(codigo)
}

func (s *servidor) generarDeTablaArchivos(c *gecko.Context) error {
	tbl, err := s.codeGenerator.GetTabla(c.PathInt("tabla_id"))
	if err != nil {
		return err
	}
	call := tbl.TblGenerarArchivos(c.PathVal("tipo"))
	err = call.Generar()
	if err != nil {
		return err
	}
	return c.StatusOk("Código generado en " + call.Destino())
}

// ================================================================ //
// ========== CONSULTA ============================================ //

func (s *servidor) generarDeConsulta(c *gecko.Context) error {
	Consulta, err := codegenerator.GetConsulta(c.PathInt("consulta_id"), s.ddd)
	if err != nil {
		return err
	}

	if c.QueryVal("modo") == "archivo" {
		err = s.codeGenerator.QryGenerarArchivos(Consulta, c.QueryVal("tipo")).Generar()
		if err != nil {
			return err
		}
		return c.StatusOk("Generado")
	}

	codigo, err := s.codeGenerator.GenerarDeConsultaStringNew(Consulta, c.QueryVal("tipo"))
	if err != nil {
		return err
	}
	if c.EsHTMX() {
		return c.StatusOk(html.EscapeString(codigo))
	}
	return c.StatusOk(codigo)
}

// ================================================================ //
// ================================================================ //

func (s *servidor) generarDePaqueteArchivos(c *gecko.Context) error {
	paq, err := s.ddd.GetPaquete(c.PathInt("paquete_id"))
	if err != nil {
		return err
	}
	reporte := "ARCHIVOS GENERADOS:\n\n"
	errores := []error{}
	tablas, consultas, err := codegenerator.GetTablasYConsultas(paq.PaqueteID, s.ddd)
	if err != nil {
		return err
	}
	for _, tbl := range tablas {
		call := tbl.TblGenerarArchivos(c.PathVal("tipo"))
		reporte += call.Destino() + "\n"
		err = call.Generar()
		if err != nil {
			errores = append(errores, err)
		}
	}
	for _, con := range consultas {
		call := s.codeGenerator.QryGenerarArchivos(&con, c.PathVal("tipo"))
		reporte += call.Destino() + "\n"
		err = call.Generar()
		if err != nil {
			errores = append(errores, err)
		}
	}
	if len(errores) > 0 {
		reporte += "\nERRORES:\n\n"
		for _, e := range errores {
			reporte += e.Error() + "\n\n"
		}
	}
	return c.StatusOk(reporte)
}
