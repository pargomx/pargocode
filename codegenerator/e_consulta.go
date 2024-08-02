package codegenerator

import (
	"fmt"
	"monorepo/ddd"
)

type Consulta struct {
	Paquete     ddd.Paquete
	Consulta    ddd.Consulta
	TablaOrigen ddd.Tabla
	From        Tabla
	Campos      []CampoConsulta
	Relaciones  []Relacion

	CamposSeleccionados []CampoConsulta // Multiuso. Default vacío.
}

type Relacion struct {
	ConsultaID  int
	Posicion    int
	TipoJoin    ddd.TipoJoin
	JoinTablaID int
	JoinAs      string
	JoinOn      string
	FromTablaID int

	Join Tabla
	From *Tabla
}

type CampoConsulta struct {
	ConsultaID  int
	Posicion    int
	CampoID     *int
	Expresion   string
	AliasSql    string
	NombreCampo string
	TipoGo      string
	Pk          bool
	Filtro      bool
	GroupBy     bool
	Descripcion string

	Consulta      *Consulta
	OrigenTabla   *ddd.Tabla
	OrigenPaquete *ddd.Paquete
	OrigenCampo   *ddd.Campo
}

// ================================================================ //
// ================================================================ //

func (con *Consulta) NombreItem() string {
	return con.Consulta.NombreItem
}
func (con *Consulta) NombreItems() string {
	return con.Consulta.NombreItems
}
func (con *Consulta) NombreAbrev() string {
	return con.Consulta.Abrev
}

func (con *Consulta) Directrices() []Directriz {
	return ToDirectrices(con.Consulta.Directrices)
}

func (con *Consulta) BuscarCampo(nombre string) (*CampoConsulta, error) {
	if nombre == "" {
		return nil, fmt.Errorf("nombre de campo a buscar en la consulta vacío")
	}
	for _, campo := range con.Campos {
		if campo.NombreCampo == nombre {
			return &campo, nil
		}
		if campo.AliasSql == nombre {
			return &campo, nil
		}
		if campo.Expresion == nombre {
			return &campo, nil
		}
	}
	return nil, fmt.Errorf("campo no encontrado: %v", nombre)
}

func (con *Consulta) TieneCamposFiltro() bool {
	for _, campo := range con.Campos {
		if campo.Filtro {
			return true
		}
	}
	return false
}

func (con *Consulta) CamposFiltro() []CampoConsulta {
	campos := []CampoConsulta{}
	for _, campo := range con.Campos {
		if campo.Filtro {
			campos = append(campos, campo)
		}
	}
	return campos
}

func (con *Consulta) TieneCamposGroupBy() bool {
	for _, campo := range con.Campos {
		if campo.GroupBy {
			return true
		}
	}
	return false
}

func (con *Consulta) CamposGroupBy() []CampoConsulta {
	campos := []CampoConsulta{}
	for _, campo := range con.Campos {
		if campo.GroupBy {
			campos = append(campos, campo)
		}
	}
	return campos
}

func (con *Consulta) PrimaryKeys() []CampoConsulta {
	campos := []CampoConsulta{}
	for _, campo := range con.Campos {
		if campo.Pk {
			campos = append(campos, campo)
		}
	}
	return campos
}

// ================================================================ //
// ================================================================ //
