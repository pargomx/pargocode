package main

import (
	"html"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== TABLA =============================================== //

func (s *servidor) generarDeTabla(c *gecko.Context) error {
	gen, err := s.generador.DeTabla(c.PathInt("tabla_id"))
	if err != nil {
		return err
	}
	codigo, err := gen.GenerarToString(c.QueryVal("tipo"))
	if err != nil {
		return err
	}
	if c.EsHTMX() {
		return c.StatusOk(html.EscapeString(codigo))
	}
	return c.StatusOk(codigo)
}

func (s *servidor) generarDeTablaArchivos(c *gecko.Context) error {
	gen, err := s.generador.DeTabla(c.PathInt("tabla_id"))
	if err != nil {
		return err
	}
	err = gen.GenerarToFile(c.PathVal("tipo"))
	if err != nil {
		return err
	}
	return c.StatusOk("Código generado en " + gen.GetInfoDestino(c.PathVal("tipo")))
}

// ================================================================ //
// ========== CONSULTA ============================================ //

func (s *servidor) generarDeConsulta(c *gecko.Context) error {
	Consulta, err := s.generador.GetConsulta(c.PathInt("consulta_id"))
	if err != nil {
		return err
	}

	if c.QueryVal("modo") == "archivo" {
		err = s.generador.QryGenerarArchivos(Consulta, c.QueryVal("tipo")).Generar()
		if err != nil {
			return err
		}
		return c.StatusOk("Generado")
	}

	codigo, err := s.generador.GenerarDeConsultaStringNew(Consulta, c.QueryVal("tipo"))
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
	tablas, consultas, err := s.generador.GetTablasYConsultas(paq.PaqueteID)
	if err != nil {
		return err
	}
	for _, tbl := range tablas {
		reporte += tbl.GetInfoDestino(c.PathVal("tipo")) + "\n"
		err := tbl.GenerarToFile(c.PathVal("tipo"))
		if err != nil {
			errores = append(errores, err)
		}
	}
	for _, con := range consultas {
		call := s.generador.QryGenerarArchivos(&con, c.PathVal("tipo"))
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
