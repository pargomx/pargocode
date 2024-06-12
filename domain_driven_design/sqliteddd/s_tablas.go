package sqliteddd

import (
	"monorepo/domain_driven_design/ddd"
)

func (s *Repositorio) GetTablaByNombre(NombreRepo string) (*ddd.Tabla, error) {
	return s.GetTablaByNombreRepo(NombreRepo)
}
