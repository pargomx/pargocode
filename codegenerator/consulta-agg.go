package codegenerator

import (
	"monorepo/ddd"
)

type Relacion struct {
	ConsultaID  int
	Posicion    int
	TipoJoin    ddd.TipoJoin
	JoinTablaID int
	JoinAs      string
	JoinOn      string
	FromTablaID int

	Join tabla
	From *tabla
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

	Consulta      *consulta
	OrigenTabla   *ddd.Tabla
	OrigenPaquete *ddd.Paquete
	OrigenCampo   *ddd.Campo
}
