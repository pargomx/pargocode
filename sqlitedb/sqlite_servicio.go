package sqlitedb

import (
	"database/sql"
	"io/fs"
	"os"
	"path"

	_ "github.com/glebarez/go-sqlite"
	"github.com/pargomx/gecko/gko"
)

// INFO: esta es la versión más actualizada del paquete: 2025-05-28.

// Hay configuraciones que se aplican mediante SQL a la base de datos.
var configPragmaDSN = "?_pragma=foreign_keys(1)&_busy_timeout=1000"

// Wrapper para "database/sql" con sqlite que permite loggear sentencias.
type SqliteDB struct {
	dbPath     string // ruta al archivo de base de datos.
	db         *sql.DB
	backupsDir string // directorio en donde poner backups de base de datos.
	log        bool
}

// Utilizado para que los repositorios del dominio puedan usar DB o Transaccion.
type Ejecutor interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

// ================================================================ //

// Activa o desactiva el log de SQL Statements a la terminal.
func (s *SqliteDB) ToggleLog() {
	s.log = !s.log
}

// Cerrar base de datos.
func (s *SqliteDB) Close() error {
	op := gko.Op("sqlitedb.Close")
	err := s.db.Close()
	if err != nil {
		return op.Err(err)
	}
	return nil
}

// Cerrar base de datos y comprobar que todo esté contenido en un solo archivo.
// En WAL mode para un archivo "app.db" se generan "app.db-shm" y "app.db-wal".
// Si aún están estos archivos puede que algo los mantenga abiertos y por lo tanto
// puede usarse el error para no continuar en operaciones que se quieran hacer con
// el archivo de base de datos.
func (s *SqliteDB) CloseFully() error {
	op := gko.Op("sqlitedb.CloseFully")
	err := s.db.Close()
	if err != nil {
		return op.Err(err)
	}
	if _, err := os.Stat(s.dbPath + "-wal"); err == nil {
		return op.Strf("current db still open: WAL file exists (%v)", s.dbPath+"-wal")
	}
	if _, err := os.Stat(s.dbPath + "-shm"); err == nil {
		return op.Strf("current db still open: SHM file exists (%v)", s.dbPath+"-shm")
	}
	return nil
}

// Abre el archivo de base de datos y se conecta con la configPragmaDSN.
// Confía en que el dbPath ya se comprobó.
func (s *SqliteDB) openDatabase() error {
	var err error
	s.db, err = sql.Open("sqlite", s.dbPath+configPragmaDSN)
	if err != nil {
		return err
	}
	return nil
}

// ================================================================ //

// Inicia una conexión con una base de datos SQLite. Ejemplo "database.db".
// Si el archivo o su directorio no existen, intenta crearlos.
func NuevoRepositorio(dbPath string, migracionesFS fs.FS) (*SqliteDB, error) {
	op := gko.Op("sqlitedb.NewRepo")

	if dbPath == "" {
		return nil, op.Str("database path no especificada")
	}

	// Crear directorio si no existe.
	_, err := os.Stat(path.Dir(dbPath))
	if os.IsNotExist(err) {
		err := os.MkdirAll(path.Dir(dbPath), 0750)
		if err != nil {
			return nil, op.Err(err).Op("NewDatabaseDir")
		}
		gko.LogInfof("SqliteDB: Creado directorio '%v'", path.Dir(dbPath))
	} else if err != nil {
		return nil, op.Err(err)
	}

	// Verificar o crear archivo para base de datos.
	_, err = os.Stat(dbPath)
	if os.IsNotExist(err) {
		err = os.WriteFile(dbPath, []byte{}, 0640)
		if err != nil {
			return nil, op.Err(err).Op("NewDatabaseFile")
		}
		gko.LogInfof("SqliteDB: Creado archivo '%v'", dbPath)
	} else if err != nil {
		return nil, op.Err(err)
	}

	// Abrir repositorio.
	repo := &SqliteDB{dbPath: dbPath, db: nil}
	err = repo.openDatabase()
	if err != nil {
		return nil, op.Err(err)
	}

	err = repo.verificarMigraciones(migracionesFS)
	if err != nil {
		repo.Close()
		return nil, err
	}

	// Debug config de la conexión
	// var pragma int
	// sqliteDB.QueryRow("PRAGMA foreign_keys").Scan(&pragma)
	// fmt.Println("foreign_keys: ", pragma)
	// sqliteDB.QueryRow("PRAGMA busy_timeout").Scan(&pragma)
	// fmt.Println("busy_timeout: ", pragma)

	// Para evitar error database locked. https://github.com/mattn/go-sqlite3/issues/274
	repo.db.SetMaxOpenConns(1)

	return repo, nil
}
