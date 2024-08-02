package codegenerator

import (
	"io"
	"monorepo/textutils"
	"strings"

	"github.com/pargomx/gecko/gko"
)

func (tbl *Tabla) renderDeTabla(tipo string, buf io.Writer, separador bool) (err error) {
	op := gko.Op("renderDeTabla").Ctx("tipo", tipo)
	if tbl == nil {
		return op.Msg("tabla es nil")
	}
	op.Ctx("tabla", tbl.Tabla.NombreRepo)
	if tbl.NombreItem() == "" {
		return op.Msg("nombre de modelo indefinido")
	}
	if len(tbl.Campos) == 0 {
		return op.Msg("debe haber al menos una columna")
	}
	if len(tbl.PrimaryKeys()) == 0 {
		return op.Msg("debe haber al menos una clave primaria")
	}
	if tipo == "mysql/list_by" && len(tbl.CamposSeleccionados) == 0 {
		return op.Msg("no se seleccionón ningún campo para list_by")
	}
	if tipo == "mysql/get_by" && len(tbl.CamposSeleccionados) == 0 {
		return op.Msg("no se seleccionón ningún campo para get_by")
	}
	data := map[string]any{
		"TablaOrConsulta": tbl,
		"Tabla":           tbl,
	}
	if tipo == "mysql/servicio" || tipo == "sqlite/servicio" {
		separador = false // nunca porque es header
	}
	if separador {
		textutils.ImprimirSeparador(buf, strings.ToUpper(tipo))
	}

	switch {
	//	case strings.HasPrefix(tipo, "gk"):
	//		data, err = getTablasFK(tbl) // TODO: JOIN
	//		if err != nil {
	//			return err
	//		}
	//		return s.renderer.HaciaBufferGo(tipo, data, buf)

	case strings.HasPrefix(tipo, "html"):
		return tbl.Generador.renderer.HaciaBufferHTML(tipo, data, buf)

	case strings.Contains(tipo, "create_table"):
		//	if err := yamlutils.PopularTablaFKs(tbl); err != nil { // TODO: JOIN
		//		return ctx.Err(err)
		//	}
		return tbl.Generador.renderer.HaciaBuffer(tipo, data, buf)

	default:
		return tbl.Generador.renderer.HaciaBufferGo(tipo, data, buf)
	}
}
