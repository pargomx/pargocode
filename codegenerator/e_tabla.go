package codegenerator

import (
	"errors"
	"monorepo/ddd"
	"strings"
)

type Tabla struct {
	Tabla   ddd.Tabla
	Campos  []CampoTabla
	Paquete ddd.Paquete

	CamposSeleccionados []CampoTabla
}

// ================================================================ //
// ================================================================ //

func (tbl *Tabla) NombreRepo() string {
	return tbl.Tabla.NombreRepo
}
func (tbl *Tabla) NombreItem() string {
	return tbl.Tabla.NombreItem
}
func (tbl *Tabla) NombreItems() string {
	return tbl.Tabla.NombreItems
}
func (tbl *Tabla) NombreAbrev() string {
	return tbl.Tabla.Abrev
}
func (tbl *Tabla) NombreNominativo() string {
	if tbl.Tabla.EsFemenino {
		return "la " + strings.ToLower(tbl.Tabla.Humano)
	}
	return "el " + strings.ToLower(tbl.Tabla.Humano)
}

func (tbl *Tabla) Directrices() []Directriz {
	return ToDirectrices(tbl.Tabla.Directrices)
}

func (tbl *Tabla) TablaOrigen() ddd.Tabla {
	return tbl.Tabla
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

func (tbl *Tabla) PrimerCampo() CampoTabla {
	return tbl.Campos[0]
}

func (tbl *Tabla) TieneCamposFiltro() bool {
	for _, campo := range tbl.Campos {
		if campo.Campo.Filtro {
			return true
		}
	}
	return false
}

func (tbl *Tabla) CamposFiltro() []CampoTabla {
	campos := []CampoTabla{}
	for _, campo := range tbl.Campos {
		if campo.Campo.Filtro {
			campos = append(campos, campo)
		}
	}
	return campos
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

func (tbl *Tabla) UniqueKeys() []CampoTabla {
	res := []CampoTabla{}
	for _, c := range tbl.Campos {
		if c.Campo.Unique() {
			res = append(res, c)
		}
	}
	return res
}

func (tbl *Tabla) CamposEspeciales() []CampoTabla {
	res := []CampoTabla{}
	for _, c := range tbl.Campos {
		if c.Campo.Especial {
			res = append(res, c)
		}
	}
	return res
}

func (tbl *Tabla) CamposEditables() []CampoTabla {
	res := []CampoTabla{}
	for _, c := range tbl.Campos {
		if !c.Campo.ReadOnly() {
			res = append(res, c)
		}
	}
	return res
}

// Retorna los campos que son requeridos para escribir (required y PKs).
func (tbl *Tabla) CamposRequeridosOrPK() []CampoTabla {
	res := []CampoTabla{}
	for _, c := range tbl.Campos {
		if c.Required() || c.PrimaryKey {
			res = append(res, c)
		}
	}
	return res
}
