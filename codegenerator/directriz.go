package codegenerator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pargomx/gecko/gko"
)

// Directriz es una cadena de texto que contiene una clave y un valor separados por dos puntos.
// Si el valor es una lista de valores, estos se separan por comas.
// Elimina todos los espacios.
//
// Ejemplo: "clave:valor", "clave:valor1,valor2,valor3"
//
// Ejemplo custom: "custom_l:BySomeValue:param1,param2:JOIN s WHERE x > 2 ORDER BY x:param1,param1,param2"
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
		if strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(v, "")) == "" {
			continue
		}
		dirs = append(dirs, Directriz(strings.TrimSpace(v)))
	}
	return dirs
}

// ================================================================ //

type customList struct {
	CompFunc string // sufijo para el nombre de la función: "Nuevos" -> "ListUsuariosNuevos"
	ArgsFunc string // parámetros de la función: "param1,param1,param2"
	CompSQL  string // SQL después del FROM: "JOIN s WHERE x > 2 ORDER BY x"
	ArgsSQL  string // parámetros del query: "param1,param1,param2"
}

// Ejemplo custom: "custom_l:BySomeValue:param1,param2:JOIN s WHERE x > 2 ORDER BY x:param1,param1,param2"
func (d Directriz) CustomList() (*customList, error) {
	str := strings.TrimSpace(d.String())
	kv := strings.Split(str, ":")
	if len(kv) > 5 {
		return nil, gko.Op("CustomList").Msg("directriz con más de 4 secciones").Ctx("d", d.String())
	}
	if len(kv) < 4 {
		return nil, gko.Op("CustomList").Msg("directriz con menos de 3 secciones").Ctx("d", d.String())
	}
	cl := customList{
		CompFunc: strings.TrimSpace(kv[1]),
		ArgsFunc: strings.TrimSpace(kv[2]),
		CompSQL:  strings.TrimSpace(kv[3]),
	}
	if len(kv) == 5 {
		cl.ArgsSQL = strings.TrimSpace(kv[4])
	}
	return &cl, nil
}
