package sqliteddd

import (
	"monorepo/sqlitedb"
)

// ================================================================ //
// ========== Repositorio ========================================= //

type Repositorio struct {
	db sqlitedb.Ejecutor
}

// Se puede pasar la DB directamente, o bien una transacción ya iniciada.
// Es responsabilidad de quien invoca hacer Rollback o Commit.
func NuevoRepositorio(db sqlitedb.Ejecutor) *Repositorio {
	return &Repositorio{
		db: db,
	}
}
