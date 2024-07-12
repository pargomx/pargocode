package gencode

import (
	"bytes"
	"fmt"
	"io"
	"monorepo/domain_driven_design/dpaquete"
	"monorepo/fileutils"
	"monorepo/textutils"
	"os"
	"path/filepath"
	"strings"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== ARCHIVOS ============================================ //

type qryGenCall struct {
	filename string
	tipo     string
	qry      *dpaquete.Consulta
	gen      *Generador
	mkdir    bool
}

func (c qryGenCall) Destino() string {
	if c.mkdir {
		return c.filename + " (creará directorio)"
	}
	return c.filename
}

func (c qryGenCall) Generar() error {
	codigo, err := c.gen.GenerarDeConsultaString(c.qry, c.tipo)
	if err != nil {
		return err
	}
	if c.mkdir { // Crear directorio padre si no existe
		err := os.MkdirAll(filepath.Dir(c.filename), 0755)
		if err != nil {
			return err
		}
	}
	return fileutils.GuardarGoCode(c.filename, codigo)
}

func (gen *Generador) QryGenerarArchivos(qry *dpaquete.Consulta, tipo string) qryGenCall {
	c := qryGenCall{
		filename: "generado.go",
		tipo:     tipo,
		qry:      qry,
		gen:      gen,
	}
	switch tipo {
	case "entidad":
		c.filename = filepath.Join(qry.Paquete.Directorio, qry.Paquete.Nombre,
			"q_"+qry.Consulta.NombreItem+".go")

	case "mysql":
		c.filename = filepath.Join(qry.Paquete.Directorio, "mysql"+qry.Paquete.Nombre,
			"q_"+qry.Consulta.NombreItem+".go")
		c.tipo = "mysql-directriz"

	case "sqlite", "sqlite-directriz":
		c.filename = filepath.Join(qry.Paquete.Directorio, "sqlite"+qry.Paquete.Nombre, "s_"+qry.Consulta.NombreItem+"_gen.go")
		c.tipo = "mysql-directriz"

	}
	c.filename = strings.TrimSuffix(c.filename, "/") // debe ser relativa desde workdir
	if !fileutils.Existe(filepath.Dir(c.filename)) {
		c.mkdir = true
	}
	return c
}

// ================================================================ //
// ========== DIRECTRIZ =========================================== //

// Helper que admite ejecutar colecciones de plantillas en un mismo string.
func (s *Generador) GenerarDeConsultaMySQLDirectriz(qry *dpaquete.Consulta, buf io.Writer) error {

	if len(qry.Directrices()) == 0 {
		return fmt.Errorf("no hay directrices para %v", qry.Consulta.NombreItem)
	}

	if qry.Directrices()[0].Key() == "sqlite" {
		fmt.Fprintf(buf, "package sqlite%v\n\n", qry.Paquete.Nombre)
	} else {
		fmt.Fprintf(buf, "package mysql%v\n\n", qry.Paquete.Nombre)
	}

	// err := s.GenerarDeConsultaNew(qry, "mysql/paquete", buf, false)
	// if err != nil {
	// 	return err
	// }

	err := s.GenerarDeConsultaNew(qry, "mysql/constantes", buf, true)
	if err != nil {
		return err
	}

	scanRow := true
	scanRows := true
	filtros := qry.TieneCamposFiltro()

	for _, directriz := range qry.Directrices() {
		qry.CamposSeleccionados = nil
		switch directriz.Key() {

		case "sqlite":
			// Solo cambia el nombre del paquete

		case "fetch":
			if scanRow {
				if err = s.GenerarDeConsultaNew(qry, "mysql/scan-row", buf, true); err != nil {
					return err
				}
				scanRow = false
			}
			err = s.GenerarDeConsultaNew(qry, "mysql/fetch", buf, true)

		case "get":
			if scanRow {
				if err = s.GenerarDeConsultaNew(qry, "mysql/scan-row", buf, true); err != nil {
					return err
				}
				scanRow = false
			}
			err = s.GenerarDeConsultaNew(qry, "mysql/get", buf, true)

		case "get_by":
			if scanRow {
				if err = s.GenerarDeConsultaNew(qry, "mysql/scan-row", buf, true); err != nil {
					return err
				}
				scanRow = false
			}
			for _, v := range directriz.Values() {
				campo, err := qry.BuscarCampo(v)
				if err != nil {
					return err
				}
				qry.CamposSeleccionados = append(qry.CamposSeleccionados, *campo)
			}
			err = s.GenerarDeConsultaNew(qry, "mysql/get_by", buf, true)

		case "list":
			if scanRows {
				if err = s.GenerarDeConsultaNew(qry, "mysql/scan-rows", buf, true); err != nil {
					return err
				}
				scanRows = false
			}
			if filtros {
				if err = s.GenerarDeConsultaNew(qry, "mysql/qry-filtros", buf, true); err != nil {
					return err
				}
				filtros = false
			}
			err = s.GenerarDeConsultaNew(qry, "mysql/list", buf, true)

		case "list_by":
			if scanRows {
				if err = s.GenerarDeConsultaNew(qry, "mysql/scan-rows", buf, true); err != nil {
					return err
				}
				scanRows = false
			}
			if filtros {
				if err = s.GenerarDeConsultaNew(qry, "mysql/qry-filtros", buf, true); err != nil {
					return err
				}
				filtros = false
			}
			for _, v := range directriz.Values() {
				campo, err := qry.BuscarCampo(v)
				if err != nil {
					return err
				}
				qry.CamposSeleccionados = append(qry.CamposSeleccionados, *campo)
			}
			err = s.GenerarDeConsultaNew(qry, "mysql/list_by", buf, true)

		default:
			return fmt.Errorf("directriz '%v' no aplica para mysql de consulta", directriz.Key())
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// ================================================================ //
// ========== TO STRING =========================================== //

func (s *Generador) GenerarDeConsultaString(consulta *dpaquete.Consulta, tipo string) (string, error) {
	buf := &bytes.Buffer{}
	var err error
	switch tipo {
	case "mysql":
		err = s.GenerarDeConsultaAllMySQL(consulta, buf)

	case "entidad":
		tipo = "go/qry_struct"
		err = s.GenerarDeConsultaNew(consulta, tipo, buf, false)

	case "mysql-directriz":
		err = s.GenerarDeConsultaMySQLDirectriz(consulta, buf)

	default:
		err = s.GenerarDeConsultaNew(consulta, tipo, buf, false)
	}
	return buf.String(), err
}

func (s *Generador) GenerarDeConsultaStringNew(consulta *dpaquete.Consulta, tipo string) (string, error) {
	buf := &bytes.Buffer{}
	var err error
	switch tipo {

	case "entidad":
		tipo = "go/qry_struct"
		err = s.GenerarDeConsultaNew(consulta, tipo, buf, false)

	case "mysql", "mysql-directriz":
		tipo = "mysql-directriz"
		err = s.GenerarDeConsultaMySQLDirectriz(consulta, buf)

	default:
		err = s.GenerarDeConsultaNew(consulta, tipo, buf, false)
	}
	return buf.String(), err
}

// ================================================================ //
// ========== COLECCIONES ========================================= //

func (s *Generador) GenerarDeConsultaAllMySQL(qry *dpaquete.Consulta, buf io.Writer) error {
	ctx := gecko.NewErr(800).Op("GenerarDeConsultaAllMySQL")
	err := s.GenerarDeConsultaNew(qry, "mysql/paquete", buf, false)
	if err != nil {
		return err
	}
	porGenerar := []string{
		"mysql/constantes",
		"mysql/scan-row",
		"mysql/get",
		"mysql/fetch",

		"mysql/qry-filtros",
		"mysql/scan-rows",
		"mysql/list",
	}
	for _, tipo := range porGenerar {
		if err := s.GenerarDeConsultaNew(qry, tipo, buf, true); err != nil {
			return ctx.Err(err)
		}
	}
	return nil
}

// ================================================================ //
// ========== GENERAR ============================================= //

func (s *Generador) GenerarDeConsultaNew(consulta *dpaquete.Consulta, tipo string, buf io.Writer, separador bool) error {
	ctx := gecko.NewErr(700).Op("generar").Ctx("tipo", tipo)
	if consulta == nil {
		return ctx.Msg("consulta es nil")
	}
	ctx.Ctx("consulta", consulta.Consulta.NombreItem)
	if consulta.Consulta.NombreItem == "" {
		return ctx.Msg("nombre de modelo indefinido")
	}
	if len(consulta.Campos) == 0 {
		return ctx.Msg("debe haber al menos una columna")
	}
	// if len(consulta.PrimaryKeys) == 0 {
	// 	return ctx.Msg("debe haber al menos una clave primaria")
	// }
	// if (tipo == "mysql/list_by" || tipo == "mysql/get_by") && len(consulta.Campos) == 0 {
	// 	return ctx.Msg("no se seleccionón ningún campo para list_by o get_by")
	// }
	data := map[string]any{
		"TablaOrConsulta":  consulta,
		"Consulta":         consulta,
		"AgregadoConsulta": consulta,
	}
	if separador {
		textutils.ImprimirSeparador(buf, strings.ToUpper(tipo))
	}
	switch {
	case tipo == "mysql/query":
		return s.renderer.HaciaBuffer(tipo, data, buf)
	default:
		return s.renderer.HaciaBufferGo(tipo, data, buf)
	}
}
