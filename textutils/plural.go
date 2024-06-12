package textutils

import (
	"strings"
)

// ========================================================================== //
// ========================================================================== //

// Deduce el nomre plural a partir de su singular.
// Ejemplo: "Cosa" se transfoma en "Cosas".
// Ejemplo: "mucho camión" > "muchos camiones"
func DeducirNombrePlural(input string) (plural string) {
	cosaDeAlgo := false

	palabras := strings.Split(input, " ")
	for _, sing := range palabras {

		if len(sing) == 0 {
			continue
		}

		// Excepción: en "cosa de algo" lo segundo se conserva en singular.
		if sing == "de" || cosaDeAlgo {
			cosaDeAlgo = true
			plural = plural + sing + " "
			continue
		}

		// Caso especial
		if strings.HasSuffix(sing, "ón") {
			plural = plural + strings.TrimSuffix(sing, "ón") + "ones "
			continue
		}

		// Caso normal: última letra de la palabra
		switch sing[len(sing)-1:] {

		case "a", "e", "i", "o", "u":
			plural = plural + sing + "s "

		case "í", "n", "r":
			plural = plural + sing + "es "

		case "z":
			plural = plural + sing[:len(sing)-1] + "ces "

		case "s":
			plural = plural + sing + " "

		default: // Cualquier otra consonante
			plural = plural + sing + "es "
		}

	}
	return strings.TrimSuffix(plural, " ")
}
