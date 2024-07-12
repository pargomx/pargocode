package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"monorepo/assets"
	"monorepo/domain_driven_design/gencode"
	"monorepo/domain_driven_design/sqliteddd"
	"monorepo/htmltmpl"
	"monorepo/migraciones"
	"monorepo/sqlitedb"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/plantillas"
)

// Generador de código
var codeGenerator *gencode.Generador

var puerto = ""       // Default: 5050
var directorio = ""   // Default: directorio actual
var databasePath = "" // Default: _pargo/pargo.sqlite

type servidor struct {
	db  *sqlitedb.SqliteDB
	ddd *sqliteddd.Repositorio
}

func main() {

	// Setup
	gecko.MostrarMensajeEnErrores = true
	gecko.PrintLogTimestamps = false

	// Parámetros de ejecución
	flag.StringVar(&directorio, "dir", "", "directorio raíz del proyecto")
	flag.StringVar(&databasePath, "db", "entidades.db", "ubicación de la db sqlite")
	flag.StringVar(&puerto, "p", "5050", "el servidor escuchará en este puerto")
	flag.Parse()
	if directorio != "" {
		err := os.Chdir(directorio)
		if err != nil {
			fatal("directorio de proyecto inválido: " + err.Error())
		}
	}

	// Repositorio
	sqliteDB, err := sqlitedb.NuevoRepositorio(databasePath, migraciones.MigracionesFS)
	if err != nil {
		fatal(err.Error())
	}

	// Generador de código
	codeGenerator, err = gencode.NuevoGeneradorDeCodigo()
	if err != nil {
		fatal(err.Error())
	}

	// Servicios
	dddrepo := sqliteddd.NuevoRepositorio(sqliteDB)
	srv := &servidor{db: sqliteDB, ddd: dddrepo}

	tpls, err := plantillas.NuevoServicioPlantillasEmbebidas(htmltmpl.PlantillasFS, "")
	if err != nil {
		fatal(err.Error())
	}
	g := gecko.New()
	g.Renderer = tpls
	g.TmplBaseLayout = "app/layout"

	g.StaticFS("/assets", assets.AssetsFS)
	g.FileFS("/favicon.ico", "img/favicon.ico", assets.AssetsFS)

	// ================================================================ //
	// ================================================================ //

	g.GET("/", getInicio)

	g.GET("/mapa", srv.getMapaEntidadRelacion)
	g.GET("/paquetes", srv.getPaquetes)
	g.POS("/paquetes", srv.agregarPaquete)
	g.PUT("/paquetes/{paquete_id}", srv.actualizarPaquete)
	g.DEL("/paquetes/{paquete_id}", srv.eliminarPaquete)
	g.GET("/paquetes/{paquete_id}", srv.getMapaEntidadRelacionPaquete)

	g.GET("/tablas", srv.getPaquetes)              // 1. Tablas en el proyecto
	g.GET("/tablas/nueva", srv.getTablaNueva)      // 2. Formulario para nueva tabla
	g.POS("/tablas/nueva", srv.postTablaNueva)     // 3. Crear nueva tabla
	g.GET("/tablas/{tabla_id}", srv.getTabla)      // 4. Dashboard para tabla
	g.PUT("/tablas/{tabla_id}", srv.putTabla)      // 5. Actualizar datos de la tabla
	g.DEL("/tablas/{tabla_id}", srv.eliminarTabla) // 6. Eliminar tabla
	g.POS("/tablas/{tabla_id}/campos", srv.postCampo)
	g.PUT("/tablas/{tabla_id}/campos", srv.putCampo)
	g.GET("/tablas/{tabla_id}/generar", srv.generarDeTabla)
	g.PUT("/tablas/{tabla_id}/generar/{tipo}", srv.generarDeTablaArchivos)

	g.PUT("/campos/{campo_id}", srv.updateCampo)
	g.DEL("/campos/{campo_id}", srv.deleteCampo)
	g.GET("/campos/{campo_id}/form", srv.formCampo)
	g.PUT("/campos/{campo_id}/reordenar", srv.reordenarCampo)
	g.GET("/campos/{campo_id}/enum", srv.getEnumCampo)
	g.POS("/campos/{campo_id}/enum", srv.postEnumCampo)

	g.POS("/consultas", srv.crearConsulta)
	g.GET("/consultas/nueva", srv.formNuevaConsulta)
	g.GET("/consultas/{consulta_id}", srv.getConsulta)
	g.DEL("/consultas/{consulta_id}", srv.deleteConsulta)
	g.PUT("/consultas/{consulta_id}", srv.actualizarConsulta)
	g.GET("/consultas/{consulta_id}/generar", srv.generarDeConsulta)
	g.POS("/consultas/{consulta_id}/relaciones", srv.postRelacionConsulta)
	g.PUT("/consultas/{consulta_id}/relaciones/{posicion}", srv.actualizarRelacionConsulta)
	g.DEL("/consultas/{consulta_id}/relaciones/{posicion}", srv.eliminarRelacionConsulta)

	g.POS("/consultas/{consulta_id}/campos", srv.postCampoConsulta)
	g.PUT("/consultas/{consulta_id}/campos/{posicion}", srv.actualizarCampoConsulta)
	g.DEL("/consultas/{consulta_id}/campos/{posicion}", srv.eliminarCampoConsulta)
	g.PUT("/consultas/{consulta_id}/reordenar-campo", srv.reordenarCampoConsulta)

	// LOG SQLITE
	g.GET("/log", func(c *gecko.Context) error { sqliteDB.ToggleLog(); return c.StatusOk("Log toggled") })
	// sqliteDB.ToggleLog()

	// ================================================================ //
	// ================================================================ //

	// Handle interrupt.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range ch {
			err = sqliteDB.Close()
			if err != nil {
				fmt.Println("sqliteDB.Close: ", err.Error())
			}
			fmt.Println("")
			gecko.LogInfof("servidor terminado: " + sig.String())
			os.Exit(0)
		}
	}()

	// Listen and serve
	serv := http.Server{
		Addr:    ":" + puerto,
		Handler: g,
	}
	gecko.LogInfof("pargo escuchando en :%v", puerto)
	err = serv.ListenAndServe()
	if err != nil {
		fatal(err.Error())
	}
}

// ================================================================ //
// ========== UTILS =============================================== //

func fatal(msg string) {
	fmt.Println("[FATAL] " + msg)
	os.Exit(1)
}
