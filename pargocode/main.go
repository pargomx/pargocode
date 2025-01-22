package main

import (
	"flag"
	"fmt"
	"monorepo/htmltmpl"
	"monorepo/migraciones"
	"monorepo/sqlitedb"
	"monorepo/sqliteddd"
	"monorepo/textutils"
	"os"
	"os/signal"
	"syscall"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/plantillas"
)

// Información de compilación establecida con:
//
//	BUILD_INFO="$(date -I):$(git log --format="%H" -n 1)"
//	go build -ldflags "-X main.BUILD_INFO=$BUILD_INFO -X main.AMBIENTE=DEV"
var BUILD_INFO string // Información de compilación [ fecha:commit_hash ]
var AMBIENTE string   // Ambiente de ejecución [ DEV / PROD ]

type servidor struct {
	cfg   configs
	gecko *gecko.Gecko
	db    *sqlitedb.SqliteDB
	ddd   *sqliteddd.Repositorio
	txt   *textutils.Utils
}

type configs struct {
	puerto       int    // Puerto TCP del servidor
	directorio   string // Directorio raíz de la aplicación
	databasePath string // Ruta del archivo base de datos
	logDB        bool   // Loggear consultas a la base de datos
	sourceDir    string // Directorio con assets y htmltmpl para no usar embeded
}

func main() {
	gko.LogInfof("Versión:%s:%s", BUILD_INFO, AMBIENTE)
	s := servidor{}
	var err error

	// Configuraciones
	flag.IntVar(&s.cfg.puerto, "p", 5051, "el servidor escuchará en este puerto")
	flag.StringVar(&s.cfg.directorio, "dir", "", "directorio raíz de la aplicación")
	flag.StringVar(&s.cfg.databasePath, "db", "entidades.db", "ubicación de la db sqlite")
	flag.BoolVar(&s.cfg.logDB, "logdb", false, "loggear consultas a la base de datos")
	flag.StringVar(&s.cfg.sourceDir, "src", "", "directorio con assets y htmltmpl para no usar embeded")
	flag.Parse()

	// Directorio raíz de ejecución
	if s.cfg.directorio != "" {
		err := os.Chdir(s.cfg.directorio)
		if err != nil {
			gko.FatalExit("directorio raíz inválido: " + err.Error())
		}
	}

	// Repositorio
	s.db, err = sqlitedb.NuevoRepositorio(s.cfg.databasePath, migraciones.MigracionesFS)
	if err != nil {
		gko.FatalError(err)
	}
	if s.cfg.logDB {
		s.db.ToggleLog()
	}
	s.ddd = sqliteddd.NuevoRepositorio(s.db)

	// Plantillas HTML
	s.gecko = gecko.New()
	if s.cfg.sourceDir != "" {
		gko.LogInfo("Usando plantillas en " + s.cfg.sourceDir + "/htmltmpl")
		s.gecko.Renderer, err = plantillas.NuevoServicioPlantillas(s.cfg.sourceDir+"/htmltmpl", AMBIENTE == "DEV")
	} else {
		s.gecko.Renderer, err = plantillas.NuevoServicioPlantillasEmbebidas(htmltmpl.PlantillasFS, "")
	}
	if err != nil {
		gko.FatalError(err)
	}
	s.gecko.TmplBaseLayout = "app/layout"
	s.gecko.TmplError = "app/error"

	s.txt = textutils.NewTextUtils()

	// Rutas
	s.registrarRutas()

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
