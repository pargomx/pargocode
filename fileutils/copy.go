package fileutils

import (
	"fmt"
	"os"
)

// Copia el archivo source a dest con los mismos permisos.
func Copy(src, dst string) error {

	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	input, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	err = os.WriteFile(dst, input, srcStat.Mode())
	if err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	return nil
}
