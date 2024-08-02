package codegenerator

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pargomx/gecko/gko"
)

func (tbl *Tabla) GenerarDeTablaString(tipo string) (string, error) {
	buf := &bytes.Buffer{}
	var err error
	switch tipo {
	case "entidad":
		err = tbl.renderTablaEntidad(buf)
	case "mysql-directriz", "sqlite":
		err = tbl.renderTablaRepoSQL(buf)
	default:
		err = tbl.renderDeTabla(tipo, buf, false)
	}
	return buf.String(), err
}

func (tbl *Tabla) renderTablaEntidad(buf io.Writer) error {
	ctx := gko.Op("GenerarDeTablaEntidad")
	if err := tbl.renderDeTabla("go/tbl_struct", buf, false); err != nil {
		return ctx.Err(err)
	}
	if err := tbl.renderDeTabla("go/tbl_errores", buf, false); err != nil {
		return ctx.Err(err)
	}
	if err := tbl.renderDeTabla("go/tbl_propiedades", buf, false); err != nil {
		return ctx.Err(err)
	}
	return nil
}

func (tbl *Tabla) renderTablaRepoSQL(buf io.Writer) error {

	if len(tbl.Directrices()) == 0 {
		return fmt.Errorf("no hay directrices para tabla %v", tbl.Tabla.NombreRepo)
	}

	if tbl.Directrices()[0].Key() == "sqlite" {
		fmt.Fprintf(buf, "package sqlite%v\n\n", tbl.Paquete.Nombre)
	} else {
		fmt.Fprintf(buf, "package mysql%v\n\n", tbl.Paquete.Nombre)
	}

	err := tbl.renderDeTabla("mysql/constantes", buf, true)
	if err != nil {
		return err
	}

	var generado struct {
		scanRow  bool
		scanRows bool
		existe   bool
	}

	filtros := tbl.TieneCamposFiltro()

	for _, directriz := range tbl.Directrices() {
		tbl.CamposSeleccionados = nil
		switch directriz.Key() {

		case "sqlite":
			// Solo cambia el nombre del paquete
		case "insert":
			err = tbl.renderDeTabla("mysql/tbl-insert", buf, true)
		case "update":
			err = tbl.renderDeTabla("mysql/tbl-update", buf, true)
		case "insert_update":
			err = tbl.renderDeTabla("mysql/tbl-insert_update", buf, true)
		case "delete":
			if !generado.existe {
				if err = tbl.renderDeTabla("mysql/existe", buf, true); err != nil {
					return err
				}
				generado.existe = false
			}
			err = tbl.renderDeTabla("mysql/tbl-delete", buf, true)

		case "fetch":
			if !generado.scanRow {
				if err = tbl.renderDeTabla("mysql/scan-row", buf, true); err != nil {
					return err
				}
				generado.scanRow = true
			}
			err = tbl.renderDeTabla("mysql/fetch", buf, true)

		case "get":
			if !generado.scanRow {
				if err = tbl.renderDeTabla("mysql/scan-row", buf, true); err != nil {
					return err
				}
				generado.scanRow = true
			}
			err = tbl.renderDeTabla("mysql/get", buf, true)

		case "get_by":
			if !generado.scanRow {
				if err = tbl.renderDeTabla("mysql/scan-row", buf, true); err != nil {
					return err
				}
				generado.scanRow = true
			}
			for _, v := range directriz.Values() {
				campo, err := tbl.BuscarCampo(v)
				if err != nil {
					return err
				}
				tbl.CamposSeleccionados = append(tbl.CamposSeleccionados, *campo)
			}
			err = tbl.renderDeTabla("mysql/get_by", buf, true)

		case "list":
			if generado.scanRows {
				if err = tbl.renderDeTabla("mysql/scan-rows", buf, true); err != nil {
					return err
				}
				generado.scanRows = true
			}
			if filtros {
				if err = tbl.renderDeTabla("mysql/tbl-filtros", buf, true); err != nil {
					return err
				}
				filtros = false
			}
			err = tbl.renderDeTabla("mysql/list", buf, true)

		case "list_by":
			if generado.scanRows {
				if err = tbl.renderDeTabla("mysql/scan-rows", buf, true); err != nil {
					return err
				}
				generado.scanRows = true
			}
			if filtros {
				if err = tbl.renderDeTabla("mysql/tbl-filtros", buf, true); err != nil {
					return err
				}
				filtros = false
			}
			for _, v := range directriz.Values() {
				campo, err := tbl.BuscarCampo(v)
				if err != nil {
					return err
				}
				tbl.CamposSeleccionados = append(tbl.CamposSeleccionados, *campo)
			}
			err = tbl.renderDeTabla("mysql/list_by", buf, true)

		default:
			return fmt.Errorf("directriz '%v' no aplica para mysql de tabla", directriz.Key())
		}
		if err != nil {
			return err
		}
	}
	return nil
}
