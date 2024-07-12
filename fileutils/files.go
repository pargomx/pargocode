package fileutils

import (
	"errors"
	"os"

	"github.com/pargomx/gecko"
)

// NuevaCarpeta intenta crear una nueva carpeta
// con los permisos definidos en dirPerms. Si la
// carpeta ya existía no pasa nada. Si no puede
// crearla aborta el programa con un error fatal.
func NuevaCarpeta(path string) {

	err := os.Mkdir(path, dirPerms)
	if errors.Is(err, os.ErrExist) {
		return
	}
	if err != nil {
		gecko.FatalFmt("nuevaCarpeta: ["+path+"]", err)
	}
}

// Pregunta al usuario antes de intentar crear la carpeta.
// No pregunta si ya existe el directorio.
func NuevaCarpetaConfirm(path string) {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return
	}
	if !confirmadoPorUsuario(path) {
		return
	}
	err = os.Mkdir(path, dirPerms)
	if errors.Is(err, os.ErrExist) {
		return
	}
	if err != nil {
		gecko.FatalFmt("nuevaCarpeta: ["+path+"]", err)
	}
}

// NuevoArchivo intenta crear un nuevo archivo. Si el archivo
// ya existía no pasa nada, lo deja como está. Si no puede
// crearlo aborta el programa con un error fatal.
func NuevoArchivo(path string) {

	_, err := os.Stat(path)
	if err == nil {
		return // no hacer nada si ya existe
	}

	if !errors.Is(err, os.ErrNotExist) {
		gecko.FatalFmt("NuevoArchivo: ["+path+"]", err)
	}

	_, err = os.Create(path)
	if err != nil {
		gecko.FatalFmt("NuevoArchivo: ["+path+"]", err)
	}
}
