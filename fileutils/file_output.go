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
	if !OuptutToConsole && !OutputToFile {
		gko.FatalExit("No se especificó a donde mandar el resultado: fileutils.OuptutToConsole/OutputToFile es false")
	}
	if OuptutToConsole {
		gko.LogEventof("Imprimiendo:")
		fmt.Print(codigo)
	}
	if !OutputToFile {
		return nil
	}

	fileOut, err := os.Create(filename)
	if err != nil {
		return err
	}
	if _, err = fileOut.WriteString(codigo); err != nil {
		return err
	}
	fileOut.Close()
	cmd := exec.Command("goimports", "-w", filename)
	if errOut, err := cmd.CombinedOutput(); err != nil {
		gko.LogWarnf(string(errOut))
	}
	cmd = exec.Command("goimports-reviser", "-project-name", "monorepo", "-file-path", filename)
	if errOut, err := cmd.CombinedOutput(); err != nil {
		gko.LogWarnf(string(errOut))
	}

	return nil
}
