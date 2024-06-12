package gecko

import (
	"errors"
	"fmt"
	"net/http"
)

var MostrarMensajeEnErrores bool = false

type Gkerror struct {
	// HTTP status code que define el tipo de error
	codigo int

	// Mensaje para el usuario
	mensaje string

	// Operación que se estaba intentando realizar
	operación string

	// Contexto de la acción a realizar. Ej. "editar usuario_id=1234"
	contexto string

	// Error genérico de una dependencia externa
	err error
}

// Error satisface la interfaz `error` componiendo el mensaje
// de una manera comprensible y completa para poner en los logs.
//
// Evitar visibilizar al usuario porque da todo el contexto.
func (e *Gkerror) Error() string {
	msg := ""
	if e.codigo > 0 {
		msg += fmt.Sprintf("[%d]", e.codigo)
	}
	if e.operación != "" {
		msg += " " + e.operación
	}
	if e.err != nil {
		msg += " " + e.err.Error()
	}
	if e.mensaje != "" && MostrarMensajeEnErrores {
		msg += ": " + e.mensaje + "."
	}
	if e.contexto != "" {
		msg += " {" + e.contexto + "}"
	}
	return msg
}

// ================================================================ //
// ========== S E T T E R S ======================================= //

// Define un nuevo status code para el error.
//
// Subsecuentes llamadas sustituyen el código anterior.
func (e *Gkerror) Code(code int) *Gkerror {
	e.codigo = code
	return e
}

func (e *Gkerror) GetCode() int {
	return e.codigo
}

// Mensaje dirigido al usuario.
// Subsecuentes llamadas se concatenan con `: `.
func (e *Gkerror) Msg(msg string) *Gkerror {
	if e.mensaje == "" {
		e.mensaje = msg
	} else {
		e.mensaje += ": " + msg
	}
	return e
}

// Mensaje dirigido al usuario.
// Subsecuentes llamadas se concatenan con `: `.
func (e *Gkerror) Msgf(format string, a ...any) *Gkerror {
	if e.mensaje == "" {
		e.mensaje = fmt.Sprintf(format, a...)
	} else {
		e.mensaje += ": " + fmt.Sprintf(format, a...)
	}
	return e
}

// Operación que se intentaba realizar.
// Subsecuentes llamadas se concatenan con ` > `.
func (e *Gkerror) Op(op string) *Gkerror {
	if e.operación == "" {
		e.operación = op
	} else {
		e.operación = op + " > " + e.operación
	}
	return e
}

// Contexto en forma de "clave=valor".
// Subsecuentes llamadas se concatenan con ` `.
func (e *Gkerror) Ctx(key string, val any) *Gkerror {
	if e.contexto == "" {
		e.contexto = fmt.Sprintf("%s=%v", key, val)
	} else {
		e.contexto += fmt.Sprintf(" %s=%v", key, val)
	}
	return e
}

func (e *Gkerror) Err(err error) *Gkerror {
	if err == nil {
		if e.err == nil {
			e.err = errors.New("err nil")
		} else {
			e.err = errors.Join(e.err, errors.New("err nil"))
		}
		return e
	}

	ne, esGecko := err.(*Gkerror)
	if !esGecko {
		if e.err == nil {
			e.err = err
		} else {
			e.err = errors.Join(e.err, err)
		}
		return e
	}

	// es gecko
	e.codigo = ne.codigo

	if ne.mensaje != "" {
		if ne.mensaje == "" {
			e.mensaje = ne.mensaje
		} else {
			e.mensaje += ": " + ne.mensaje
		}
	}

	if ne.operación != "" {
		if e.operación == "" {
			e.operación = ne.operación
		} else {
			e.operación = ne.operación + " > " + e.operación
		}
	}

	if ne.contexto != "" {
		if e.contexto == "" {
			e.contexto = ne.contexto
		} else {
			e.contexto += " " + ne.contexto
		}
	}

	if ne.err != nil {
		if e.err == nil {
			e.err = ne.err
		} else {
			e.err = errors.Join(e.err, ne.err)
		}
		return e
	}

	// fmt.Println("gk: " + err.Error())

	return e
}

// ================================================================ //
// ========== C O N S T R U C T O R E S =========================== //

// Nuevo error gecko con http status code.
func NewErr(code int) *Gkerror {
	return &Gkerror{
		codigo: code,
	}
}

// Nuevo error gecko con la operación que se intenta realizar.
func NewOp(op string) *Gkerror {
	return &Gkerror{
		operación: op,
	}
}

// ================================================================ //
// ========== A S S E R T I O N S ================================= //

// Convierte un interface error en *GeckoError usando type assertion.
func Err(err error) *Gkerror {
	if err == nil {
		return &Gkerror{
			err: errors.New("err nil"),
		}
	}
	if errGecko, ok := err.(*Gkerror); ok {
		return errGecko
	}
	return &Gkerror{
		err: err,
	}
}

// Reporta si el código del error es 404 NotFound.
func (e *Gkerror) EsNotFound() bool {
	if e == nil {
		return false
	}
	return e.codigo == http.StatusNotFound
}

// Reporta si el código del error es 409 Conflict.
func (e *Gkerror) EsAlreadyExists() bool {
	if e == nil {
		return false
	}
	return e.codigo == http.StatusConflict
}

// Reporta si el código del error concreto es 404 NotFound.
func EsErrNotFound(err error) bool {
	if err == nil {
		return false
	}
	e, ok := err.(*Gkerror)
	if !ok {
		return false
	}
	if e == nil {
		return false
	}
	return e.codigo == http.StatusNotFound
}

// Reporta si el código del error concreto es 409 Conflict.
func EsErrAlreadyExists(err error) bool {
	if err == nil {
		return false
	}
	e, ok := err.(*Gkerror)
	if !ok {
		return false
	}
	if e == nil {
		return false
	}
	return e.codigo == http.StatusConflict
}
