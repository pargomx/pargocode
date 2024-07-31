package tmplutils

import (
	"strings"
	"text/template"
	"time"

	"monorepo/textutils"

	"github.com/pargomx/gecko/gko"
)

// funcMap contiene funciones útiles para usar dentro de las plantillas.
var funcMap = template.FuncMap{
	"suma": func(num ...int) int {
		var res int
		for _, n := range num {
			res = res + n
		}
		return res
	},
	"resta": func(num ...int) int {
		if len(num) == 0 {
			gko.LogWarnf("template.resta: llamada func resta sin argumentos")
			return 0
		}
		var res int
		for i, n := range num {
			if i == 0 {
				res = n
			} else {
				res = res - n
			}
		}
		return res
	},
	"mult": func(a, b int) int {
		return a * b
	},
	"div": func(a, b int) int {
		return a / b
	},
	"divf": func(a, b float64) float64 {
		return a / b
	},

	"br": func() string { // Agregar salto de línea
		return "\n"
	},

	"timestamp": func() string {
		return time.Now().Format("2006-01-02 15:04:05 MST")
	},

	"separador":       textutils.Separador,
	"separadorSimple": textutils.SeparadorSimple,

	"lower": strings.ToLower,
	"upper": strings.ToUpper,
}
