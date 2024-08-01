package ddd

import "errors"

// Paquete corresponde a un elemento de la tabla 'paquetes'.
type Paquete struct {
	PaqueteID   int    // `paquetes.paquete_id`
	GoModule    string // `paquetes.go_module`  Path de importación del proyecto entero
	Directorio  string // `paquetes.directorio`  Directorio relativo desde la raíz del proyecto que contiene el paquete
	Nombre      string // `paquetes.nombre`  Package name
	Descripcion string // `paquetes.descripcion`  Descripción del bounded context
}

var (
	ErrPaqueteNotFound      error = errors.New("el paquete no se encuentra")
	ErrPaqueteAlreadyExists error = errors.New("el paquete ya existe")
)

func (paq *Paquete) Validar() error {

	return nil
}
