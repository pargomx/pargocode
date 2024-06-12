package main

import (
	"html"
	"monorepo/domain_driven_design/dpaquete"
	"monorepo/gecko"
)

// ================================================================ //
// ========== GENERAR ============================================= //

func (s *servidor) generarDeConsulta(c *gecko.Context) error {
	agregadoConsulta, err := dpaquete.GetAgregadoConsulta(c.PathInt("consulta_id"), s.ddd)
	if err != nil {
		return err
	}

	if c.QueryVal("modo") == "archivo" {
		err = codeGenerator.QryGenerarArchivos(agregadoConsulta, c.QueryVal("tipo")).Generar()
		if err != nil {
			return err
		}
		return c.StatusOk("Generado")
	}

	codigo, err := codeGenerator.GenerarDeConsultaStringNew(agregadoConsulta, c.QueryVal("tipo"))
	if err != nil {
		return err
	}
	if c.EsHTMX() {
		return c.StatusOk(html.EscapeString(codigo))
	}
	return c.StatusOk(codigo)
}
