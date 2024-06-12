package plantillas

import (
	"fmt"

	numaletra "github.com/jtorz/num-a-letra"
)

// Convierte un número a su representación textual.
// Por ejemplo 97 se convierte en "noventa y siete".
func NumeroEnLetras(n int) string {
	res, err := numaletra.IntLetra(n)
	if err != nil {
		fmt.Println("plantillas.NumeroEnLetras:", err)
	}
	return res
}
