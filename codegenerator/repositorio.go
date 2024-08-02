package codegenerator

import "monorepo/ddd"

type Repositorio interface {
	GetPaquete(paqueteID int) (*ddd.Paquete, error)
	ListTablasByPaqueteID(PaqueteID int) ([]ddd.Tabla, error)
	ListConsultasByPaqueteID(PaqueteID int) ([]ddd.Consulta, error)

	GetTabla(tablaID int) (*ddd.Tabla, error)
	ListCamposByTablaID(tablaID int) ([]ddd.Campo, error)
	GetCampo(campoID int) (*ddd.Campo, error)
	GetValoresEnum(CampoID int) ([]ddd.ValorEnum, error)

	GetConsulta(consultaID int) (*ddd.Consulta, error)
	ListConsultaRelacionesByConsultaID(consultaID int) ([]ddd.ConsultaRelacion, error)
	ListConsultaCamposByConsultaID(consultaID int) ([]ddd.ConsultaCampo, error)
}
