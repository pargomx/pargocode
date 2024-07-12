package sqlitedb

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"strconv"
	"strings"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== MIGRACIONES ========================================= //

// Verifica que las migraciones estén aplicadas y las aplica si no lo están.
//
// Deben estar en el directorio "migraciones" y tener de prefijo un número
// consecutivo seguido de un guión bajo. Ejemplo: "03_usuarios.sql".
//
// Toda migración debe registrarse a sí misma en la tabla "migraciones" con
// el mismo número de ID que tiene en su prefijo.
func (s *SqliteDB) verificarMigraciones(migracionesFS fs.FS) error {
	const selectMigraciones = "SELECT id, fecha, detalles FROM migraciones"

	rows, err := s.db.Query(selectMigraciones)
	if err != nil {
		// Inicializar base de datos si aún no tiene ni siquiera tabla de migraciones.
		gecko.LogEventof("Aplicando migración 00_setup.sql")
		migracionCero, err := fs.ReadFile(migracionesFS, "00_setup.sql")
		if err != nil {
			return err
		}
		_, err = s.db.Exec(string(migracionCero)) // no se puede configurar db dentro de una transacción.
		if err != nil {
			return err
		}
		rows, err = s.db.Query(selectMigraciones)
		if err != nil {
			return err
		}
		// gecko.LogOkeyf("Inicializada")
	}

	// Obtener migraciones aplicadas.
	type migracionAplicada struct {
		id       int    // `migraciones.id`
		fecha    string // `migraciones.fecha`
		detalles string // `migraciones.detalles`
	}
	aplicadas := []migracionAplicada{}
	for rows.Next() {
		apli := migracionAplicada{}
		err := rows.Scan(&apli.id, &apli.fecha, &apli.detalles)
		if err != nil {
			return err
		}
		aplicadas = append(aplicadas, apli)
	}

	// Obtener migraciones disponibles.
	type migraSourcefile struct {
		id        int
		filename  string
		contenido string
	}
	disponibles := []migraSourcefile{}
	files, err := fs.ReadDir(migracionesFS, ".")
	if err != nil {
		return err
	}
	for i, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			// fmt.Println("ignorando migración no .sql", file.Name())
			continue
		}
		id, err := strconv.Atoi(strings.Split(file.Name(), "_")[0])
		if err != nil {
			return err
		}
		bytes, err := fs.ReadFile(migracionesFS, file.Name())
		if err != nil {
			return err
		}
		if id != i {
			return fmt.Errorf("migración %v tiene id %v pero debería ser %v para ser consecutivo", file.Name(), id, i)
		}
		disponibles = append(disponibles, migraSourcefile{
			id:        id,
			filename:  file.Name(),
			contenido: string(bytes),
		})
	}

	for i, migra := range disponibles {
		// Aplicar migración si no está aplicada.
		if len(aplicadas)-1 < i {
			gecko.LogEventof("Aplicando migración %v", migra.filename)
			tx, err := s.db.Begin()
			if err != nil {
				return err
			}
			_, err = tx.Exec(migra.contenido)
			if err != nil {
				tx.Rollback()
				return err
			}
			aplicado := migracionAplicada{}
			err = tx.QueryRow(selectMigraciones+" WHERE id = ?", migra.id).
				Scan(&aplicado.id, &aplicado.fecha, &aplicado.detalles)
			if err != nil {
				tx.Rollback()
				if errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("%v no se registra en tabla de migraciones con id %v", migra.filename, migra.id)
				}
				return err
			}
			if aplicado.id != migra.id {
				tx.Rollback()
				return fmt.Errorf("migración %v registra un id %v pero el archivo tiene número %v", migra.filename, aplicado.id, migra.id)
			}
			err = tx.Commit()
			if err != nil {
				return err
			}
			aplicadas = append(aplicadas, aplicado)
			// gecko.LogOkeyf("aplicada")
		}
		// Verificar que coincida la migración aplicada con la disponible.
		if migra.id != aplicadas[i].id {
			return fmt.Errorf("migración aplicada con id %v pero debería ser %v para %v", aplicadas[i].id, migra.id, migra.filename)
		}
		if !strings.Contains(migra.contenido, aplicadas[i].detalles) {
			gecko.LogWarnf("Migración " + migra.filename + " no coincide con mensaje aplicado '" + aplicadas[i].detalles + "'")
		}
	}

	// Flush base de datos y volver a abrir.
	err = s.db.Close()
	if err != nil {
		return err
	}
	s.db, err = sql.Open("sqlite", s.dbPath+pragmaConfig)
	if err != nil {
		return err
	}
	return nil
}
