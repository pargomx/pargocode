package fileutils

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pargomx/gecko/gko"
)

var (
	OuptutToConsole bool = false // Imprime a consola si true.
	OutputToFile    bool = true  // Imprime a nuevo archivo si true.
)

func GuardarGoCodeConfirm(filename string, codigo string) error {
	if !confirmadoPorUsuario(filename) {
		return nil
	}
	return GuardarGoCode(filename, codigo)
}

func GuardarGoCode(filename string, codigo string) error {
	op := gko.Op("GuardarGoCode")
	if !OuptutToConsole && !OutputToFile {
		gko.ErrNoDisponible.Msg("No se especificó a donde mandar el resultado: fileutils.OuptutToConsole/OutputToFile es false")
	}
	if OuptutToConsole {
		gko.LogEventof("Imprimiendo:")
		fmt.Print(codigo)
	}
	if !OutputToFile {
		return nil
	}
	if !Existe("go.mod") {
		return op.E(gko.ErrNoDisponible).Str("Pargo no parece estar en la raíz de un proyecto Go")
	}

	fileOut, err := os.Create(filename)
	if err != nil {
		return op.Err(err)
	}
	if _, err = fileOut.WriteString(codigo); err != nil {
		return op.Err(err)
	}
	fileOut.Close()
	cmd := exec.Command("goimports", "-w", filename)
	if _, err := cmd.CombinedOutput(); err != nil {
		op.Err(err).Log()
	}
	cmd = exec.Command("goimports-reviser", "-project-name", "github.com/pargomx/pargocode", "-file-path", filename)
	if _, err := cmd.CombinedOutput(); err != nil {
		op.Err(err).Log()
	}

	return nil
}

func GuardarPlainText(filename string, txt string) error {
	op := gko.Op("GuardarPlainText")
	if !OuptutToConsole && !OutputToFile {
		gko.ErrNoDisponible.Msg("No se especificó a donde mandar el resultado: fileutils.OuptutToConsole/OutputToFile es false")
	}
	if OuptutToConsole {
		gko.LogEventof("Imprimiendo:")
		fmt.Print(txt)
	}
	if !OutputToFile {
		return nil
	}
	if !Existe("go.mod") {
		return op.E(gko.ErrNoDisponible).Str("Pargo no parece estar en la raíz de un proyecto Go")
	}

	fileOut, err := os.Create(filename)
	if err != nil {
		return op.Err(err)
	}
	if _, err = fileOut.WriteString(txt); err != nil {
		return op.Err(err)
	}
	fileOut.Close()

	return nil
}
