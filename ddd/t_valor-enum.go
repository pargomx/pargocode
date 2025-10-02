package ddd

import (
	"errors"
	"strings"

	"github.com/pargomx/pargocode/textutils"
)

// ValorEnum corresponde a un elemento de la tabla 'valores_enum'.
type ValorEnum struct {
	CampoID     int    // `valores_enum.campo_id`
	Numero      int    // `valores_enum.numero`
	Clave       string // `valores_enum.clave`
	Etiqueta    string // `valores_enum.etiqueta`
	Descripcion string // `valores_enum.descripcion`
}

var (
	ErrValorEnumNotFound      error = errors.New("el valor enum no se encuentra")
	ErrValorEnumAlreadyExists error = errors.New("el valor enum ya existe")
)

func (val *ValorEnum) Validar() error {

	return nil
}

func (ve ValorEnum) Camel() string {
	return textutils.SnakeToCamel(strings.ToLower(ve.Clave))
}

func (ve ValorEnum) ClaveLower() string {
	return strings.ToLower(ve.Clave)
}
