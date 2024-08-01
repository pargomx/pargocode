package appdominio

import (
	"errors"
	"fmt"
	"monorepo/ddd"

	"github.com/pargomx/gecko/gko"

	"github.com/davecgh/go-spew/spew"
)

// El paquete no puede tener dos tablas con el mismo nombre, ni repetir NombreItem con tablas o consultas.
func AgregarTabla(tbl ddd.Tabla, repo Repositorio) error {

	// Parámetros necesarios
	if tbl.TablaID == 0 {
		return errors.New("tabla sin TablaID")
	}
	if tbl.NombreRepo == "" {
		return errors.New("tabla sin NombreRepo")
	}
	if tbl.NombreItem == "" {
		return errors.New("tabla sin NombreItem")
	}
	if tbl.NombreItems == "" {
		return errors.New("tabla sin NombreItems")
	}
	if tbl.Abrev == "" {
		return errors.New("tabla sin Abrev")
	}
	if tbl.Kebab == "" {
		return errors.New("tabla sin Kebab")
	}

	// Para integridad de la base de datos
	pkg, err := GetPaquete(repo, tbl.PaqueteID)
	if err != nil {
		return err
	}

	// Validar pertenencia al paquete
	if tbl.PaqueteID != pkg.Paquete.PaqueteID {
		return fmt.Errorf("agregar tabla al paquete '%v' con paquete_id equivocado", pkg.Paquete.Nombre)
	}

	// Que no haya conflicto de nombre con otras tablas o consultas
	for _, t := range pkg.Tablas {
		if t.NombreRepo == tbl.NombreRepo {
			return fmt.Errorf("ya existe una tabla con el NombreRepo '%v'", tbl.NombreRepo)
		}
		if t.NombreItem == tbl.NombreItem {
			return fmt.Errorf("ya existe una tabla con el NombreItem '%v'", tbl.NombreItem)
		}
		if t.Abrev == tbl.Abrev {
			return fmt.Errorf("ya existe una tabla con el Abrev '%v'", tbl.Abrev)
		}
		if t.Kebab == tbl.Kebab {
			return fmt.Errorf("ya existe una tabla con el Kebab '%v'", tbl.Kebab)
		}
	}
	for _, c := range pkg.Consultas {
		if c.NombreItem == tbl.NombreItem {
			return fmt.Errorf("conflicto: existe una consulta con el NombreItem '%v'", tbl.NombreItem)
		}
	}

	// Agregar en la base de datos
	err = repo.InsertTabla(tbl)
	if err != nil {
		spew.Dump(tbl)
		return gko.Err(err).Ctx("tabla", tbl)
	}

	return nil
}

func ActualizarTabla(tablaID int, new ddd.Tabla, repo Repositorio) error {
	old, err := repo.GetTabla(tablaID)
	if err != nil {
		return err
	}

	if old.TablaID != new.TablaID {
		return gko.ErrDatoInvalido().Msg("No se puede cambiar el ID de la tabla").Ctx("oldID", old.TablaID).Ctx("newID", new.TablaID)
	}

	if new.TablaID == 0 {
		return gko.ErrDatoInvalido().Msg("No se puede quitar el ID")
	}
	if new.PaqueteID == 0 {
		return gko.ErrDatoInvalido().Msg("No se puede quitar el paquete")
	}
	if new.NombreRepo == "" {
		return gko.ErrDatoInvalido().Msg("El nombre de la tabla no puede estar vacío")
	}
	if new.NombreItem == "" {
		return gko.ErrDatoInvalido().Msg("El nombre del item no puede estar vacío")
	}
	if new.NombreItems == "" {
		return gko.ErrDatoInvalido().Msg("El nombre items plural no puede estar vacío")
	}
	if new.Humano == "" {
		return gko.ErrDatoInvalido().Msg("El nombre humano no puede estar vacío")
	}
	if new.Abrev == "" {
		return gko.ErrDatoInvalido().Msg("La abreviatura no puede estar vacía")
	}

	old.PaqueteID = new.PaqueteID
	old.NombreRepo = new.NombreRepo
	old.NombreItem = new.NombreItem
	old.NombreItems = new.NombreItems
	old.Abrev = new.Abrev
	old.Humano = new.Humano
	old.HumanoPlural = new.HumanoPlural
	old.Kebab = new.Kebab
	old.EsFemenino = new.EsFemenino
	old.Descripcion = new.Descripcion
	old.Directrices = new.Directrices

	err = repo.UpdateTabla(*old)
	if err != nil {
		return err
	}
	return nil
}
