package codegenerator

import (
	"monorepo/fileutils"
	"os"
	"path/filepath"
	"strings"

	"github.com/pargomx/gecko/gko"
)

type tblGenCall struct {
	filename string
	tipo     string
	tbl      *Tabla
	gen      *Generador
	mkdir    bool
}

func (c tblGenCall) Destino() string {
	if c.mkdir {
		return c.filename + " (creará directorio)"
	}
	return c.filename
}

func (c tblGenCall) Generar() error {
	codigo, err := c.tbl.GenerarDeTablaString(c.tipo)
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

func (tbl *Tabla) TblGenerarArchivos(tipo string) tblGenCall {
	c := tblGenCall{
		filename: "generado.go",
		tipo:     tipo,
		tbl:      tbl,
		gen:      tbl.Generador,
	}
	switch tipo {
	case "entidad":
		c.filename = filepath.Join(tbl.Paquete.Directorio, tbl.Paquete.Nombre, "t_"+tbl.Tabla.Kebab+".go")

	case "mysql", "mysql-directriz":
		c.filename = filepath.Join(tbl.Paquete.Directorio, "mysql"+tbl.Paquete.Nombre, "s_"+tbl.Tabla.NombreRepo+"_gen.go")
		c.tipo = "mysql-directriz"

	case "sqlite", "sqlite-directriz":
		c.filename = filepath.Join(tbl.Paquete.Directorio, "sqlite"+tbl.Paquete.Nombre, "s_"+tbl.Tabla.NombreRepo+"_gen.go")
		c.tipo = "mysql-directriz"

		if !fileutils.Existe(filepath.Join(tbl.Paquete.Directorio, "sqlite"+tbl.Paquete.Nombre, "servicio_repo.go")) {
			err := tbl.TblGenerarArchivos("sqlite/servicio").Generar()
			if err != nil {
				gko.LogError(err)
			}
		}

	case "sqlite/servicio":
		c.filename = filepath.Join(tbl.Paquete.Directorio, "sqlite"+tbl.Paquete.Nombre, "servicio_repo.go")

	case "mysql/servicio":
		c.filename = filepath.Join(tbl.Paquete.Directorio, "mysql"+tbl.Paquete.Nombre, "servicio_repo.go")

	}
	c.filename = strings.TrimSuffix(c.filename, "/") // debe ser relativa desde workdir
	if !fileutils.Existe(filepath.Dir(c.filename)) {
		c.mkdir = true
	}
	return c
}
