package codegenerator

import (
	"bytes"
	"embed"
	"fmt"
	"monorepo/ddd"
	"monorepo/fileutils"
	"monorepo/textutils"
	"monorepo/tmplutils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pargomx/gecko/gko"
)

// ================================================================ //
// ========== GENERADOR DE CÓDIGO ================================= //

//go:embed plantillas
var plantillasFS embed.FS

// Colección para generar de varias tablas y/o consultas.
type Generadores []generador

type generador struct {
	renderer *tmplutils.Renderer

	destinos []genDest
	hechos   []string
	errores  []error

	tbl *tabla
	con *consulta

	sinTitulos bool
}

type genDest struct {
	jobs     []genJob      // trabajos por realizar
	buf      *bytes.Buffer // buffer para escribir el código
	filename string        // nombre del archivo destino (si aplica)
	mkdir    bool          // crear directorio si no existe
}

type genJob struct {
	tmpl      string          // nombre de la plantilla a renderizar
	titulo    string          // título del bloque para separador
	custom    *customList     // parámetros para consulta custom
	camposTbl []CampoTabla    // campos seleccionados
	camposCon []CampoConsulta // campos seleccionados
}

type tabla struct {
	Tabla   ddd.Tabla
	Campos  []CampoTabla
	Paquete ddd.Paquete

	CamposSeleccionados []CampoTabla
	CustomList          *customList

	Sqlite bool
	MySQL  bool
}

type consulta struct {
	Paquete     ddd.Paquete
	Consulta    ddd.Consulta
	TablaOrigen ddd.Tabla
	From        tabla
	Campos      []CampoConsulta
	Relaciones  []Relacion

	CamposSeleccionados []CampoConsulta
	CustomList          *customList

	Sqlite bool
	MySQL  bool
}

// ================================================================ //

func (c *generador) addErr(err error) {
	if err != nil {
		c.errores = append(c.errores, err)
	}
}
func (c *generador) Errs() string {
	txt := ""
	for _, err := range c.errores {
		txt += err.Error() + "\n"
	}
	txt = strings.TrimSuffix(txt, "\n")
	return txt
}

func (c *generador) addHecho(hecho string) {
	c.hechos = append(c.hechos, hecho)
}
func (c *generador) GetHechos() []string {
	return c.hechos
}

// ================================================================ //
// ========== ORIGEN DEL GENERADOR DE CÓDIGO ====================== //

// Nuevo generador de código para una tabla.
func NuevoDeTabla(repo Repositorio, tablaID int) (*generador, error) {
	op := gko.Op("codegen.NuevoDeTabla")
	renderer, err := tmplutils.NuevoRenderer(plantillasFS, "plantillas")
	if err != nil {
		return nil, op.Err(err)
	}
	if repo == nil {
		return nil, op.Str("repo es nil")
	}
	tbl, err := getTabla(repo, tablaID)
	if err != nil {
		return nil, op.Err(err)
	}
	tblGenCall := generador{
		renderer: renderer,
		tbl:      tbl,
	}
	return &tblGenCall, nil
}

func NuevoDeConsulta(repo Repositorio, consultaID int) (*generador, error) {
	op := gko.Op("codegen.NuevoDeConsulta")
	renderer, err := tmplutils.NuevoRenderer(plantillasFS, "plantillas")
	if err != nil {
		return nil, op.Err(err)
	}
	if repo == nil {
		return nil, op.Str("repo es nil")
	}
	con, err := getConsulta(repo, consultaID)
	if err != nil {
		return nil, op.Err(err)
	}
	tblGenCall := generador{
		renderer: renderer,
		con:      con,
	}
	return &tblGenCall, nil
}

func NuevoDePaquete(repo Repositorio, paqueteID int) (Generadores, error) {
	op := gko.Op("codegen.NuevoDePaquete")
	renderer, err := tmplutils.NuevoRenderer(plantillasFS, "plantillas")
	if err != nil {
		return nil, op.Err(err)
	}
	if repo == nil {
		return nil, op.Str("repo es nil")
	}
	Generadores := Generadores{}
	tablas, err := repo.ListTablasByPaqueteID(paqueteID)
	if err != nil {
		return nil, op.Err(err)
	}
	for _, t := range tablas {
		tbl, err := getTabla(repo, t.TablaID)
		if err != nil {
			return nil, op.Err(err)
		}
		call := generador{
			renderer: renderer,
			tbl:      tbl,
		}
		Generadores = append(Generadores, call)
	}
	consultas, err := repo.ListConsultasByPaqueteID(paqueteID)
	if err != nil {
		return nil, op.Err(err)
	}
	for _, c := range consultas {
		con, err := getConsulta(repo, c.ConsultaID)
		if err != nil {
			return nil, op.Err(err)
		}
		call := generador{
			renderer: renderer,
			con:      con,
		}
		Generadores = append(Generadores, call)
	}
	return Generadores, nil
}

// ================================================================ //
// ========== DEFINIR TRABAJO ===================================== //

// Definir qué se va a generar. Dependiendo de lo solicitado se mandará guardar
// al lugar adecuado. Pueden ser diferentes archivos o buffers de destino.
func (c *generador) PrepararJob(tipo string) *generador {
	if tipo == "" {
		c.addErr(gko.ErrDatoIndef().Str("tipo de trabajo indefinido").Op("PrepararJob"))
		return c
	}
	switch tipo {
	case "entidad":
		c.addJobsEntidad()
	case "mysql":
		if c.tbl != nil {
			c.addJobsRepoSQLTabla(false)
		} else if c.con != nil {
			c.addJobsRepoSQLConsulta(false)
		}
	case "sqlite":
		if c.tbl != nil {
			c.addJobsRepoSQLTabla(true)
		} else if c.con != nil {
			c.addJobsRepoSQLConsulta(true)
		}
	default:
		c.addJob(tipo, "generado.go", "")
	}
	return c
}

// Imprime en consola los trabajos del generador.
func (c *generador) DescribirJobs() {
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

func (c *generador) SinTitulos() *generador {
	c.sinTitulos = true
	return c
}

// ================================================================ //
// ========== PREPARAR TRABAJO PENDIENTE ========================== //

// Agregar un trabjo a la lista de pendientes con el destino adecuado.
// El título es opcional y se usa para agregar un comentario separador.
// Si ya se había agregado un trabajo con el mismo tmpl, destino y título
// se devolverá el trabajo existente sin agregar uno nuevo.
func (c *generador) addJob(tmpl string, destino string, titulo string) *genJob {
	op := gko.Op("addJob").Ctx("tmpl", tmpl)
	if c.tbl == nil && c.con == nil {
		c.addErr(op.Str("tabla o consulta indefinida"))
		return nil
	}
	if tmpl == "" {
		c.addErr(op.Str("plantilla no definida"))
		return nil
	}
	// no agregar trabajo si ya se declaró
	for i, dest := range c.destinos {
		if dest.filename == destino {
			for y, job := range dest.jobs {
				if job.tmpl == tmpl && job.titulo == titulo {
					return &c.destinos[i].jobs[y]
				}
			}
		}
	}
	// agregar trabajo al destino adecuado
	destIdx, err := c.addDestino(destino)
	if err != nil {
		c.addErr(op.Err(err))
		return nil
	}
	job := genJob{
		tmpl:   tmpl,
		titulo: titulo,
	}
	c.destinos[destIdx].jobs = append(c.destinos[destIdx].jobs, job)
	return &c.destinos[destIdx].jobs[len(c.destinos[destIdx].jobs)-1]
}

// ================================================================ //

func (job *genJob) setCustomList(custom *customList) {
	job.custom = custom
}

// agregar campos seleccionados para get_by y list_by
func (job *genJob) setCamposSelec(c *generador, campos []string) {
	if c.tbl != nil {
		for _, v := range campos {
			campo, err := c.tbl.BuscarCampo(v)
			if err != nil {
				c.addErr(gko.Op("withCampos").Err(err).Ctx("tbl", c.tbl.NombreRepo()))
				continue
			}
			job.camposTbl = append(job.camposTbl, *campo)
		}
	} else if c.con != nil {
		for _, v := range campos {
			campo, err := c.con.BuscarCampo(v)
			if err != nil {
				c.addErr(gko.Op("withCampos").Err(err).Ctx("con", c.con.Consulta.NombreItem))
				continue
			}
			job.camposCon = append(job.camposCon, *campo)
		}
	}
}

// ================================================================ //

// Devuelve el idx del addDestino declarado (lo agrega si no existe).
func (c *generador) addDestino(filename string) (int, error) {
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

func (c *generador) addJobsEntidad() {
	if c.tbl != nil {
		destino := filepath.Join(c.tbl.Paquete.Directorio, c.tbl.Paquete.Nombre, "t_"+c.tbl.Tabla.Kebab+".go")
		oldFilename(filepath.Join(c.tbl.Paquete.Directorio, c.tbl.Paquete.Nombre, c.tbl.Tabla.Kebab+".go"), destino)
		c.addJob("go/tbl_struct", destino, "")
		c.addJob("go/tbl_propiedades", destino, "")

	} else if c.con != nil {
		destino := filepath.Join(c.con.Paquete.Directorio, c.con.Paquete.Nombre, "q_"+c.con.Consulta.NombreItem+".go")
		oldFilename(filepath.Join(c.con.Paquete.Directorio, c.con.Paquete.Nombre, c.con.TablaOrigen.NombreRepo+"_extendido.go"), destino)
		c.addJob("go/qry_struct", destino, "")
	}
}

func (generadores Generadores) GenerarSchemaSQLite(tipoJob string, tipoDB string) (hechos []string, err error) {
	op := gko.Op("GenerarSchemaSQLite")

	buf := &strings.Builder{}
	filename := "migraciones/new_schema.sql"

	for _, c := range generadores {
		if c.tbl == nil {
			continue
		}
		destino := filepath.Join(c.tbl.Paquete.Directorio, c.tbl.Paquete.Nombre+".sql")
		if tipoJob == "migracion" {
			c.addJob("sqlite/insert_into_select", destino, "")
			filename = "migraciones/migracion_full.sql"
		} else if tipoDB == "sqlite" {
			c.addJob("sqlite/create_table", destino, "")
		} else if tipoDB == "mysql" {
			c.addJob("mysql/create_table", destino, "")
		} else {
			return nil, op.Msgf("Trabajo inválido: %v/%v", tipoDB, tipoJob)
		}
		c.SinTitulos()
		err = c.Generar()
		if err != nil {
			fmt.Fprintf(buf, "\n/* ERROR\n\t%v\n*/\n", err.Error())
		}
		fmt.Fprintf(buf, "\n%v\n", c.ToString())
	}

	// Guardar archivo
	if !fileutils.Existe(filepath.Dir(filename)) {
		err = os.MkdirAll(filepath.Dir(filename), 0755)
		if err != nil {
			return nil, op.Err(err)
		}
		hechos = append(hechos, filepath.Dir("[NEW_DIR] "+filename))
	}
	sobreescrito := false
	if fileutils.Existe(filename) {
		sobreescrito = true
	}
	err = fileutils.GuardarPlainText(filename, buf.String())
	if err != nil {
		return nil, op.Err(err)
	}
	if sobreescrito {
		hechos = append(hechos, "[REWRITE] "+filename)
	} else {
		hechos = append(hechos, "[NEWFILE] "+filename)
	}
	return hechos, nil
}

// ================================================================ //

func (c *generador) addJobsRepoSQLTabla(sqlite bool) {
	op := gko.Op("addJobsRepoSQL")
	if c.tbl == nil {
		c.addErr(op.Str("tabla es nil"))
		return
	}
	op.Ctx("tabla", c.tbl.Tabla.NombreRepo)
	if len(c.tbl.Directrices()) == 0 {
		c.addErr(op.Str("no hay directrices"))
	}
	c.tbl.Sqlite = sqlite
	c.tbl.MySQL = !sqlite

	// Si el paquete aún no tiene el archivo de servicio, se agrega.
	if sqlite && !fileutils.Existe(filepath.Join(c.tbl.Paquete.Directorio, "sqlite"+c.tbl.Paquete.Nombre, "servicio_repo.go")) {
		c.addJob("sqlite/servicio", filepath.Join(c.tbl.Paquete.Directorio, "sqlite"+c.tbl.Paquete.Nombre, "servicio_repo.go"), "")

	} else if !sqlite && !fileutils.Existe(filepath.Join(c.tbl.Paquete.Directorio, "mysql"+c.tbl.Paquete.Nombre, "servicio_repo.go")) {
		c.addJob("mysql/servicio", filepath.Join(c.tbl.Paquete.Directorio, "mysql"+c.tbl.Paquete.Nombre, "servicio_repo.go"), "")
	}

	// Destino diferente dependiendo si es repo mysql o sqlite.
	destino := "generado.go"
	if sqlite {
		destino = filepath.Join(c.tbl.Paquete.Directorio, "sqlite"+c.tbl.Paquete.Nombre, "s_"+c.tbl.Tabla.NombreRepo+"_gen.go")
	} else {
		destino = filepath.Join(c.tbl.Paquete.Directorio, "mysql"+c.tbl.Paquete.Nombre, "s_"+c.tbl.Tabla.NombreRepo+"_gen.go")
		oldFilename(filepath.Join(c.tbl.Paquete.Directorio, "mysql"+c.tbl.Paquete.Nombre, c.tbl.Tabla.Kebab+".go"), destino)
	}

	c.addJob("mysql/paquete", destino, "")
	for _, directriz := range c.tbl.Directrices() {
		switch directriz.Key() {
		case "insert":
			c.addJob("mysql/tbl-insert", destino, "INSERT")

		case "update":
			c.addJob("mysql/tbl-update", destino, "UPDATE")

		case "insert_update":
			c.addJob("mysql/tbl-insert_update", destino, "INSERT_UPDATE")

		case "existe":
			c.addJob("mysql/existe", destino, "EXISTE")

		case "delete":
			c.addJob("mysql/existe", destino, "EXISTE")
			c.addJob("mysql/tbl-delete", destino, "DELETE")

		case "fetch":
			c.addJob("mysql/constantes", destino, "CONSTANTES")
			c.addJob("mysql/scan-row", destino, "SCAN")
			c.addJob("mysql/fetch", destino, "FETCH")

		case "get":
			c.addJob("mysql/constantes", destino, "CONSTANTES")
			c.addJob("mysql/scan-row", destino, "SCAN")
			c.addJob("mysql/get", destino, "GET")

		case "list":
			if c.tbl.TieneCamposFiltro() {
				c.addJob("mysql/tbl-filtros", destino, "FILTROS")
			}
			c.addJob("mysql/constantes", destino, "CONSTANTES")
			c.addJob("mysql/scan-rows", destino, "SCAN")
			c.addJob("mysql/list", destino, "LIST")

		case "get_by":
			c.addJob("mysql/constantes", destino, "CONSTANTES")
			c.addJob("mysql/scan-row", destino, "SCAN")
			c.addJob("mysql/get_by", destino, "GET_BY "+strings.Join(directriz.Values(), " ")).setCamposSelec(c, directriz.Values())

		case "list_by":
			if c.tbl.TieneCamposFiltro() {
				c.addJob("mysql/tbl-filtros", destino, "FILTROS")
			}
			c.addJob("mysql/constantes", destino, "CONSTANTES")
			c.addJob("mysql/scan-rows", destino, "SCAN")
			c.addJob("mysql/list_by", destino, "LIST_BY "+strings.Join(directriz.Values(), " ")).setCamposSelec(c, directriz.Values())

		case "list_custom":
			customList, err := directriz.CustomList()
			if err != nil {
				c.addErr(op.Err(err))
				continue
			}
			c.addJob("mysql/constantes", destino, "CONSTANTES")
			c.addJob("mysql/scan-rows", destino, "SCAN")
			c.addJob("mysql/list_custom", destino, "LIST "+customList.CompFunc).setCustomList(customList)

		default:
			c.addErr(op.Msgf("directriz no reconocida: '%v'", directriz.Key()))
		}
	}
}

// ================================================================ //

func (c *generador) addJobsRepoSQLConsulta(sqlite bool) {
	op := gko.Op("addJobsConsultaRepoSQL")
	if c.con == nil {
		c.addErr(op.Str("consulta es nil"))
		return
	}
	op.Ctx("consulta", c.con.Consulta.NombreItem)
	if len(c.con.Directrices()) == 0 {
		c.addErr(op.Str("no hay directrices"))
	}
	c.con.Sqlite = sqlite
	c.con.MySQL = !sqlite

	// Si el paquete aún no tiene el archivo de servicio, se agrega.
	if sqlite && !fileutils.Existe(filepath.Join(c.con.Paquete.Directorio, "sqlite"+c.con.Paquete.Nombre, "servicio_repo.go")) {
		c.addJob("sqlite/servicio", filepath.Join(c.con.Paquete.Directorio, "sqlite"+c.con.Paquete.Nombre, "servicio_repo.go"), "")
	} else if !sqlite && !fileutils.Existe(filepath.Join(c.con.Paquete.Directorio, "mysql"+c.con.Paquete.Nombre, "servicio_repo.go")) {
		c.addJob("mysql/servicio", filepath.Join(c.con.Paquete.Directorio, "mysql"+c.con.Paquete.Nombre, "servicio_repo.go"), "")
	}

	// Destino diferente dependiendo si es repo mysql o sqlite.
	destino := "generado.go"
	if sqlite {
		destino = filepath.Join(c.con.Paquete.Directorio, "sqlite"+c.con.Paquete.Nombre, "s_"+c.con.Consulta.NombreItem+"_gen.go")
	} else {
		destino = filepath.Join(c.con.Paquete.Directorio, "mysql"+c.con.Paquete.Nombre, "s_"+c.con.Consulta.NombreItem+"_gen.go")
		oldFilename(filepath.Join(c.con.Paquete.Directorio, "mysql"+c.con.Paquete.Nombre, c.con.TablaOrigen.NombreRepo+"_extendido.go"), destino)
	}

	c.addJob("mysql/paquete", destino, "")
	for _, directriz := range c.con.Directrices() {
		switch directriz.Key() {

		case "fetch":
			c.addJob("mysql/constantes", destino, "CONSTANTES")
			c.addJob("mysql/scan-row", destino, "SCAN")
			c.addJob("mysql/fetch", destino, "FETCH")

		case "get":
			c.addJob("mysql/constantes", destino, "CONSTANTES")
			c.addJob("mysql/scan-row", destino, "SCAN")
			c.addJob("mysql/get", destino, "GET")

		case "list":
			if c.con.TieneCamposFiltro() {
				c.addJob("mysql/qry-filtros", destino, "FILTROS "+c.con.Consulta.NombreItem)
			}
			c.addJob("mysql/constantes", destino, "CONSTANTES")
			c.addJob("mysql/scan-rows", destino, "SCAN")
			c.addJob("mysql/list", destino, "LIST")

		case "get_by":
			c.addJob("mysql/constantes", destino, "CONSTANTES")
			c.addJob("mysql/scan-row", destino, "SCAN")
			c.addJob("mysql/get_by", destino, "GET_BY "+strings.Join(directriz.Values(), " ")).setCamposSelec(c, directriz.Values())

		case "list_by":
			if c.con.TieneCamposFiltro() {
				c.addJob("mysql/qry-filtros", destino, "FILTROS "+c.con.Consulta.NombreItem)
			}
			c.addJob("mysql/constantes", destino, "CONSTANTES")
			c.addJob("mysql/scan-rows", destino, "SCAN")
			c.addJob("mysql/list_by", destino, "LIST_BY "+strings.Join(directriz.Values(), " ")).setCamposSelec(c, directriz.Values())

		case "list_custom":
			customList, err := directriz.CustomList()
			if err != nil {
				c.addErr(op.Err(err))
				continue
			}
			c.addJob("mysql/constantes", destino, "CONSTANTES")
			c.addJob("mysql/scan-rows", destino, "SCAN")
			c.addJob("mysql/list_custom", destino, "LIST "+customList.CompFunc).setCustomList(customList)

		default:
			c.addErr(op.Msgf("directriz no reconocida: '%v'", directriz.Key()))
		}
	}
}

// ================================================================ //
// ========== RENDERIZAR TRABAJOS ================================= //

// Valida y ejecuta todos los trabajos de generación de código en buffer.
func (c *generador) Generar() (err error) {
	op := gko.Op("tblgen.Generar")
	if len(c.errores) > 0 {
		return op.Msg(c.Errs())
	}
	if c.tbl == nil && c.con == nil {
		return op.Msg("tabla o consulta es nil")
	}
	if len(c.destinos) == 0 {
		return op.Msg("no hay trabajos por realizar")
	}

	// Validar tabla
	if c.tbl != nil {
		op.Ctx("tabla", c.tbl.Tabla.NombreRepo)
		if c.tbl.NombreItem() == "" {
			return op.Msg("nombre de tabla indefinido")
		}
		if len(c.tbl.Campos) == 0 {
			return op.Msg("debe haber al menos una columna")
		}
		if len(c.tbl.PrimaryKeys()) == 0 {
			return op.Msg("debe haber al menos una clave primaria")
		}
	}

	// Validar consulta
	if c.con != nil {
		op.Ctx("consulta", c.con.Consulta.NombreItem)
		if c.con.Consulta.NombreItem == "" {
			return op.Msg("nombre de consulta indefinido")
		}
		if len(c.con.Campos) == 0 {
			return op.Msg("debe haber al menos una columna")
		}
	}

	// Generar todos los trabajos pendientes en buffer para cada destino.
	for i := range c.destinos {
		for _, job := range c.destinos[i].jobs {
			data := map[string]any{}

			if c.tbl != nil {
				if job.tmpl == "mysql/list_by" && len(job.camposTbl) == 0 {
					return op.Msg("no se seleccionón ningún campo para list_by")
				}
				if job.tmpl == "mysql/get_by" && len(job.camposTbl) == 0 {
					return op.Msg("no se seleccionón ningún campo para get_by")
				}
				c.tbl.CamposSeleccionados = nil
				c.tbl.CamposSeleccionados = job.camposTbl
				c.tbl.CustomList = nil
				c.tbl.CustomList = job.custom
				data = map[string]any{
					"TablaOrConsulta": c.tbl,
					"Tabla":           c.tbl,
				}
			}

			if c.con != nil {
				if job.tmpl == "mysql/list_by" && len(job.camposCon) == 0 {
					return op.Msg("no se seleccionón ningún campo para list_by")
				}
				if job.tmpl == "mysql/get_by" && len(job.camposCon) == 0 {
					return op.Msg("no se seleccionón ningún campo para get_by")
				}
				c.con.CamposSeleccionados = nil
				c.con.CamposSeleccionados = job.camposCon
				c.con.CustomList = nil
				c.con.CustomList = job.custom
				data = map[string]any{
					"TablaOrConsulta":  c.con,
					"Consulta":         c.con,
					"AgregadoConsulta": c.con,
				}
			}

			if job.titulo != "" && !c.sinTitulos {
				textutils.ImprimirSeparador(c.destinos[i].buf, strings.ToUpper(job.titulo))
			}

			switch {
			case strings.HasPrefix(job.tmpl, "html"):
				err = c.renderer.HaciaBufferHTML(job.tmpl, data, c.destinos[i].buf)

			case strings.Contains(job.tmpl, "create_table"):
				err = c.renderer.HaciaBuffer(job.tmpl, data, c.destinos[i].buf)

			case strings.Contains(job.tmpl, "insert_into_select"):
				err = c.renderer.HaciaBuffer(job.tmpl, data, c.destinos[i].buf)

			case strings.Contains(job.tmpl, "query"):
				err = c.renderer.HaciaBuffer(job.tmpl, data, c.destinos[i].buf)

			default:
				err = c.renderer.HaciaBufferGo(job.tmpl, data, c.destinos[i].buf)
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
func (c *generador) ToFile() (err error) {
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

func (c *generador) ToString() string {
	var buf bytes.Buffer
	for _, dest := range c.destinos {
		if !c.sinTitulos {
			textutils.ImprimirSeparador(&buf, dest.filename)
		}
		buf.WriteString(dest.buf.String())
	}
	return buf.String()
}

// Renombrar archivo generado con nombre antiguo con git mv.
func oldFilename(oldFilename, newFilename string) {
	if !fileutils.Existe(oldFilename) {
		return
	}
	gko.LogWarn("Renombrando archivo antiguo: " + oldFilename)
	err := exec.Command("git", "mv", oldFilename, newFilename).Run()
	if err != nil {
		gko.FatalError(err)
	}
}
