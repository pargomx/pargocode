package ddd

import "errors"

// Consulta corresponde a un elemento de la tabla 'consultas'.
type Consulta struct {
	ConsultaID  int    // `consultas.consulta_id`
	PaqueteID   int    // `consultas.paquete_id`
	TablaID     int    // `consultas.tabla_id`  Tabla para iniciar el FROM
	NombreItem  string // `consultas.nombre_item`
	NombreItems string // `consultas.nombre_items`
	Abrev       string // `consultas.abrev`
	EsFemenino  bool   // `consultas.es_femenino`
	Descripcion string // `consultas.descripcion`
	Directrices string // `consultas.directrices`
}

var (
	ErrConsultaNotFound      error = errors.New("la consulta no se encuentra")
	ErrConsultaAlreadyExists error = errors.New("la consulta ya existe")
)

func (con *Consulta) Validar() error {

	return nil
}
