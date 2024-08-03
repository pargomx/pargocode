package codegenerator

import (
	"bytes"
	"fmt"
	"monorepo/ddd"
	"monorepo/fileutils"
	"monorepo/textutils"
	"os"
	"path/filepath"
	"strings"

	"github.com/pargomx/gecko/gko"
)

// ================================================================ //
// ========== GENERADOR DE TABLA ================================== //

type tblGenerator struct {
	gen *Generador
	tbl *tabla

	destinos []genDest
	hechos   []string
	errores  []error
}

type tabla struct {
	Tabla   ddd.Tabla
	Campos  []CampoTabla
	Paquete ddd.Paquete

	CamposSeleccionados []CampoTabla

	Sqlite bool
	MySQL  bool
}

type genDest struct {
	jobs     []genJob      // trabajos por realizar
	buf      *bytes.Buffer // buffer para escribir el código
	filename string        // nombre del archivo destino (si aplica)
	mkdir    bool          // crear directorio si no existe
}

type genJob struct {
	tmpl string // nombre de la plantilla a renderizar
	// separador bool         // agregar comentario separador
	titulo string       // título del bloque para separador
	campos []CampoTabla // campos seleccionados
}

// ================================================================ //

func (c *tblGenerator) addErr(err error) {
	if err != nil {
		c.errores = append(c.errores, err)
	}
}
func (c *tblGenerator) Errs() string {
	txt := ""
	for _, err := range c.errores {
		txt += err.Error() + "\n"
	}
	txt = strings.TrimSuffix(txt, "\n")
	return txt
}

func (c *tblGenerator) addHecho(hecho string) {
	c.hechos = append(c.hechos, hecho)
}
func (c *tblGenerator) GetHechos() []string {
	return c.hechos
}

// ================================================================ //
// ========== NUEVO GENERADOR PARA TABLA ========================== //

// Nuevo generador de código para una tabla.
func (gen *Generador) DeTabla(tablaID int) (*tblGenerator, error) {
	tbl, err := gen.getTabla(tablaID)
	if err != nil {
		return nil, err
	}
	tblGenCall := tblGenerator{
		gen: gen,
		tbl: tbl,
	}
	return &tblGenCall, nil
}

// ================================================================ //
// ========== DEFINIR TRABAJO ===================================== //

// Definir qué se va a generar. Dependiendo de lo solicitado se mandará guardar
// al lugar adecuado. Pueden ser diferentes archivos o buffers de destino.
func (c *tblGenerator) PrepararJob(tipo string) *tblGenerator {
	if tipo == "" {
		c.addErr(gko.ErrDatoIndef().Str("tipo de trabajo indefinido").Op("PrepararJob"))
		return c
	}
	switch tipo {
	case "entidad":
		c.addJobsEntidad()
	case "mysql":
		c.addJobsRepoSQL(false)
	case "sqlite":
		c.addJobsRepoSQL(true)
	default:
		c.addJob(tipo, "generado.go", "", nil)
	}
	return c
}

// Imprime en consola los trabajos del generador.
func (c *tblGenerator) DescribirJobs() {
	porHacer := ""
	for _, dest := range c.destinos {
		porHacer += fmt.Sprintf("\033[34m%v\033[0m\n", dest.filename)
		for _, job := range dest.jobs {
			porHacer += fmt.Sprintf("  \033[33m%v\033[0m\n", job.tmpl)
		}
	}
	gko.LogInfo("Generando código:")
	fmt.Println(porHacer)
}

// ================================================================ //
// ========== PREPARAR TRABAJO PENDIENTE ========================== //

// Agregar un trabjo a la lista de pendientes con el destino adecuado.
// El título es opcional y se usa para agregar un comentario separador.
func (c *tblGenerator) addJob(tmpl string, destino string, titulo string, campos []string) {
	op := gko.Op("addJob").Ctx("tmpl", tmpl)
	if c.tbl == nil {
		c.addErr(op.Str("tabla no definida"))
		return
	}
	if tmpl == "" {
		c.addErr(op.Str("plantilla no definida"))
		return
	}
	// no agregar trabajo si ya se declaró
	for _, dest := range c.destinos {
		if dest.filename == destino {
			for _, job := range dest.jobs {
				if job.tmpl == tmpl {
					return
				}
			}
		}
	}
	// agregar trabajo al destino adecuado
	destIdx, err := c.addDestino(destino)
	if err != nil {
		c.addErr(op.Err(err))
		return
	}
	job := genJob{
		tmpl:   tmpl,
		titulo: titulo,
	}
	// agregar campos seleccionados si aplica
	for _, v := range campos {
		campo, err := c.tbl.BuscarCampo(v)
		if err != nil {
			c.addErr(op.Err(err))
			continue
		}
		job.campos = append(job.campos, *campo)
	}
	c.destinos[destIdx].jobs = append(c.destinos[destIdx].jobs, job)
}

// ================================================================ //

// Devuelve el idx del addDestino declarado (lo agrega si no existe).
func (c *tblGenerator) addDestino(filename string) (int, error) {
	filename = strings.TrimSuffix(filename, "/") // debe ser un archivo
	filename = strings.TrimPrefix(filename, "/") // debe ser relativa desde workdir
	if filename == "" {
		return 0, gko.ErrDatoIndef().Str("filename es vacío").Op("nuevoDestino")
	}
	for i, dest := range c.destinos {
		if dest.filename == filename {
			return i, nil
		}
	}
	dest := genDest{
		filename: filename,
		buf:      new(bytes.Buffer),
	}
	if !fileutils.Existe(filepath.Dir(filename)) {
		dest.mkdir = true
	}
	c.destinos = append(c.destinos, dest)
	if len(c.destinos) == 0 {
		return 0, gko.ErrInesperado().Str("no se agregó destino").Op("nuevoDestino")
	}
	return len(c.destinos) - 1, nil
}

// ================================================================ //
// ========== COLECCIONES DE TRABAJO ============================== //

func (c *tblGenerator) addJobsEntidad() {
	destino := filepath.Join(c.tbl.Paquete.Directorio,
		c.tbl.Paquete.Nombre, "t_"+c.tbl.Tabla.Kebab+".go")
	c.addJob("go/tbl_struct", destino, "", nil)
	c.addJob("go/tbl_errores", destino, "", nil)
	c.addJob("go/tbl_propiedades", destino, "", nil)
}

// ================================================================ //

func (c *tblGenerator) addJobsRepoSQL(sqlite bool) {
	op := gko.Op("addJobsRepoSQL").Ctx("tabla", c.tbl.Tabla.NombreRepo)
	if len(c.tbl.Directrices()) == 0 {
		c.addErr(op.Str("no hay directrices"))
	}

	c.tbl.Sqlite = sqlite
	c.tbl.MySQL = !sqlite

	// Si el paquete aún no tiene el archivo de servicio, se agrega.
	if sqlite && !fileutils.Existe(filepath.Join(c.tbl.Paquete.Directorio, "sqlite"+c.tbl.Paquete.Nombre, "servicio_repo.go")) {
		c.addJob("sqlite/servicio", filepath.Join(c.tbl.Paquete.Directorio, "sqlite"+c.tbl.Paquete.Nombre, "servicio_repo.go"), "", nil)

	} else if !sqlite && !fileutils.Existe(filepath.Join(c.tbl.Paquete.Directorio, "mysql"+c.tbl.Paquete.Nombre, "servicio_repo.go")) {
		c.addJob("mysql/servicio", filepath.Join(c.tbl.Paquete.Directorio, "mysql"+c.tbl.Paquete.Nombre, "servicio_repo.go"), "", nil)
	}

	// Destino diferente dependiendo si es repo mysql o sqlite.
	destino := "generado.go"
	if sqlite {
		destino = filepath.Join(c.tbl.Paquete.Directorio,
			"sqlite"+c.tbl.Paquete.Nombre, "s_"+c.tbl.Tabla.NombreRepo+"_gen.go")
	} else {
		destino = filepath.Join(c.tbl.Paquete.Directorio,
			"mysql"+c.tbl.Paquete.Nombre, "s_"+c.tbl.Tabla.NombreRepo+"_gen.go")
	}

	c.addJob("mysql/paquete", destino, "", nil)
	for _, directriz := range c.tbl.Directrices() {
		switch directriz.Key() {
		case "insert":
			c.addJob("mysql/tbl-insert", destino, directriz.Key(), nil)

		case "update":
			c.addJob("mysql/tbl-update", destino, directriz.Key(), nil)

		case "insert_update":
			c.addJob("mysql/tbl-insert_update", destino, directriz.Key(), nil)

		case "delete":
			c.addJob("mysql/existe", destino, "EXISTE", nil)
			c.addJob("mysql/tbl-delete", destino, "DELETE", nil)

		case "fetch":
			c.addJob("mysql/constantes", destino, "CONSTANTES", nil)
			c.addJob("mysql/scan-row", destino, "SCAN", nil)
			c.addJob("mysql/fetch", destino, "FETCH", nil)

		case "get":
			c.addJob("mysql/constantes", destino, "CONSTANTES", nil)
			c.addJob("mysql/scan-row", destino, "SCAN", nil)
			c.addJob("mysql/get", destino, "GET", nil)

		case "list":
			if c.tbl.TieneCamposFiltro() {
				c.addJob("mysql/tbl-filtros", destino, "FILTROS", nil)
			}
			c.addJob("mysql/constantes", destino, "CONSTANTES", nil)
			c.addJob("mysql/scan-rows", destino, "SCAN", nil)
			c.addJob("mysql/list", destino, "LIST", nil)

		case "get_by":
			c.addJob("mysql/constantes", destino, "CONSTANTES", nil)
			c.addJob("mysql/scan-row", destino, "SCAN", nil)
			c.addJob("mysql/get_by", destino, "GET_BY", directriz.Values())

		case "list_by":
			if c.tbl.TieneCamposFiltro() {
				c.addJob("mysql/tbl-filtros", destino, "FILTROS", nil)
			}
			c.addJob("mysql/constantes", destino, "CONSTANTES", nil)
			c.addJob("mysql/scan-rows", destino, "SCAN", nil)
			c.addJob("mysql/list_by", destino, "LIST_BY", directriz.Values())

		default:
			c.addErr(op.Msgf("directriz no reconocida: '%v'", directriz.Key()))
		}
	}
}

// ================================================================ //
// ========== RENDERIZAR TRABAJOS ================================= //

// Valida y ejecuta todos los trabajos de generación de código en buffer.
func (c *tblGenerator) Generar() (err error) {
	op := gko.Op("tblgen.Generar")
	if len(c.errores) > 0 {
		return op.Msg(c.Errs())
	}
	if c.tbl == nil {
		return op.Msg("tabla es nil")
	}
	if len(c.destinos) == 0 {
		return op.Msg("no hay trabajos por realizar")
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
	for i := range c.destinos {
		// Generar todos los trabajos pendientes en buffer
		for _, job := range c.destinos[i].jobs {
			if job.tmpl == "mysql/list_by" && len(job.campos) == 0 {
				return op.Msg("no se seleccionón ningún campo para list_by")
			}
			if job.tmpl == "mysql/get_by" && len(job.campos) == 0 {
				return op.Msg("no se seleccionón ningún campo para get_by")
			}
			c.tbl.CamposSeleccionados = nil
			c.tbl.CamposSeleccionados = job.campos

			data := map[string]any{
				"TablaOrConsulta": c.tbl,
				"Tabla":           c.tbl,
			}
			if job.titulo != "" {
				textutils.ImprimirSeparador(c.destinos[i].buf, strings.ToUpper(job.titulo))
			}

			switch {
			case strings.HasPrefix(job.tmpl, "html"):
				err = c.gen.renderer.HaciaBufferHTML(job.tmpl, data, c.destinos[i].buf)

			case strings.Contains(job.tmpl, "create_table"):
				err = c.gen.renderer.HaciaBuffer(job.tmpl, data, c.destinos[i].buf)

			default:
				err = c.gen.renderer.HaciaBufferGo(job.tmpl, data, c.destinos[i].buf)
			}
			if err != nil {
				return op.Err(err)
			}
		}
	}
	return nil
}

// ================================================================ //
// ========== ESCRIBIR ============================================ //

// Guarda los buffers generados en los archivos correspondientes.
func (c *tblGenerator) ToFile() (err error) {
	op := gko.Op("tblgen.ToFile")
	for _, dest := range c.destinos {
		if dest.mkdir {
			err := os.MkdirAll(filepath.Dir(dest.filename), 0755)
			if err != nil {
				return op.Err(err)
			}
			c.addHecho(filepath.Dir("[NEW_DIR] " + dest.filename))
		}
		sobreescrito := false
		if fileutils.Existe(dest.filename) {
			sobreescrito = true
		}
		err = fileutils.GuardarGoCode(dest.filename, dest.buf.String())
		if err != nil {
			return op.Err(err)
		}
		if sobreescrito {
			c.addHecho("[REWRITE] " + dest.filename)
		} else {
			c.addHecho("[NEWFILE] " + dest.filename)
		}
	}
	return nil
}

func (c *tblGenerator) ToString() string {
	var buf bytes.Buffer
	for _, dest := range c.destinos {
		textutils.ImprimirSeparador(&buf, dest.filename)
		buf.WriteString(dest.buf.String())
	}
	return buf.String()
}
