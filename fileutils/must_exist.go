package fileutils

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pargomx/gecko"
)

func Existe(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrNotExist)
}

// Verifica que el archivo exista.
func FileExist(absolutePath string) bool {
	if !filepath.IsAbs(absolutePath) {
		gecko.FatalFmt("fileExist: path debe ser absoluta: %v", absolutePath)
	}
	stat, err := os.Stat(absolutePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false
		}
		gecko.FatalFmt("fileExist: %v", err)
	}
	if stat.IsDir() {
		gecko.FatalFmt("fileExist: es un directorio: %v", absolutePath)
	}
	return true
}

// Verifica que el archivo exista. Si no cumple los requisitos
// sale del programa con un error fatal.
func FileMustExist(absolutePath string) {
	if !filepath.IsAbs(absolutePath) {
		gecko.FatalFmt("fileMustExist: path debe ser absoluta: %v", absolutePath)
	}
	stat, err := os.Stat(absolutePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			gecko.FatalFmt("fileMustExist: %v", absolutePath)
		}
		gecko.FatalFmt("fileMustExist: %v", err)
	}
	if stat.IsDir() {
		gecko.FatalFmt("fileMustExist: es un directorio: %v", absolutePath)
	}
}

// Verifica que el directorio exista. Si no cumple los
// requisitos sale del programa con un error fatal.
func DirMustExist(absolutePath string) {
	if !filepath.IsAbs(absolutePath) {
		gecko.FatalFmt("dirMustExist: path debe ser absoluta: %v", absolutePath)
	}
	stat, err := os.Stat(absolutePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			gecko.FatalFmt("dirMustExist: %v", absolutePath)
		}
		gecko.FatalFmt("dirMustExist: %v", err)
	}
	if !stat.IsDir() {
		gecko.FatalFmt("dirMustExist: es un directorio")
	}
}
