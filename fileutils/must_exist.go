package fileutils

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pargomx/gecko/gko"
)

func Existe(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrNotExist)
}

// Verifica que el archivo exista.
func FileExist(absolutePath string) bool {
	if !filepath.IsAbs(absolutePath) {
		gko.FatalExitf("fileExist: path debe ser absoluta: %v", absolutePath)
	}
	stat, err := os.Stat(absolutePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false
		}
		gko.FatalExitf("fileExist: %v", err)
	}
	if stat.IsDir() {
		gko.FatalExitf("fileExist: es un directorio: %v", absolutePath)
	}
	return true
}

// Verifica que el archivo exista. Si no cumple los requisitos
// sale del programa con un error fatal.
func FileMustExist(absolutePath string) {
	if !filepath.IsAbs(absolutePath) {
		gko.FatalExitf("fileMustExist: path debe ser absoluta: %v", absolutePath)
	}
	stat, err := os.Stat(absolutePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			gko.FatalExitf("fileMustExist: %v", absolutePath)
		}
		gko.FatalExitf("fileMustExist: %v", err)
	}
	if stat.IsDir() {
		gko.FatalExitf("fileMustExist: es un directorio: %v", absolutePath)
	}
}

// Verifica que el directorio exista. Si no cumple los
// requisitos sale del programa con un error fatal.
func DirMustExist(absolutePath string) {
	if !filepath.IsAbs(absolutePath) {
		gko.FatalExitf("dirMustExist: path debe ser absoluta: %v", absolutePath)
	}
	stat, err := os.Stat(absolutePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			gko.FatalExitf("dirMustExist: %v", absolutePath)
		}
		gko.FatalExitf("dirMustExist: %v", err)
	}
	if !stat.IsDir() {
		gko.FatalExitf("dirMustExist: es un directorio")
	}
}
