package textutils

import (
	"fmt"
	"io"
	"strings"
)

// Imprime un separador visual de 70 caracteres de largo con título.
func ImprimirSeparador(buf io.Writer, titulo string) {
	fmt.Fprint(buf, Separador(titulo))
}

// Imprime un separador visual de 70 caracteres de largo de una línea
func ImprimirSeparadorSimple(buf io.Writer) {
	fmt.Fprint(buf, SeparadorSimple())
}

// Imprime un separador visual de 70 caracteres de largo con título.
func Separador(titulo string) (separador string) {
	separador += fmt.Sprintln()
	separador += fmt.Sprintln("// ", strings.Repeat("=", 64), " //")
	separador += fmt.Sprintln("// ", strings.Repeat("=", 10), titulo, strings.Repeat("=", 52-len(titulo)), " //")
	separador += fmt.Sprintln()
	return separador
}

// Imprime un separador visual de 70 caracteres de largo de una línea
func SeparadorSimple() (separador string) {
	separador += fmt.Sprintln()
	separador += fmt.Sprintln("// ", strings.Repeat("=", 64), " //")
	separador += fmt.Sprintln()
	return separador
}
