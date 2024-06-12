package textutils

import (
	"errors"
	"regexp"
	"strings"

	"monorepo/gecko"

	"github.com/manifoldco/promptui"
)

// NuevaValidación devuelve una función de validación
// que usa la expresión regular compilada y devuelve
// el mensaje de error dado si el input no la cumple.
func NuevaValidación(reg *regexp.Regexp, errMsg string) promptui.ValidateFunc {
	if reg == nil {
		gecko.FatalFmt("NuevaValidación: regexp es nil")
	}
	return func(input string) error {
		if input == "" {
			return errors.New("input requerido")
		}
		if len(input) < 3 {
			return errors.New("debe tener al menos 3 caracteres")
		}
		if !reg.MatchString(input) {
			return errors.New(errMsg)
		}
		return nil
	}
}

// Acepta espacios, letras y acentos.
func RegexLetrasAcentosEspacios() *regexp.Regexp {
	reg, err := regexp.Compile(`^[a-zA-ZáéíóúÁÉÍÓÚ ]+$`)
	if err != nil {
		gecko.FatalFmt("regex: %v", err)
	}
	return reg
}

// Acepta letras minúsculas y guiones.
func RegexMinusculasGuiones() *regexp.Regexp {
	reg, err := regexp.Compile(`^[a-z-]+$`)
	if err != nil {
		gecko.FatalFmt("regex: %v", err)
	}
	return reg
}

// Acepta letras minúsculas sin espacios.
func RegexMinusculas() *regexp.Regexp {
	reg, err := regexp.Compile(`^[a-z]+$`)
	if err != nil {
		gecko.FatalFmt("regex: %v", err)
	}
	return reg
}

// Acepta solo 3 letras minúsculas. Ej. "abc".
func RegexTresLetras() *regexp.Regexp {
	reg, err := regexp.Compile(`^[a-z-]{3}$`)
	if err != nil {
		gecko.FatalFmt("regex: %v", err)
	}
	return reg
}

// Acepta letras minúsculas y mayúsculas, números y guiones.
func RegexAlfanumericoGuiones() *regexp.Regexp {
	reg, err := regexp.Compile(`^[a-zA-Z1-9-]+$`)
	if err != nil {
		gecko.FatalFmt("regex: %v", err)
	}
	return reg
}

// QuitarEspacios elimina todos los espacios dobles del texto.
func QuitarEspacios(txt string) string {
	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(txt, " "))
}
