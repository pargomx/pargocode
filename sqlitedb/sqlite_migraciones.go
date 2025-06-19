package sqlitedb

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/pargomx/gecko/gko"
)

// ================================================================ //
// ========== INICIALIZAR ========================================= //

// Hay configuraciones PRAGMA que aplican para cada conexión y se ponen en el DSN.
const configPragmaSQL = `PRAGMA journal_mode = WAL;`

const createTableMigraciones = `
CREATE TABLE migraciones (
	major INT NOT NULL,
	minor INT NOT NULL,
	fecha TEXT NOT NULL,
	detalles TEXT NOT NULL,
	PRIMARY KEY (major,minor),
	UNIQUE(detalles)
);`

// Inicializa la base de datos con el configPragmaSQL y crea la tabla para
// migraciones. Confía en que el archivo en dbPath existe y está limpio.
func (s *SqliteDB) initDatabase() error {
	op := gko.Op("initDB")
	// Configurar database debe ser fuera de transacciones.
	_, err := s.Exec(configPragmaSQL)
	if err != nil {
		return op.Err(err)
	}
	// Crear tabla para automatizar migraciones.
	_, err = s.ExecInTransaction(createTableMigraciones)
	if err != nil {
		return op.Err(err)
	}
	gko.LogInfof("SQLiteDB: nuevo archivo %v configurado", s.dbPath)
	return nil
}

// ================================================================ //
// ========== MIGRACIONES ========================================= //

type migracionAplicada struct {
	major    int    // `migraciones.major`
	minor    int    // `migraciones.minor`
	fecha    string // `migraciones.fecha`
	detalles string // `migraciones.detalles`
}

type migracionDisponible struct {
	major     int
	minor     int
	filename  string
	contenido string // El contenido del script
}

const selectMigraciones = "SELECT major, minor, fecha, detalles FROM migraciones"

// LEGACY: Actualizar del viejo método de automatizar migraciones.
const oldSelectMigraciones = "SELECT id, fecha, detalles FROM migraciones"
const updateTablaMigraciones = `
ALTER TABLE migraciones RENAME TO migraciones_old;

CREATE TABLE migraciones (
	major INT NOT NULL,
	minor INT NOT NULL,
	fecha TEXT NOT NULL,
	detalles TEXT NOT NULL,
	PRIMARY KEY (major,minor),
	UNIQUE(detalles)
);

INSERT INTO migraciones(major, minor, fecha, detalles)
  SELECT 1, id, fecha, detalles FROM migraciones_old;

DROP TABLE migraciones_old;
`

// Inicializa o actualiza la base de datos:
// Verifica que las migraciones estén aplicadas y las aplica si no lo están.
//
// Deben estar en el directorio "migraciones" y tener de prefijo un número
// consecutivo seguido de un guión bajo. Ejemplo: "v1/03_usuarios.sql".
//
// Toda migración debe registrarse a sí misma en la tabla "migraciones" con
// el mismo número de versión major que su directorio y minor que tiene en su prefijo.
func (s *SqliteDB) verificarMigraciones(migracionesFS fs.FS) error {
	op := gko.Op("sqlitedb.Migraciones")

	// LEGACY: Puede tener el sistema anterior para migraciones.
	_, err := s.db.Query(oldSelectMigraciones)
	if err == nil { // Error nulo es que sí tiene el esquema anterior.
		_, err = s.ExecInTransaction(updateTablaMigraciones)
		if err != nil {
			return op.Err(err).Op("UpdateTablaMigraciones")
		}
		gko.LogEvento("SQLiteDB: migraciones mejoradas")
	}

	// Conocer todas las migraciones aplicadas.
	aplicadas := make(map[[2]int]migracionAplicada) // key: major
	minVersionAplicada := 0

	// Puede no haya tabla migraciones, señal que se debe inicializar db.
	rows, err := s.db.Query(selectMigraciones)
	if err != nil {
		err = s.initDatabase()
		if err != nil {
			return op.Err(err)
		}

	} else {
		// Si la tabla migraciones existe, traer todos los registros.
		for rows.Next() {
			apli := migracionAplicada{}
			err := rows.Scan(&apli.major, &apli.minor, &apli.fecha, &apli.detalles)
			if err != nil {
				return op.Err(err).Op("GetMigracionesAplicadas")
			}
			aplicadas[[2]int{apli.major, apli.minor}] = apli
		}
		// Conocer la primer versión major aplicada para no aplicar migraciones anteriores innecesarias.
		err = s.db.QueryRow("SELECT coalesce(min(major),0) FROM migraciones").Scan(&minVersionAplicada)
		if err != nil { // Innecesario && !errors.Is(err, sql.ErrNoRows) por coalesce.
			return op.Err(err).Str("can't get first major version applied")
		}
	}

	// Obtener migraciones disponibles en orden (ej. 1.0, 1.1, 1.2, 1.3, 2.0, 2.1)
	// para comprobar y aplicar las que hagan falta.
	disponibles, err := s.getMigracionesDisponibles(migracionesFS)
	if err != nil {
		return op.Err(err)
	}
	// Para no aplicar migraciones viejas innecesarias a instancias nuevas se necesita saber:
	maxVersionDisponible := getLastMajorDisponible(disponibles)
	ceroAplicadas := len(aplicadas) == 0 // si la base de datos está vacía porque es una instancia nueva.

	// para verificar que sean consecutivas.
	majorLastLoop, minorLastLoop := 0, 0

	// para migrar datos entre versiones mayores.
	migracionDatosForUpgrade := migracionDisponible{}

	// Verificar cada una de las migraciones disponibles que apliquen.
	for _, dispo := range disponibles {

		if dispo.major < 1 {
			return op.Msg("Migración inválida").Strf("deben comenzar en v1/... no %v", dispo.filename)
		}
		if dispo.minor < 0 {
			return op.Msg("Migración inválida").Strf("no números negativos: %v", dispo.filename)
		}
		if dispo.major == 1 && dispo.minor == 0 {
			return op.Msg("Migración inválida").
				Strf("ignorando migración %v porque vX/0 es para migrar datos desde versión anterior y v1/0 es inválido", dispo.filename)
		}

		// Verificar que sean consecutivas.
		if majorLastLoop+1 == dispo.major {
			majorLastLoop = dispo.major // puede subir de versión
			minorLastLoop = 0
			if majorLastLoop > 1 {
				minorLastLoop = -1 // la v2+ empeiza en 0.
			}
		} else if majorLastLoop != dispo.major { // puede seguir en esta versión
			return op.Msg("Migración inválida").Strf("migración mayor no consecutiva: %v", dispo.filename)
		}
		if minorLastLoop+1 == dispo.minor {
			minorLastLoop = dispo.minor
		} else {
			return op.Msg("Migración inválida").Strf("migración menor no consecutiva: %v", dispo.filename)
		}

		// La migración X.0 se guarda para ejecutarse solo
		// después de hacer upgrade mayor hacia X.1 (esquema de datos).
		if dispo.minor == 0 {
			migracionDatosForUpgrade = dispo
			continue
		}

		// Cuando haya una versión mayor superior para una nueva instancia,
		// saltar versiones anteriores porque no hay datos que migrar.
		// If emptyDB AND lastMigra is v3.x, THEN skip v1.x v2.x v3.0
		if ceroAplicadas && dispo.major < maxVersionDisponible {
			continue
		}

		// Ignorar versiones major anteriores a la aplicada al crear esta instancia de la app.
		// If aplicada v2.x v3.1 v3.2 AND disponibles v3.3 v4.1 v4.2 THEN skip v1.x
		if minVersionAplicada > dispo.major {
			continue
		}

		// Identificar la migración que se va a comprobar.
		key := [2]int{dispo.major, dispo.minor}

		// Si ya está aplicada no hacer nada y saltar a la siguiente.
		if _, aplicada := aplicadas[key]; aplicada {
			continue
		}

		// Casos de las migraciones:
		// [1.0] se ignora
		// [1.1+] no necesitan upgrade porque vienen de archivo.db en blanco.
		// [2+.1] requiere upgrade y aplicar [2+.0] en la misma transacción para pasar los datos.
		//        a menos que sea empty, entonces solo se aplica [2.1] sin upgrade.
		// [2+.3+] no nencesitan upgrade.
		if dispo.major >= 2 && dispo.minor == 1 && !ceroAplicadas {
			err = s.upgradeVersionEsquema(dispo, migracionDatosForUpgrade, aplicadas)
			if err != nil {
				return op.Err(err)
			}
		} else if dispo.minor > 0 {
			err = s.aplicarMigraciones(aplicadas, dispo)
			if err != nil {
				return op.Err(err)
			}
		} else {
			return op.Strf("SQLiteDB: ignorando migración %v", dispo.filename)
		}

		// Verificar que coincida la migración aplicada con la disponible.
		aplicada, ok := aplicadas[key]
		if !ok {
			return op.Msg("Revisar DB").Strf("migración aplicada '%v' no se registró en migraciones como (%v,%v) y se aplicará otra vez en reinicio", dispo.filename, dispo.major, dispo.minor)
		}
		if !strings.Contains(dispo.contenido, aplicada.detalles) {
			// only warning to allow manual overrides.
			gko.LogWarnf("REVISAR DB: migración '%v' no coincide con mensaje aplicado '%v'", dispo.filename, aplicada.detalles)
		}
	}

	// Reconectar db para hacer flush luego de migraciones.
	err = s.Close()
	if err != nil {
		return op.Err(err)
	}
	err = s.openDatabase()
	if err != nil {
		return op.Err(err)
	}
	return nil
}

// ================================================================ //

// Crea un nuevo archivo de base de datos con el esquema desde cero y
// luego pasa los datos desde la versión anterior (si los hay).
//
// La migración de esquema debe ser X.1 y la de datos X.0, donde X > 1.
func (s *SqliteDB) upgradeVersionEsquema(migEsquema migracionDisponible, migDatos migracionDisponible, aplicadas map[[2]int]migracionAplicada) error {
	op := gko.Op("UpgradeEsquema")

	// Para la primera migración solo hay que inicializar, no hay que migrar.
	if migEsquema.major < 2 {
		return op.Strf("migración de esquema solo aplica a partir de v2, no %v.%v", migEsquema.major, migEsquema.minor)
	}
	if migEsquema.minor != 1 {
		return op.Strf("migración de esquema inválida: %v.%v", migEsquema.major, migEsquema.minor)
	}
	if migEsquema.major != migDatos.major || migDatos.minor != 0 {
		return op.Strf("migración de esquema %v.%v requiere migración de datos %v.0, no %v.%v",
			migEsquema.major, migEsquema.minor, migEsquema.major, migDatos.major, migDatos.minor)
	}

	err := s.Backup() // TODO: muchas migraciones conflicto en el mismo segundo? agregar sufijo con migracion_id mejor
	if err != nil {
		return op.Err(err)
	}
	err = s.CloseFully() // asegurando que todo se contenga en un solo archivo.
	if err != nil {
		return op.Err(err)
	}

	// Comprobar que no exista otro archivo con el mismo nombre para la nueva db.
	newTempFilename := "newDatabaseFile.db"
	newTempPath := path.Join(path.Dir(s.dbPath), newTempFilename)
	if _, err := os.Stat(newTempPath); err == nil {
		return op.Strf("new temp db file ya existe: %v", newTempPath)
	} else if !os.IsNotExist(err) {
		return op.Err(err).Strf("new temp db file ya existe? %v", newTempPath)
	}

	// Al final siempre se debe volver la conexión a la ruta original.
	oldDatabasePath := s.dbPath
	defer func() {
		s.db.Close()
		s.dbPath = oldDatabasePath
		s.openDatabase()
		os.Remove(newTempPath)
	}()

	// Cambiar el dbPath hacia un nuevo archivo vacío y conectarse.
	s.dbPath = newTempPath
	err = s.openDatabase()
	if err != nil {
		return op.Err(err)
	}

	err = s.initDatabase()
	if err != nil {
		return op.Err(err)
	}

	// Hacer el attach a la base de datos vieja para poder pasar los datos a la nueva.
	_, err = s.db.Exec("ATTACH DATABASE ? AS old_schema", oldDatabasePath)
	if err != nil {
		return op.Err(err)
	}

	// Pasar registro de migraciones aplicadas anteriormente.
	_, err = s.db.Exec("INSERT INTO main.migraciones SELECT * FROM old_schema.migraciones")
	if err != nil {
		return op.Err(err)
	}

	// Aplicar migración con el esquema nuevo desde cero y luego migrar datos del esquema viejo.
	err = s.aplicarMigraciones(aplicadas, migEsquema, migDatos)
	if err != nil {
		return op.Err(err)
	}

	// Cerrar para poder reemplazar viejo archivo por el nuevo recién migrado.
	err = s.CloseFully()
	if err != nil {
		return op.Err(err).Msg("Cleanup required").
			Strf("se aplicaron las migraciones a '%v' pero no se pudo reemplazar la db", newTempFilename)
	}
	// dbPath es el newTempFile, entonces verificar manualmente que esté cerrada tambien el old.
	if _, err := os.Stat(oldDatabasePath + "-wal"); err == nil {
		return op.Strf("old db still open: WAL file exists (%v)", oldDatabasePath+"-wal")
	}
	if _, err := os.Stat(oldDatabasePath + "-shm"); err == nil {
		return op.Strf("old db still open: SHM file exists (%v)", oldDatabasePath+"-shm")
	}

	// Conservar original como extra backup.
	err = os.Rename(oldDatabasePath, fmt.Sprintf(oldDatabasePath+".v%v.bak", migEsquema.major-1))
	if err != nil {
		return op.Err(err)
	}

	// Poner el nuevo en su lugar.
	err = os.Rename(newTempFilename, oldDatabasePath)
	if err != nil {
		return op.Err(err)
	}
	s.dbPath = oldDatabasePath

	err = s.openDatabase()
	if err != nil {
		return op.Err(err)
	}
	return nil
}

// Aplica una o varias migraciones dentro de una sola transacción.
// Cualquier error hace un rollback completo (excepto para el commit).
func (s *SqliteDB) aplicarMigraciones(aplicadas map[[2]int]migracionAplicada, porAplicar ...migracionDisponible) error {
	op := gko.Op("Aplicar")
	if len(porAplicar) == 0 {
		return op.Str("nada por aplicar")
	}
	tx, err := s.db.Begin()
	if err != nil {
		return op.Err(err).Strf("no se puede aplicar migración '%v'", porAplicar[0].filename)
	}
	outputMsgs := []string{}
	defer func() {
		if len(outputMsgs) > 0 {
			gko.LogInfof("SQLiteDB: mensajes desde migración:\n %v", strings.Join(outputMsgs, "\n "))
		}
	}()
	for _, migra := range porAplicar {
		// Ejecutar cada statement por separado para facilitar debug porque sqlite no da info.
		statements := strings.Split(migra.contenido, ";")
		for i, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			// Al primer error dentro de este loop se provoca un rollback con esta info.
			execErr := gko.Op("Aplicar").Msgf("Rollback %v por fallo en sentencia %v", migra.filename, i+1)
			if len(stmt) > 240 {
				execErr.Str(stmt[:240] + "...")
			} else {
				execErr.Str(stmt)
			}

			// Para mensajes debug se recibe una sola columna de texto.
			if strings.HasPrefix(strings.ToUpper(stmt), "SELECT") {
				rows, err := tx.Query(stmt)
				if err != nil {
					tx.Rollback()
					return execErr.Err(err)
				}
				defer rows.Close()
				for rows.Next() {
					var msg string
					err := rows.Scan(&msg)
					if err != nil {
						return execErr.Err(err)
					}
					outputMsgs = append(outputMsgs, fmt.Sprintf("%v(%02d): %v", migra.filename, i+1, msg))
				}

			} else {
				_, err := tx.Exec(stmt)
				if err != nil {
					tx.Rollback()
					return execErr.Err(err)
				}
			}
		}

		// Comprobar que se haya registrado la aplicación en la tabla de migraciones.
		aplicada := migracionAplicada{}
		migRegErr := gko.Op("Aplicar").Msg("Migración fallida (rollback)").Str(migra.filename)
		err = tx.QueryRow(selectMigraciones+" WHERE major = ? AND minor = ?",
			migra.major, migra.minor).
			Scan(&aplicada.major, &aplicada.minor, &aplicada.fecha, &aplicada.detalles)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				return migRegErr.Strf("se debería registrar como %v.%v en db", migra.major, migra.minor)
			}
			return migRegErr.Err(err)
		}
		if aplicada.major != migra.major || aplicada.minor != migra.minor {
			tx.Rollback()
			return migRegErr.Strf("se registra como %v.%v en db pero debería ser %v.%v",
				aplicada.major, aplicada.minor, migra.major, migra.minor)
		}
		aplicadas[[2]int{migra.major, migra.minor}] = aplicada
		gko.LogEventof("SqliteDB: migración aplicada %v", migra.filename)

	}
	err = tx.Commit()
	if err != nil {
		return op.Err(err).Strf("failed commit for '%v'", porAplicar[0].filename)
	}
	return nil
}

// Obtiene todas las migraciones en orden (minor,major) y lee su contenido.
//
//   - Deben estar contenidas en directorios v1, v2, v3.
//   - Deben comenzar por un número, guion bajo, y terminar con .sql
//   - La migración 0 de una v2 o superior es para migrar los datos
//     desde la última migración de la versión anterior.
//   - Los arcivos que comienzan por "_guion.sql" se ignoran.
//
// Ejemplos válidos:
//
//   - v1/1_startSchema.sql
//
//   - v1/1_addTable.sql
//
//   - v2/0_dataTransform.sql
//
//   - v2/1_fullSchema.sql
//
//   - v2/2_addColumn.sql
//
// Se pasa como argumento para el colector de basura.
func (s *SqliteDB) getMigracionesDisponibles(migracionesFS fs.FS) ([]migracionDisponible, error) {
	op := gko.Op("getMigracionesMajorDirs")
	disponibles := []migracionDisponible{}
	majorDirs, err := fs.ReadDir(migracionesFS, ".")
	if err != nil {
		return nil, op.Err(err)
	}
	for _, majorDir := range majorDirs {
		if !majorDir.IsDir() || !strings.HasPrefix(majorDir.Name(), "v") {
			continue // Only process directories like v1, v2, ...
		}
		major, err := strconv.Atoi(strings.TrimPrefix(majorDir.Name(), "v"))
		if err != nil {
			return nil, op.Err(err).Op("parseMajorVersion")
		}
		files, err := fs.ReadDir(migracionesFS, majorDir.Name())
		if err != nil {
			return nil, op.Err(err).Op("getMigracionesFiles")
		}
		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
				continue // Ignore non-sql files and directories
			}
			if strings.HasPrefix(file.Name(), "_") {
				continue // Ignore files starting with underscore
			}
			minorStr := strings.Split(file.Name(), "_")[0]
			minor, err := strconv.Atoi(minorStr)
			if err != nil {
				return nil, op.Err(err).Op("parseMinorVersion")
			}
			path := majorDir.Name() + "/" + file.Name()
			bytes, err := fs.ReadFile(migracionesFS, path)
			if err != nil {
				return nil, op.Err(err).Op("readFile")
			}
			disponibles = append(disponibles, migracionDisponible{
				major:     major,
				minor:     minor,
				filename:  path,
				contenido: string(bytes),
			})
		}
	}
	// Sort disponibles by major, then minor
	sort.Slice(disponibles, func(i, j int) bool {
		if disponibles[i].major == disponibles[j].major {
			return disponibles[i].minor < disponibles[j].minor
		}
		return disponibles[i].major < disponibles[j].major
	})
	return disponibles, nil
}

// Obtener la última migración mayor disponible.
func getLastMajorDisponible(disponibes []migracionDisponible) int {
	max := 0
	for _, migra := range disponibes {
		if migra.major > max {
			max = migra.major
		}
	}
	return max
}
