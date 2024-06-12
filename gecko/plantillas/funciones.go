package plantillas

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"
)

// funcMap contiene funciones útiles para transformar datos a texto.
var funcMap = template.FuncMap{

	// * ARITMÉTICA
	"suma": func(num ...int) int {
		var res int
		for _, n := range num {
			res = res + n
		}
		return res
	},
	"resta": func(num ...int) int {
		if len(num) == 0 {
			fmt.Sprintln("WARN: Llamada func resta sin argumentos")
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

	"lower": strings.ToLower,
	"upper": strings.ToUpper,

	// * POINTERS
	"derefInt": func(num *int) int {
		if num == nil {
			return 0
		}
		return *num
	},

	// * DINERO
	"dinero":       CentavosToDinero,
	"dineroPaypal": CentavosToDineroPayPal,
	"pesos":        CentavosToPesos,

	// * TIEMPO
	"horasMinutos":   fmtHorasMinutos,
	"numeroEnLetras": NumeroEnLetras,
	"edad":           EdadPtr,

	"fechaCompleta":    FechaEsp,
	"fechaCompletaHoy": func() string { return FechaEsp(time.Now()) },

	// * ARCHIVOS
	"filesize": ByteCountSI,

	// * STRINGS
	"concat": func(args ...any) string {
		return fmt.Sprintf(strings.Repeat("%v", len(args)), args...)
	},

	// ================================================================ //
	//  OTRAS COSAS

	"ciclo": func(ciclo int) string {
		switch ciclo {
		case -1:
			return "N/A"
		case -2:
			return "S/D"
		default:
			return strconv.Itoa(ciclo)
		}
	},

	"colorHash": func(val ...interface{}) string {
		str := fmt.Sprint(val...)
		hash := sha1.New()
		hash.Write([]byte(str))
		bytes := hash.Sum(nil)
		encodedStr := hex.EncodeToString(bytes)
		estilo := fmt.Sprintf(`style="color: #%v;"`, encodedStr[2:8])
		return estilo

	},

	"colorCalif": func(cal *int) string {
		if cal == nil {
			return ""
		}
		if *cal < 60 {
			return "text-red"
		}
		return ""
	},

	"addQueryParam": addQueryParam,
	"addQueryNum":   addQueryNum,
}
