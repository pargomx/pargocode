package appdominio

import (
	"errors"

	"github.com/pargomx/pargocode/ddd"
)

type Tabla struct {
	Tabla   ddd.Tabla
	Campos  []CampoTabla
	Paquete ddd.Paquete

	CamposSeleccionados []CampoTabla
}

// Devuelve una compia del campo si la tabla lo tiene.
// Evalúa nombre Kebab, Snake, Camel, Humano.
func (tbl *Tabla) BuscarCampo(nombre string) (*CampoTabla, error) {
	if nombre == "" {
		return nil, errors.New("no se puede buscar campo sin especificar nombre")
	}
	campo := CampoTabla{}
	for _, c := range tbl.Campos {
		if nombre == c.Campo.NombreCampo ||
			nombre == c.Campo.NombreColumna ||
			nombre == c.Campo.NombreHumano {
			campo = c
			return &campo, nil
		}
	}
	return nil, errors.New("la tabla '" + tbl.Tabla.NombreRepo + "' no tiene campo '" + nombre + "'")
}

// Retorna una copia de todos los campos PK.
func (tbl *Tabla) PrimaryKeys() []CampoTabla {
	res := []CampoTabla{}
	for _, c := range tbl.Campos {
		if c.Campo.PrimaryKey {
			res = append(res, c)
		}
	}
	return res
}

// Retorna una copia de todos los campos FK.
func (tbl *Tabla) ForeignKeys() []CampoTabla {
	res := []CampoTabla{}
	for _, c := range tbl.Campos {
		if c.Campo.ForeignKey {
			res = append(res, c)
		}
	}
	return res
}
