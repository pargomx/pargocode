package codegenerator

import (
	"bytes"
	"fmt"
	"io"
	"monorepo/fileutils"
	"monorepo/textutils"
	"os"
	"path/filepath"
	"strings"

	"github.com/pargomx/gecko/gko"
)

// ================================================================ //
// ========== ARCHIVOS ============================================ //

type tblGenCall struct {
	filename string
	tipo     string
	tbl      *Tabla
	gen      *Generador
	mkdir    bool
}

func (c tblGenCall) Destino() string {
	if c.mkdir {
		return c.filename + " (creará directorio)"
	}
	return c.filename
}

func (c tblGenCall) Generar() error {
	codigo, err := c.gen.GenerarDeTablaString(c.tbl, c.tipo)
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

func (gen *Generador) TblGenerarArchivos(tbl *Tabla, tipo string) tblGenCall {
	c := tblGenCall{
		filename: "generado.go",
		tipo:     tipo,
		tbl:      tbl,
		gen:      gen,
	}
	switch tipo {
	case "entidad":
		c.filename = filepath.Join(tbl.Paquete.Directorio, tbl.Paquete.Nombre, "t_"+tbl.Tabla.Kebab+".go")

	case "mysql", "mysql-directriz":
		c.filename = filepath.Join(tbl.Paquete.Directorio, "mysql"+tbl.Paquete.Nombre, "s_"+tbl.Tabla.NombreRepo+"_gen.go")
		c.tipo = "mysql-directriz"

	case "sqlite", "sqlite-directriz":
		c.filename = filepath.Join(tbl.Paquete.Directorio, "sqlite"+tbl.Paquete.Nombre, "s_"+tbl.Tabla.NombreRepo+"_gen.go")
		c.tipo = "mysql-directriz"

		if !fileutils.Existe(filepath.Join(tbl.Paquete.Directorio, "sqlite"+tbl.Paquete.Nombre, "servicio_repo.go")) {
			err := gen.TblGenerarArchivos(tbl, "sqlite/servicio").Generar()
			if err != nil {
				gko.LogError(err)
			}
		}

	case "mysql-simple", "mysql-compacto", "mysql-all":
		c.filename = filepath.Join(tbl.Paquete.Directorio, "mysql"+tbl.Paquete.Nombre, "s_"+tbl.Tabla.NombreRepo+"_gen.go")

		if !fileutils.Existe(filepath.Join(tbl.Paquete.Directorio, "mysql"+tbl.Paquete.Nombre, "servicio_repo.go")) {
			err := gen.TblGenerarArchivos(tbl, "mysql/servicio").Generar()
			if err != nil {
				gko.LogError(err)
			}
		}

	case "sqlite/servicio":
		c.filename = filepath.Join(tbl.Paquete.Directorio, "sqlite"+tbl.Paquete.Nombre, "servicio_repo.go")

	case "mysql/servicio":
		c.filename = filepath.Join(tbl.Paquete.Directorio, "mysql"+tbl.Paquete.Nombre, "servicio_repo.go")

	}
	c.filename = strings.TrimSuffix(c.filename, "/") // debe ser relativa desde workdir
	if !fileutils.Existe(filepath.Dir(c.filename)) {
		c.mkdir = true
	}
	return c
}

// ================================================================ //

func (gen *Generador) TblGenerarArchivosDirecto(tbl *Tabla, tipo string) error {
	filename := "generado.go"
	switch tipo {

	case "entidad":
		filename = filepath.Join(tbl.Paquete.Directorio, tbl.Paquete.Nombre, tbl.Tabla.Kebab+".go")

	case "mysql-simple", "mysql-compacto", "mysql-all", "mysql-directriz":
		filename = filepath.Join(tbl.Paquete.Directorio, "mysql"+tbl.Paquete.Nombre, tbl.Tabla.Kebab+".go")

	case "mysql/servicio":
		filename = filepath.Join(tbl.Paquete.Directorio, "mysql"+tbl.Paquete.Nombre, "servicio_repo.go")

	case "n":
		return nil

	default:
		return fmt.Errorf("generación por cli no implementada para '%v'", tipo)
	}
	filename = strings.TrimSuffix(filename, "/") // debe ser relativa desde workdir

	// Crear directorio padre si no existe
	if !fileutils.Existe(filepath.Dir(filename)) {
		// TODO: Confirmar
		// if !aveprompt.YN(filepath.Dir(filename)+" (mkdir -p)", true) {
		// 	return nil
		// }
		err := os.MkdirAll(filepath.Dir(filename), 0755)
		if err != nil {
			return err
		}
	}

	// Confirmar generación de código
	// if !aveprompt.YN(filename, true) {
	// 	return nil
	// }
	codigo, err := gen.GenerarDeTablaString(tbl, tipo)
	if err != nil {
		return err
	}
	return fileutils.GuardarGoCode(filename, codigo)
}

// ================================================================ //
// ========== TO STRING =========================================== //

// Helper que admite ejecutar colecciones de plantillas en un mismo string.
func (s *Generador) GenerarDeTablaMySQLDirectriz(tbl *Tabla, buf io.Writer) error {

	if len(tbl.Directrices()) == 0 {
		return fmt.Errorf("no hay directrices para tabla %v", tbl.Tabla.NombreRepo)
	}

	if tbl.Directrices()[0].Key() == "sqlite" {
		fmt.Fprintf(buf, "package sqlite%v\n\n", tbl.Paquete.Nombre)
	} else {
		fmt.Fprintf(buf, "package mysql%v\n\n", tbl.Paquete.Nombre)
	}
	// err := s.GenerarDeTabla(tbl, "mysql/paquete", buf, false)
	// if err != nil {
	// return err
	// }

	err := s.GenerarDeTabla(tbl, "mysql/constantes", buf, true)
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
			err = s.GenerarDeTabla(tbl, "mysql/tbl-insert", buf, true)
		case "update":
			err = s.GenerarDeTabla(tbl, "mysql/tbl-update", buf, true)
		case "insert_update":
			err = s.GenerarDeTabla(tbl, "mysql/tbl-insert_update", buf, true)
		case "delete":
			if !generado.existe {
				if err = s.GenerarDeTabla(tbl, "mysql/existe", buf, true); err != nil {
					return err
				}
				generado.existe = false
			}
			err = s.GenerarDeTabla(tbl, "mysql/tbl-delete", buf, true)

		case "fetch":
			if !generado.scanRow {
				if err = s.GenerarDeTabla(tbl, "mysql/scan-row", buf, true); err != nil {
					return err
				}
				generado.scanRow = true
			}
			err = s.GenerarDeTabla(tbl, "mysql/fetch", buf, true)

		case "get":
			if !generado.scanRow {
				if err = s.GenerarDeTabla(tbl, "mysql/scan-row", buf, true); err != nil {
					return err
				}
				generado.scanRow = true
			}
			err = s.GenerarDeTabla(tbl, "mysql/get", buf, true)

		case "get_by":
			if !generado.scanRow {
				if err = s.GenerarDeTabla(tbl, "mysql/scan-row", buf, true); err != nil {
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
			err = s.GenerarDeTabla(tbl, "mysql/get_by", buf, true)

		case "list":
			if generado.scanRows {
				if err = s.GenerarDeTabla(tbl, "mysql/scan-rows", buf, true); err != nil {
					return err
				}
				generado.scanRows = true
			}
			if filtros {
				if err = s.GenerarDeTabla(tbl, "mysql/tbl-filtros", buf, true); err != nil {
					return err
				}
				filtros = false
			}
			err = s.GenerarDeTabla(tbl, "mysql/list", buf, true)

		case "list_by":
			if generado.scanRows {
				if err = s.GenerarDeTabla(tbl, "mysql/scan-rows", buf, true); err != nil {
					return err
				}
				generado.scanRows = true
			}
			if filtros {
				if err = s.GenerarDeTabla(tbl, "mysql/tbl-filtros", buf, true); err != nil {
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
			err = s.GenerarDeTabla(tbl, "mysql/list_by", buf, true)

		default:
			return fmt.Errorf("directriz '%v' no aplica para mysql de tabla", directriz.Key())
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Helper que admite ejecutar colecciones de plantillas en un mismo string.
func (s *Generador) GenerarDeTablaString(tbl *Tabla, tipo string) (string, error) {
	buf := &bytes.Buffer{}
	var err error
	switch tipo {
	case "entidad":
		err = s.GenerarDeTablaEntidad(tbl, buf)
	case "mysql-all":
		err = s.GenerarDeTablaAllMySQL(tbl, buf)
	case "mysql-simple":
		err = s.GenerarDeTablaMySQLSimple(tbl, buf)
	case "mysql-compacto":
		err = s.GenerarDeTablaMySQLInsertUpdate(tbl, buf)
	case "mysql-directriz":
		err = s.GenerarDeTablaMySQLDirectriz(tbl, buf)
	case "sqlite":
		err = s.GenerarDeTablaMySQLDirectriz(tbl, buf)
	default:
		err = s.GenerarDeTabla(tbl, tipo, buf, false)
	}
	return buf.String(), err
}

// ================================================================ //
// ========== COLECCIONES ========================================= //

func (s *Generador) GenerarDeTablaEntidad(tbl *Tabla, buf io.Writer) error {
	ctx := gko.Op("GenerarDeTablaEntidad")
	if err := s.GenerarDeTabla(tbl, "go/tbl_struct", buf, false); err != nil {
		return ctx.Err(err)
	}
	if err := s.GenerarDeTabla(tbl, "go/tbl_errores", buf, false); err != nil {
		return ctx.Err(err)
	}
	if err := s.GenerarDeTabla(tbl, "go/tbl_propiedades", buf, false); err != nil {
		return ctx.Err(err)
	}
	return nil
}

func (s *Generador) GenerarDeTablaAllMySQL(tbl *Tabla, buf io.Writer) error {
	ctx := gko.Op("GenerarDeTablaAllMySQL")
	err := s.GenerarDeTabla(tbl, "mysql/paquete", buf, false)
	if err != nil {
		return err
	}
	porGenerar := []string{
		// "mysql/servicio",
		"mysql/constantes",
		"mysql/scan-row",
		"mysql/get",
		"mysql/fetch",
		"mysql/tbl-insert",
		"mysql/tbl-update",
		"mysql/tbl-insert_update",
		"mysql/tbl-delete",
		"mysql/scan-rows",
		"mysql/tbl-filtros",
		"mysql/list",
	}
	for _, tipo := range porGenerar {
		if err := s.GenerarDeTabla(tbl, tipo, buf, true); err != nil {
			return ctx.Err(err)
		}
	}
	for _, campo := range tbl.UniqueKeys() {
		tbl.CamposSeleccionados = []CampoTabla{campo}
		err := s.GenerarDeTabla(tbl, "mysql/get_by", buf, true)
		if err != nil {
			return ctx.Ctx("mysql/get_by", campo.Campo.NombreCampo).Err(err)
		}
	}
	for _, fk := range tbl.ForeignKeys() {
		tbl.CamposSeleccionados = []CampoTabla{fk}
		err := s.GenerarDeTabla(tbl, "mysql/list_by", buf, true)
		if err != nil {
			return ctx.Ctx("mysql/list_by", fk.Campo.NombreCampo).Err(err)
		}
	}
	return nil
}

func (s *Generador) GenerarDeTablaMySQLInsertUpdate(tbl *Tabla, buf io.Writer) error {
	ctx := gko.Op("GenerarDeTablaMySQLInsertUpdate")
	err := s.GenerarDeTabla(tbl, "mysql/paquete", buf, false)
	if err != nil {
		return err
	}
	porGenerar := []string{
		"mysql/constantes",
		"mysql/scan-row",
		"mysql/get",
		"mysql/fetch",
		"mysql/tbl-insert_update",
		"mysql/tbl-delete",
	}
	for _, tipo := range porGenerar {
		if err := s.GenerarDeTabla(tbl, tipo, buf, true); err != nil {
			return ctx.Err(err)
		}
	}
	return nil
}

func (s *Generador) GenerarDeTablaMySQLSimple(tbl *Tabla, buf io.Writer) error {
	ctx := gko.Op("GenerarDeTablaMySQLSimple")
	err := s.GenerarDeTabla(tbl, "mysql/paquete", buf, false)
	if err != nil {
		return err
	}
	porGenerar := []string{
		"mysql/constantes",
		"mysql/scan-row",
		"mysql/get",
		"mysql/tbl-insert",
		"mysql/tbl-update",
		"mysql/tbl-delete",
	}
	for _, tipo := range porGenerar {
		if err := s.GenerarDeTabla(tbl, tipo, buf, true); err != nil {
			return ctx.Err(err)
		}
	}
	return nil
}

// ================================================================ //
// ========== GENERAR ============================================= //

func (s *Generador) GenerarDeTabla(tbl *Tabla, tipo string, buf io.Writer, separador bool) (err error) {
	ctx := gko.Op("generar").Ctx("tipo", tipo)
	if tbl == nil {
		return ctx.Msg("tabla es nil")
	}
	ctx.Ctx("tabla", tbl.Tabla.NombreRepo)
	if tbl.NombreItem() == "" {
		return ctx.Msg("nombre de modelo indefinido")
	}
	if len(tbl.Campos) == 0 {
		return ctx.Msg("debe haber al menos una columna")
	}
	if len(tbl.PrimaryKeys()) == 0 {
		return ctx.Msg("debe haber al menos una clave primaria")
	}
	if tipo == "mysql/list_by" && len(tbl.CamposSeleccionados) == 0 {
		return ctx.Msg("no se seleccionón ningún campo para list_by")
	}
	if tipo == "mysql/get_by" && len(tbl.CamposSeleccionados) == 0 {
		return ctx.Msg("no se seleccionón ningún campo para get_by")
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
		return s.renderer.HaciaBufferHTML(tipo, data, buf)

	case strings.Contains(tipo, "create_table"):
		//	if err := yamlutils.PopularTablaFKs(tbl); err != nil { // TODO: JOIN
		//		return ctx.Err(err)
		//	}
		return s.renderer.HaciaBuffer(tipo, data, buf)

	default:
		return s.renderer.HaciaBufferGo(tipo, data, buf)
	}
}
