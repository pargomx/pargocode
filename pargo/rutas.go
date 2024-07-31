package main

import (
	"flag"
	"fmt"
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
	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/plantillas"
)

// Información de compilación establecida con:
//
//	BUILD_INFO="$(date -I):$(git log --format="%H" -n 1)"
//	go BUILD_INFO -ldflags "-X main.BUILD_INFO=$BUILD_INFO -X main.ambiente=dev"
var BUILD_INFO string // Información de compilación [ fecha:commit_hash ]
var AMBIENTE string   // Ambiente de ejecución [ dev / prod ]

type configs struct {
	puerto       int    // Puerto TCP del servidor
	directorio   string // Default: directorio actual
	databasePath string // Default: _pargo/pargo.sqlite
}

// Generador de código
var codeGenerator *gencode.Generador

type servidor struct {
	cfg   configs
	gecko *gecko.Gecko
	db    *sqlitedb.SqliteDB
	ddd   *sqliteddd.Repositorio
}

func main() {
	gko.LogInfof("Versión:%s:%s", BUILD_INFO, AMBIENTE)

	s := servidor{
		gecko: gecko.New(),
	}

	flag.StringVar(&s.cfg.directorio, "dir", "", "directorio raíz del proyecto")
	flag.StringVar(&s.cfg.databasePath, "db", "entidades.db", "ubicación de la db sqlite")
	flag.IntVar(&s.cfg.puerto, "p", 5051, "el servidor escuchará en este puerto")
	flag.Parse()
	if s.cfg.directorio != "" {
		err := os.Chdir(s.cfg.directorio)
		if err != nil {
			gko.FatalError(err)
		}
	}
	var err error

	// Repositorio
	s.db, err = sqlitedb.NuevoRepositorio(s.cfg.databasePath, migraciones.MigracionesFS)
	if err != nil {
		gko.FatalError(err)
	}
	s.ddd = sqliteddd.NuevoRepositorio(s.db)

	// Generador de código
	codeGenerator, err = gencode.NuevoGeneradorDeCodigo()
	if err != nil {
		gko.FatalError(err)
	}

	tpls, err := plantillas.NuevoServicioPlantillasEmbebidas(htmltmpl.PlantillasFS, "")
	if err != nil {
		gko.FatalError(err)
	}
	s.gecko.Renderer = tpls
	s.gecko.TmplBaseLayout = "app/layout"

	s.gecko.StaticFS("/assets", assets.AssetsFS)
	s.gecko.FileFS("/favicon.ico", "img/favicon.ico", assets.AssetsFS)

	// ================================================================ //
	// ================================================================ //

	s.gecko.GET("/", getInicio)

	s.gecko.GET("/mapa", s.getMapaEntidadRelacion)
	s.gecko.GET("/paquetes", s.getPaquetes)
	s.gecko.POS("/paquetes", s.agregarPaquete)
	s.gecko.PUT("/paquetes/{paquete_id}", s.actualizarPaquete)
	s.gecko.DEL("/paquetes/{paquete_id}", s.eliminarPaquete)
	s.gecko.GET("/paquetes/{paquete_id}", s.getMapaEntidadRelacionPaquete)

	s.gecko.GET("/tablas", s.getPaquetes)              // 1. Tablas en el proyecto
	s.gecko.GET("/tablas/nueva", s.getTablaNueva)      // 2. Formulario para nueva tabla
	s.gecko.POS("/tablas/nueva", s.postTablaNueva)     // 3. Crear nueva tabla
	s.gecko.GET("/tablas/{tabla_id}", s.getTabla)      // 4. Dashboard para tabla
	s.gecko.PUT("/tablas/{tabla_id}", s.putTabla)      // 5. Actualizar datos de la tabla
	s.gecko.DEL("/tablas/{tabla_id}", s.eliminarTabla) // 6. Eliminar tabla
	s.gecko.POS("/tablas/{tabla_id}/campos", s.postCampo)
	s.gecko.PUT("/tablas/{tabla_id}/campos", s.putCampo)
	s.gecko.GET("/tablas/{tabla_id}/generar", s.generarDeTabla)
	s.gecko.PUT("/tablas/{tabla_id}/generar/{tipo}", s.generarDeTablaArchivos)

	s.gecko.PUT("/campos/{campo_id}", s.updateCampo)
	s.gecko.DEL("/campos/{campo_id}", s.deleteCampo)
	s.gecko.GET("/campos/{campo_id}/form", s.formCampo)
	s.gecko.PUT("/campos/{campo_id}/reordenar", s.reordenarCampo)
	s.gecko.GET("/campos/{campo_id}/enum", s.getEnumCampo)
	s.gecko.POS("/campos/{campo_id}/enum", s.postEnumCampo)

	s.gecko.POS("/consultas", s.crearConsulta)
	s.gecko.GET("/consultas/nueva", s.formNuevaConsulta)
	s.gecko.GET("/consultas/{consulta_id}", s.getConsulta)
	s.gecko.DEL("/consultas/{consulta_id}", s.deleteConsulta)
	s.gecko.PUT("/consultas/{consulta_id}", s.actualizarConsulta)
	s.gecko.GET("/consultas/{consulta_id}/generar", s.generarDeConsulta)
	s.gecko.POS("/consultas/{consulta_id}/relaciones", s.postRelacionConsulta)
	s.gecko.PUT("/consultas/{consulta_id}/relaciones/{posicion}", s.actualizarRelacionConsulta)
	s.gecko.DEL("/consultas/{consulta_id}/relaciones/{posicion}", s.eliminarRelacionConsulta)

	s.gecko.POS("/consultas/{consulta_id}/campos", s.postCampoConsulta)
	s.gecko.PUT("/consultas/{consulta_id}/campos/{posicion}", s.actualizarCampoConsulta)
	s.gecko.DEL("/consultas/{consulta_id}/campos/{posicion}", s.eliminarCampoConsulta)
	s.gecko.PUT("/consultas/{consulta_id}/reordenar-campo", s.reordenarCampoConsulta)

	// LOG SQLITE
	s.gecko.GET("/log", func(c *gecko.Context) error { s.db.ToggleLog(); return c.StatusOk("Log toggled") })

	// ================================================================ //
	// ================================================================ //

	// Handle interrupt.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range ch {
			err = s.db.Close()
			if err != nil {
				fmt.Println("sqliteDB.Close: ", err.Error())
			}
			fmt.Println("")
			gko.LogInfof("servidor terminado: " + sig.String())
			os.Exit(0)
		}
	}()

	err = s.gecko.IniciarEnPuerto(s.cfg.puerto)
	if err != nil {
		gko.FatalError(err)
	}
}
