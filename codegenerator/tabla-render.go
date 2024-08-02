package codegenerator

import (
	"io"
	"monorepo/textutils"
	"strings"

	"github.com/pargomx/gecko/gko"
)

func (c tblGenCall) renderToBuffer(tmpl string, buf io.Writer, separador bool) (err error) {
	op := gko.Op("tbl.renderToBuffer").Ctx("tipo", c.tipo)
	if c.tbl == nil {
		return op.Msg("tabla es nil")
	}
	op.Ctx("tabla", c.tbl.Tabla.NombreRepo)
	if c.tbl.NombreItem() == "" {
		return op.Msg("nombre de modelo indefinido")
	}
	if len(c.tbl.Campos) == 0 {
		return op.Msg("debe haber al menos una columna")
	}
	if len(c.tbl.PrimaryKeys()) == 0 {
		return op.Msg("debe haber al menos una clave primaria")
	}
	if tmpl == "mysql/list_by" && len(c.tbl.CamposSeleccionados) == 0 {
		return op.Msg("no se seleccionón ningún campo para list_by")
	}
	if tmpl == "mysql/get_by" && len(c.tbl.CamposSeleccionados) == 0 {
		return op.Msg("no se seleccionón ningún campo para get_by")
	}
	data := map[string]any{
		"TablaOrConsulta": c.tbl,
		"Tabla":           c.tbl,
	}
	if tmpl == "mysql/servicio" || tmpl == "sqlite/servicio" {
		separador = false // nunca porque es header
	}
	if separador {
		textutils.ImprimirSeparador(buf, strings.ToUpper(tmpl))
	}

	switch {
	//	case strings.HasPrefix(tipo, "gk"):
	//		data, err = getTablasFK(tbl) // TODO: JOIN
	//		if err != nil {
	//			return err
	//		}
	//		return s.renderer.HaciaBufferGo(tipo, data, buf)

	case strings.HasPrefix(tmpl, "html"):
		return c.tbl.Generador.renderer.HaciaBufferHTML(tmpl, data, buf)

	case strings.Contains(tmpl, "create_table"):
		//	if err := yamlutils.PopularTablaFKs(tbl); err != nil { // TODO: JOIN
		//		return ctx.Err(err)
		//	}
		return c.tbl.Generador.renderer.HaciaBuffer(tmpl, data, buf)

	default:
		return c.tbl.Generador.renderer.HaciaBufferGo(tmpl, data, buf)
	}
}
