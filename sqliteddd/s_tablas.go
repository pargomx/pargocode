package sqliteddd

import (
	"monorepo/ddd"
)

func (s *Repositorio) GetTablaByNombre(NombreRepo string) (*ddd.Tabla, error) {
	return s.GetTablaByNombreRepo(NombreRepo)
}
