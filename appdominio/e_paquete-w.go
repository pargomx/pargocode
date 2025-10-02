package appdominio

import (
	"strings"

	"github.com/pargomx/pargocode/ddd"
	"github.com/pargomx/pargocode/textutils"

	"github.com/pargomx/gecko/gko"
)

func AgregarPaquete(paq ddd.Paquete, repo Repositorio) error {
	op := gko.Op("AgregarPaquete")
	if paq.PaqueteID == 0 {
		return op.Msg("Debe especificar un nuevo paqueteID")
	}
	if paq.Nombre == "" {
		return op.Msg("Debe especificar un nombre para el nuevo paquete")
	}
	if repo.ExistePaquete(paq.PaqueteID, paq.Nombre) {
		return op.Msgf("Ya existe un paquete con ID %v o nombre '%v'", paq.PaqueteID, paq.Nombre)
	}
	err := repo.InsertPaquete(paq)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func ActualizarPaquete(paq ddd.Paquete, repo Repositorio) error {
	op := gko.Op("ActualizarPaquete")
	if paq.PaqueteID == 0 {
		return op.Msg("paqueteID no especificado")
	}
	if paq.Nombre == "" {
		return op.Msg("Debe especificar un nombre para el paquete")
	}
	if paq.GoModule == "" {
		return op.Msg("Debe especificar un go_module para el paquete")
	}
	actualizado := ddd.Paquete{
		PaqueteID:   paq.PaqueteID,
		GoModule:    strings.ToLower(textutils.QuitarAcentos(strings.ReplaceAll(paq.GoModule, " ", ""))),
		Directorio:  textutils.QuitarAcentos(strings.ReplaceAll(paq.Directorio, " ", "")),
		Nombre:      strings.ToLower(textutils.QuitarAcentos(strings.ReplaceAll(paq.Nombre, " ", ""))),
		Descripcion: paq.Descripcion,
	}
	err := repo.UpdatePaquete(actualizado)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func EliminarPaquete(paqueteID int, repo Repositorio) error {
	op := gko.Op("EliminarPaquete").Ctx("paquete_id", paqueteID)
	paquete, err := GetPaquete(repo, paqueteID)
	if err != nil {
		return op.Err(err)
	}
	if len(paquete.Tablas) != 0 || len(paquete.Consultas) != 0 {
		return op.Msgf("Es necesario primero quitar sus tablas (%d) y consultas (%v)", len(paquete.Tablas), len(paquete.Consultas))
	}
	err = repo.DeletePaquete(paqueteID)
	if err != nil {
		return op.Err(err)
	}
	return nil
}
