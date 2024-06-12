package mysql{{ .TablaOrConsulta.Paquete.Nombre }}

// ================================================================ //
// ========== Repositorio ========================================= //

type Repositorio struct {
	db mysqldb.Ejecutor
}

func NuevoRepoMySQL(db mysqldb.Ejecutor) *Repositorio {
	return &Repositorio{
		db: db,
	}
}