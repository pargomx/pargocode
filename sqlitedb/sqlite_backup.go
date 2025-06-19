package sqlitedb

import (
	"fmt"
	"os"
	"path"

	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

func (s *SqliteDB) Backup() error {
	op := gko.Op("sqlitedb.Backup")

	// Directorio para backups.
	if s.backupsDir == "" {
		s.backupsDir = "backups"
	}
	info, err := os.Stat(s.backupsDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(s.backupsDir, 0750)
		if err != nil {
			return op.Err(err).Op("NewDatabaseDir")
		}
		gko.LogInfof("SQLiteDB: directorio creado '%v'", s.backupsDir)
	} else if err != nil {
		return op.Err(err)
	} else if !info.IsDir() {
		return op.Str("Directorio para backups inválido: %v")
	}

	// Destino para el backup.
	backupName := fmt.Sprintf("%v.%v.db", path.Base(s.dbPath), gkt.Now().Format("2006-01-02_150405"))
	backupPath := path.Join(s.backupsDir, backupName)

	// Comprobar que no exista otro archivo con el mismo nombre para el backup.
	if _, err := os.Stat(backupPath); err == nil {
		return op.Strf("conflicto: backup file ya existe: %v", backupPath)
	} else if !os.IsNotExist(err) {
		return op.Err(err).Strf("backup file ya existe? %v", backupPath)
	}

	// Cerrar base de datos para que todo esté contenido en un solo archivo.
	err = s.CloseFully()
	if err != nil {
		return op.Err(err)
	}

	// Copiar archivo de base de datos.
	dstFile, err := os.OpenFile(backupPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return op.Err(err)
	}
	defer dstFile.Close()

	srcFile, err := os.Open(s.dbPath)
	if err != nil {
		return op.Err(err)
	}
	defer srcFile.Close()

	_, err = srcFile.Seek(0, 0)
	if err != nil {
		return op.Err(err)
	}

	_, err = dstFile.ReadFrom(srcFile)
	if err != nil {
		return op.Err(err)
	}

	gko.LogInfof("SqliteDB: backup saved '%v'", backupPath)

	// Volver a abrir db original.
	err = s.openDatabase()
	if err != nil {
		return op.Err(err)
	}
	return nil
}
