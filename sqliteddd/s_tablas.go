package sqliteddd

import (
	"github.com/pargomx/pargocode/ddd"
)

func (s *Repositorio) GetTablaByNombre(NombreRepo string) (*ddd.Tabla, error) {
	return s.GetTablaByNombreRepo(NombreRepo)
}
