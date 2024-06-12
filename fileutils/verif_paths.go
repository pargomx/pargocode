package fileutils

import (
	"fmt"
	"path"
	"strings"
)

// Verifica que sea /path/absoluta/ y agrega trailing slash.
func PathMustBeAbsolute(Path *string, tipo string) error {
	if Path == nil {
		return fmt.Errorf("path '%s' es nil", tipo)
	}
	if *Path == "" {
		return fmt.Errorf("path '%s' no definido", tipo)
	}
	if !path.IsAbs(*Path) {
		return fmt.Errorf("path '%s' debe ser absoluto pero fue '%s'", tipo, *Path)
	}
	*Path = strings.TrimSuffix(*Path, "/") + "/"
	return nil
}

func PathMustBeAbsoluteWithDefault(Path *string, tipo string, defaultPath string) error {
	if Path == nil {
		return fmt.Errorf("path '%s' es nil", tipo)
	}
	if *Path == "" {
		*Path = defaultPath
	}
	if *Path == "" {
		return fmt.Errorf("path '%s' no definido", tipo)
	}
	if (*Path)[0] != '/' {
		return fmt.Errorf("path '%s' debe ser absoluto pero fue '%s'", tipo, *Path)
	}
	*Path = strings.TrimSuffix(*Path, "/") + "/"
	return nil
}

func FileMustBeAbsoluteWithDefault(Path *string, tipo string, defaultPath string) error {
	if Path == nil {
		return fmt.Errorf("path '%s' es nil", tipo)
	}
	if *Path == "" {
		*Path = defaultPath
	}
	if !path.IsAbs(*Path) {
		return fmt.Errorf("path '%s' debe ser absoluto pero fue '%s'", tipo, *Path)
	}
	*Path = strings.TrimSuffix(*Path, "/")
	return nil
}
