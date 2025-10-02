package appdominio

import (
	"errors"

	"github.com/pargomx/pargocode/ddd"
)

type Paquete struct {
	Paquete   ddd.Paquete
	Tablas    []ddd.Tabla
	Consultas []ddd.Consulta

	repo Repositorio
}

func GetPaquete(repo Repositorio, paqueteID int) (Paquete, error) {
	paquete, err := repo.GetPaquete(paqueteID)
	if err != nil {
		return Paquete{}, err
	}
	tablas, err := repo.ListTablasByPaqueteID(paqueteID)
	if err != nil {
		return Paquete{}, err
	}
	consultas, err := repo.ListConsultasByPaqueteID(paqueteID)
	if err != nil {
		return Paquete{}, err
	}
	return Paquete{Paquete: *paquete, Tablas: tablas, Consultas: consultas, repo: repo}, nil
}

// ================================================================ //
// ================================================================ //

func (p *Paquete) BuscarTabla(nombre string) (*ddd.Tabla, error) {
	if nombre == "" {
		return nil, errors.New("buscar tabla sin especificar nombre")
	}
	for _, t := range p.Tablas {
		if nombre == t.NombreRepo || nombre == t.NombreItem {
			return &t, nil
		}
	}
	return nil, ddd.ErrTablaNotFound
}

func (p *Paquete) BuscarConsulta(nombre string) (*ddd.Consulta, error) {
	if nombre == "" {
		return nil, errors.New("buscar consulta sin especificar nombre")
	}
	for _, c := range p.Consultas {
		if nombre == c.NombreItem || nombre == c.NombreItems {
			return &c, nil
		}
	}
	return nil, ddd.ErrConsultaNotFound
}
