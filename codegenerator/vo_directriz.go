package codegenerator

import (
	"fmt"
	"regexp"
	"strings"
)

// Directriz es una cadena de texto que contiene una clave y un valor separados por dos puntos.
// Si el valor es una lista de valores, estos se separan por comas.
// Elimina todos los espacios.
//
// Ejemplo: "clave:valor", "clave:valor1,valor2,valor3"
type Directriz string

func (d Directriz) String() string {
	return string(d)
}

func (d Directriz) Key() string {
	return strings.Split(d.String(), ":")[0]
}

func (d Directriz) Values() []string {
	str := strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(d.String(), ""))
	kv := strings.Split(str, ":")
	if len(kv) == 2 {
		return strings.Split(kv[1], ",")

	} else if len(kv) == 1 {
		return []string{}

	} else {
		fmt.Printf("Warning: Directriz con más de un ':' (%s)", d.String())
		return []string{}
	}
}

// Convertir un string con directrices separadas por "\n".
func ToDirectrices(str string) []Directriz {
	dirs := []Directriz{}
	for _, v := range strings.Split(str, "\n") {
		v = strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(v, ""))
		if v == "" {
			continue
		}
		dirs = append(dirs, Directriz(v))
	}
	return dirs
}
