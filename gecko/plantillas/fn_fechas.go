package plantillas

import (
	"fmt"
	"time"
)

// ================================================================ //

// FechaText. TODO: mejorar y poner en español.
func FechaText(f time.Time) string {
	if f.IsZero() {
		return "Sin fecha"
	}
	return f.Format("02 Jan 2006 03:04 PM")
}

// ================================================================ //

func FechaEsp(f time.Time) string {
	mes := ""
	switch f.Month() {
	case time.January:
		mes = "enero"
	case time.February:
		mes = "febrero"
	case time.March:
		mes = "marzo"
	case time.April:
		mes = "abril"
	case time.May:
		mes = "mayo"
	case time.June:
		mes = "junio"
	case time.July:
		mes = "julio"
	case time.August:
		mes = "agosto"
	case time.September:
		mes = "septiembre"
	case time.October:
		mes = "octubre"
	case time.November:
		mes = "noviembre"
	case time.December:
		mes = "diciembre"
	default:
		mes = "error_en_mes"
	}
	return fmt.Sprintf("%v de %v de %v", f.Day(), mes, f.Year())
}
