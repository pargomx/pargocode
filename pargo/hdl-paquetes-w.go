package main

import (
	"monorepo/domain_driven_design/ddd"
	"monorepo/domain_driven_design/dpaquete"
	"monorepo/textutils"

	"github.com/pargomx/gecko"
)

func (s *servidor) agregarPaquete(c *gecko.Context) error {
	paq := ddd.Paquete{
		PaqueteID: ddd.NewPaqueteID(),
		Nombre:    textutils.QuitarAcentos(c.PromptLower()),
	}
	err := dpaquete.AgregarPaquete(paq, s.ddd)
	if err != nil {
		return err
	}
	return c.Redir("/paquetes")
}

func (s *servidor) eliminarPaquete(c *gecko.Context) error {
	err := dpaquete.EliminarPaquete(c.PathInt("paquete_id"), s.ddd)
	if err != nil {
		return err
	}
	return c.Redir("/paquetes")
}

func (s *servidor) actualizarPaquete(c *gecko.Context) error {
	paq := ddd.Paquete{
		PaqueteID:   c.PathInt("paquete_id"),
		GoModule:    c.FormVal("go_module"),
		Directorio:  c.FormVal("directorio"),
		Nombre:      c.FormVal("nombre"),
		Descripcion: c.FormVal("descripcion"),
	}
	err := dpaquete.ActualizarPaquete(paq, s.ddd)
	if err != nil {
		return err
	}
	return c.Redir("/paquetes")
}
