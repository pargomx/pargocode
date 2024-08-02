package codegenerator

import (
	"monorepo/fileutils"
	"os"
	"path/filepath"
	"strings"

	"github.com/pargomx/gecko/gko"
)

func (c tblGenCall) GenerarToFile(tipo string) error {
	c.tipo = tipo
	c.filename = c.getDestino()
	if !fileutils.Existe(filepath.Dir(c.filename)) {
		c.mkdir = true
	}
	if tipo == "sqlite" {
		if !fileutils.Existe(filepath.Join(c.tbl.Paquete.Directorio, "sqlite"+c.tbl.Paquete.Nombre, "servicio_repo.go")) {
			err := c.GenerarToFile("sqlite/servicio")
			if err != nil {
				gko.LogError(err)
			}
		}
	}

	codigo, err := c.GenerarToString(c.tipo)
	if err != nil {
		return err
	}
	if c.mkdir { // Crear directorio padre si no existe
		err := os.MkdirAll(filepath.Dir(c.filename), 0755)
		if err != nil {
			return err
		}
	}
	return fileutils.GuardarGoCode(c.filename, codigo)
}

// ================================================================ //
// ================================================================ //

func (c tblGenCall) GetInfoDestino(tipo string) string {
	c.tipo = tipo
	c.filename = c.getDestino()
	if !fileutils.Existe(filepath.Dir(c.filename)) {
		c.mkdir = true
	}
	if c.mkdir {
		return c.filename + " (creará directorio)"
	}
	return c.filename
}

func (c tblGenCall) getDestino() string {
	destino := ""
	switch c.tipo {
	case "entidad":
		destino = filepath.Join(c.tbl.Paquete.Directorio, c.tbl.Paquete.Nombre, "t_"+c.tbl.Tabla.Kebab+".go")

	case "mysql":
		destino = filepath.Join(c.tbl.Paquete.Directorio, "mysql"+c.tbl.Paquete.Nombre, "s_"+c.tbl.Tabla.NombreRepo+"_gen.go")

	case "sqlite":
		destino = filepath.Join(c.tbl.Paquete.Directorio, "sqlite"+c.tbl.Paquete.Nombre, "s_"+c.tbl.Tabla.NombreRepo+"_gen.go")

	case "sqlite/servicio":
		destino = filepath.Join(c.tbl.Paquete.Directorio, "sqlite"+c.tbl.Paquete.Nombre, "servicio_repo.go")

	case "mysql/servicio":
		destino = filepath.Join(c.tbl.Paquete.Directorio, "mysql"+c.tbl.Paquete.Nombre, "servicio_repo.go")

	default:
		destino = "generado.go"
	}
	destino = strings.TrimSuffix(destino, "/") // debe ser un archivo
	destino = strings.TrimPrefix(destino, "/") // debe ser relativa desde workdir
	return destino
}
