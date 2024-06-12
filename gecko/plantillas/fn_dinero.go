package plantillas

import (
	"fmt"
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// ================================================================ //

// CentavosToDinero convierte una cantidad entera
// de centavos en una cadena representando pesos.
//
// Ejemplo: 130400 > $ 1,304.00
func CentavosToDinero(cent int) string {
	if cent == 0 {
		return "$0.00"
	}

	if cent < 100 && cent > -100 { // Son solo centavos, negativos o positivos.
		switch {
		case cent > 9: // $0.10 - $0.99
			return "$0." + strconv.Itoa(cent)
		case cent > 0: // $0.01 - $0.09
			return "$0.0" + strconv.Itoa(cent)
		case cent > -9: // -$0.09 - $0.01
			return "-$0.0" + strconv.Itoa(-cent)
		case cent > -100: // -$0.10 - -$0.99
			return "-$0." + strconv.Itoa(-cent)
		}

	}

	m := message.NewPrinter(language.LatinAmericanSpanish)
	x := float64(cent)
	x = x / 100

	txt := m.Sprintf("%.2f", x)

	if cent < 0 { // Si es negativo se pone el "-" antes de "$"
		txt = txt[:1] + "$" + txt[1:]
	} else {
		txt = "$" + txt
	}

	return txt
}

// ================================================================ //

// CentavosToDinero convierte una cantidad entera
// de centavos en una cantidad entera de pesos.
//
// Corta los centavos.
//
// Ejemplo: 130475 > 1304
func CentavosToPesos(cent int) int {
	if cent == 0 {
		return 0
	}
	if cent < 100 && cent > -100 { // Mínimo de un peso
		fmt.Println("plantillas.CentavosToPesos: convirtiendo menos de 1 peso a pesos", cent)
		return 1
	}
	return cent / 100
}

// ================================================================ //

// CentavosToDinero convierte una cantidad entera
// de centavos en una cadena representando pesos
// que puede mandarse a PayPal como monto.
//
// Ejemplo: 130400 > 1304.00
func CentavosToDineroPayPal(cent int) string {
	return fmt.Sprintf("%.2f", float64(cent)/100)
}

// ================================================================ //
