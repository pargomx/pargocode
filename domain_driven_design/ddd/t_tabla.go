package ddd

import "errors"

// Tabla corresponde a un elemento de la tabla 'tablas'.
type Tabla struct {
	TablaID      int    // `tablas.tabla_id`
	PaqueteID    int    // `tablas.paquete_id`
	NombreRepo   string // `tablas.nombre_repo`  Nombre de la tabla SQL
	NombreItem   string // `tablas.nombre_item`  Nombre de la go struct
	NombreItems  string // `tablas.nombre_items`
	Abrev        string // `tablas.abrev`
	Humano       string // `tablas.humano`
	HumanoPlural string // `tablas.humano_plural`
	Kebab        string // `tablas.kebab`
	EsFemenino   bool   // `tablas.es_femenino`
	Descripcion  string // `tablas.descripcion`
	Directrices  string // `tablas.directrices`
}

var (
	ErrTablaNotFound      error = errors.New("la tabla no se encuentra")
	ErrTablaAlreadyExists error = errors.New("la tabla ya existe")
)

func (tab *Tabla) Validar() error {

	return nil
}
