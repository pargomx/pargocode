package appdominio

import "monorepo/ddd"

type Repositorio interface {
	GetPaquete(paqueteID int) (*ddd.Paquete, error)
	ExistePaquete(paqueteID int, nombre string) bool
	InsertPaquete(paq ddd.Paquete) error
	UpdatePaquete(paq ddd.Paquete) error
	DeletePaquete(paqueteID int) error

	InsertTabla(tab ddd.Tabla) error
	UpdateTabla(tab ddd.Tabla) error

	GetCampo(campoID int) (*ddd.Campo, error)
	GetCampoPrimaryKey(nombre string) (*ddd.Campo, error)
	GetCampoByNombre(nombre string) (*ddd.Campo, error)
	InsertCampo(campo ddd.Campo) error
	UpdateCampo(campo ddd.Campo) error
	ReordenarCampo(cam *ddd.Campo, newPosicion int) error
	GuardarValoresEnum(campoID int, vals []ddd.ValorEnum) error
	DeleteCampo(campoID int) error
	GetValoresEnum(CampoID int) ([]ddd.ValorEnum, error)

	ListTablas() ([]ddd.Tabla, error)
	ListTablasByPaqueteID(PaqueteID int) ([]ddd.Tabla, error)
	ListConsultasByPaqueteID(PaqueteID int) ([]ddd.Consulta, error)
	GetTablaByNombre(nombre string) (*ddd.Tabla, error)
	GetTabla(tablaID int) (*ddd.Tabla, error)
	ListCamposByTablaID(tablaID int) ([]ddd.Campo, error)

	GetConsulta(consultaID int) (*ddd.Consulta, error)
	InsertConsulta(consulta ddd.Consulta) error
	UpdateConsulta(consulta ddd.Consulta) error
	DeleteConsulta(consultaID int) error

	InsertConsultaRelacion(rel ddd.ConsultaRelacion) error
	DeleteRelacionConsulta(ConsultaID int, Posicion int) error
	UpdateConsultaRelacion(rel ddd.ConsultaRelacion) error
	ListConsultaRelacionesByConsultaID(consultaID int) ([]ddd.ConsultaRelacion, error)

	InsertConsultaCampo(cam ddd.ConsultaCampo) error
	DeleteConsultaCampo(ConsultaID int, Posicion int) error
	UpdateConsultaCampo(cam ddd.ConsultaCampo) error
	GetConsultaCampo(consultaID, posicion int) (*ddd.ConsultaCampo, error)
	ReordenarCampoConsulta(consultaID int, oldPosicion, newPosicion int) error
	ListConsultaCamposByConsultaID(consultaID int) ([]ddd.ConsultaCampo, error)
}
