package sqlite{{ .TablaOrConsulta.Paquete.Nombre }}

type Repositorio struct {
	db sqlitedb.Ejecutor
}

func NuevoRepo(db sqlitedb.Ejecutor) *Repositorio {
	return &Repositorio{
		db: db,
	}
}