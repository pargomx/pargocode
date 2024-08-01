package main

import (
	"monorepo/appdominio"
	"monorepo/ddd"
	"monorepo/textutils"

	"github.com/pargomx/gecko"
)

func (s *servidor) agregarPaquete(c *gecko.Context) error {
	paq := ddd.Paquete{
		PaqueteID: ddd.NewPaqueteID(),
		Nombre:    textutils.QuitarAcentos(c.PromptLower()),
	}
	err := appdominio.AgregarPaquete(paq, s.ddd)
	if err != nil {
		return err
	}
	return c.Redir("/paquetes")
}

func (s *servidor) eliminarPaquete(c *gecko.Context) error {
	err := appdominio.EliminarPaquete(c.PathInt("paquete_id"), s.ddd)
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
	err := appdominio.ActualizarPaquete(paq, s.ddd)
	if err != nil {
		return err
	}
	return c.Redir("/paquetes")
}

// ================================================================ //
// ================================================================ //

func (s *servidor) generarDePaqueteArchivos(c *gecko.Context) error {
	paq, err := s.ddd.GetPaquete(c.PathInt("paquete_id"))
	if err != nil {
		return err
	}
	reporte := "ARCHIVOS GENERADOS:\n\n"
	errores := []error{}
	tablas, consultas, err := appdominio.GetTablasYConsultas(paq.PaqueteID, s.ddd)
	if err != nil {
		return err
	}
	for _, tbl := range tablas {
		call := codeGenerator.TblGenerarArchivos(&tbl, c.PathVal("tipo"))
		reporte += call.Destino() + "\n"
		err = call.Generar()
		if err != nil {
			errores = append(errores, err)
		}
	}
	for _, con := range consultas {
		call := codeGenerator.QryGenerarArchivos(&con, c.PathVal("tipo"))
		reporte += call.Destino() + "\n"
		err = call.Generar()
		if err != nil {
			errores = append(errores, err)
		}
	}
	if len(errores) > 0 {
		reporte += "\nERRORES:\n\n"
		for _, e := range errores {
			reporte += e.Error() + "\n\n"
		}
	}
	return c.StatusOk(reporte)
}
