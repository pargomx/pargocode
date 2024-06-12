package sqlitedb

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	_ "github.com/glebarez/go-sqlite"
)

var pragmaConfig = "?_pragma=foreign_keys(1)&_busy_timeout=1000"

// Wrapper para "database/sql" con sqlite que permite loggear sentencias.
type SqliteDB struct {
	dbPath string
	db     *sql.DB
	log    bool
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

func (s *SqliteDB) Close() error {
	return s.db.Close()
}

// ================================================================ //

// Inicia una conexión con una base de datos SQLite. Ejemplo "database.db".
func NuevoRepositorio(dbPath string, migracionesFS fs.FS) (*SqliteDB, error) {

	if dbPath == "" {
		return nil, errors.New("database path no especificada")
	}

	// Crear directorio si no existe.
	_, err := os.Stat(path.Dir(dbPath))
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Creado directorio para base de datos", path.Dir(dbPath))
		err := os.MkdirAll(path.Dir(dbPath), 0755)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// Verificar o crear archivo para base de datos.
	_, err = os.Stat(dbPath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Creado archivo para base de datos", dbPath)
		err = os.WriteFile(dbPath, []byte{}, 0664)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath+pragmaConfig)
	if err != nil {
		return nil, err
	}

	repo := &SqliteDB{
		dbPath: dbPath,
		db:     db,
	}

	err = repo.verificarMigraciones(migracionesFS)
	if err != nil {
		db.Close()
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
