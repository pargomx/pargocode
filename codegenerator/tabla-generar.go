package codegenerator

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pargomx/gecko/gko"
)

func (c tblGenCall) GenerarToString(tipo string) (string, error) {
	buf := &bytes.Buffer{}
	var err error
	switch tipo {
	case "entidad":
		err = c.renderTablaEntidad(buf)
	case "mysql", "sqlite":
		err = c.renderTablaRepoSQL(buf)
	default:
		err = c.renderToBuffer(tipo, buf, false)
	}
	return buf.String(), err
}

// ================================================================ //
// ================================================================ //

func (c tblGenCall) renderTablaEntidad(buf io.Writer) error {
	ctx := gko.Op("GenerarDeTablaEntidad")
	if err := c.renderToBuffer("go/tbl_struct", buf, false); err != nil {
		return ctx.Err(err)
	}
	if err := c.renderToBuffer("go/tbl_errores", buf, false); err != nil {
		return ctx.Err(err)
	}
	if err := c.renderToBuffer("go/tbl_propiedades", buf, false); err != nil {
		return ctx.Err(err)
	}
	return nil
}

func (c tblGenCall) renderTablaRepoSQL(buf io.Writer) error {

	if len(c.tbl.Directrices()) == 0 {
		return fmt.Errorf("no hay directrices para tabla %v", c.tbl.Tabla.NombreRepo)
	}

	if c.tbl.Directrices()[0].Key() == "sqlite" {
		fmt.Fprintf(buf, "package sqlite%v\n\n", c.tbl.Paquete.Nombre)
	} else {
		fmt.Fprintf(buf, "package mysql%v\n\n", c.tbl.Paquete.Nombre)
	}

	err := c.renderToBuffer("mysql/constantes", buf, true)
	if err != nil {
		return err
	}

	var generado struct {
		scanRow  bool
		scanRows bool
		existe   bool
	}

	filtros := c.tbl.TieneCamposFiltro()

	for _, directriz := range c.tbl.Directrices() {
		c.tbl.CamposSeleccionados = nil
		switch directriz.Key() {

		case "sqlite":
			// Solo cambia el nombre del paquete
		case "insert":
			err = c.renderToBuffer("mysql/tbl-insert", buf, true)
		case "update":
			err = c.renderToBuffer("mysql/tbl-update", buf, true)
		case "insert_update":
			err = c.renderToBuffer("mysql/tbl-insert_update", buf, true)
		case "delete":
			if !generado.existe {
				if err = c.renderToBuffer("mysql/existe", buf, true); err != nil {
					return err
				}
				generado.existe = false
			}
			err = c.renderToBuffer("mysql/tbl-delete", buf, true)

		case "fetch":
			if !generado.scanRow {
				if err = c.renderToBuffer("mysql/scan-row", buf, true); err != nil {
					return err
				}
				generado.scanRow = true
			}
			err = c.renderToBuffer("mysql/fetch", buf, true)

		case "get":
			if !generado.scanRow {
				if err = c.renderToBuffer("mysql/scan-row", buf, true); err != nil {
					return err
				}
				generado.scanRow = true
			}
			err = c.renderToBuffer("mysql/get", buf, true)

		case "get_by":
			if !generado.scanRow {
				if err = c.renderToBuffer("mysql/scan-row", buf, true); err != nil {
					return err
				}
				generado.scanRow = true
			}
			for _, v := range directriz.Values() {
				campo, err := c.tbl.BuscarCampo(v)
				if err != nil {
					return err
				}
				c.tbl.CamposSeleccionados = append(c.tbl.CamposSeleccionados, *campo)
			}
			err = c.renderToBuffer("mysql/get_by", buf, true)

		case "list":
			if generado.scanRows {
				if err = c.renderToBuffer("mysql/scan-rows", buf, true); err != nil {
					return err
				}
				generado.scanRows = true
			}
			if filtros {
				if err = c.renderToBuffer("mysql/tbl-filtros", buf, true); err != nil {
					return err
				}
				filtros = false
			}
			err = c.renderToBuffer("mysql/list", buf, true)

		case "list_by":
			if generado.scanRows {
				if err = c.renderToBuffer("mysql/scan-rows", buf, true); err != nil {
					return err
				}
				generado.scanRows = true
			}
			if filtros {
				if err = c.renderToBuffer("mysql/tbl-filtros", buf, true); err != nil {
					return err
				}
				filtros = false
			}
			for _, v := range directriz.Values() {
				campo, err := c.tbl.BuscarCampo(v)
				if err != nil {
					return err
				}
				c.tbl.CamposSeleccionados = append(c.tbl.CamposSeleccionados, *campo)
			}
			err = c.renderToBuffer("mysql/list_by", buf, true)

		default:
			return fmt.Errorf("directriz '%v' no aplica para mysql de tabla", directriz.Key())
		}
		if err != nil {
			return err
		}
	}
	return nil
}
