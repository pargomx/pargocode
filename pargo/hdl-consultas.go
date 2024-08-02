package main

import (
	"html"
	"monorepo/codegenerator"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== GENERAR ============================================= //

func (s *servidor) generarDeConsulta(c *gecko.Context) error {
	Consulta, err := codegenerator.GetConsulta(c.PathInt("consulta_id"), s.ddd)
	if err != nil {
		return err
	}

	if c.QueryVal("modo") == "archivo" {
		err = codeGenerator.QryGenerarArchivos(Consulta, c.QueryVal("tipo")).Generar()
		if err != nil {
			return err
		}
		return c.StatusOk("Generado")
	}

	codigo, err := codeGenerator.GenerarDeConsultaStringNew(Consulta, c.QueryVal("tipo"))
	if err != nil {
		return err
	}
	if c.EsHTMX() {
		return c.StatusOk(html.EscapeString(codigo))
	}
	return c.StatusOk(codigo)
}
