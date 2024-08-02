package codegenerator

import (
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
