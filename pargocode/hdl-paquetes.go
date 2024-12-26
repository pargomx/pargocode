package main

import (
	"monorepo/appdominio"

	"github.com/pargomx/gecko"
)

func (s *servidor) getPaquetes(c *gecko.Context) error {
	paquetes, err := s.ddd.ListPaquetes()
	if err != nil {
		return err
	}
	PaquetesConEntidades := []appdominio.Paquete{}
	for _, paq := range paquetes {
		tablas, err := s.ddd.ListTablasByPaqueteID(paq.PaqueteID)
		if err != nil {
			return err
		}
		consultas, err := s.ddd.ListConsultasByPaqueteID(paq.PaqueteID)
		if err != nil {
			return err
		}
		PaquetesConEntidades = append(PaquetesConEntidades, appdominio.Paquete{Paquete: paq, Tablas: tablas, Consultas: consultas})
	}

	data := map[string]any{
		"Titulo":               "Paquetes",
		"PaquetesConEntidades": PaquetesConEntidades,
	}
	return c.RenderOk("ddd/paquetes", data)
}

// ================================================================ //
// ================================================================ //

func (s *servidor) getMapaEntidadRelacion(c *gecko.Context) error {
	tablas, err := s.ddd.ListTablas()
	if err != nil {
		return err
	}
	mermaid := "classDiagram\n"
	relaciones := "\n"
	for _, tbl := range tablas {
		mermaid += "class " + tbl.NombreRepo + " {\n"
		campos, err := s.ddd.ListCamposByTablaID(tbl.TablaID)
		if err != nil {
			return err
		}
		for _, col := range campos {
			mermaid += "\t" + col.TipoSql + " " + col.NombreColumna + "\n"
			if col.ReferenciaCampo != nil {
				ref, err := s.ddd.GetCampo(*col.ReferenciaCampo)
				if err != nil {
					return err
				}
				tRef, err := s.ddd.GetTabla(ref.TablaID)
				if err != nil {
					return err
				}
				relaciones += tbl.NombreRepo + " --> " + tRef.NombreRepo + " : " + ref.NombreColumna + "\n"
			}
		}
		mermaid += "}\n"
	}
	mermaid += relaciones
	return c.StatusOk(mermaid)
}

func (s *servidor) getMapaEntidadRelacionPaquete(c *gecko.Context) error {
	tablas, err := s.ddd.ListTablasByPaqueteID(c.PathInt("paquete_id"))
	if err != nil {
		return err
	}
	mermaid := "classDiagram\n"
	relaciones := "\n"
	for _, tbl := range tablas {
		mermaid += "class " + tbl.NombreRepo + " {\n"
		campos, err := s.ddd.ListCamposByTablaID(tbl.TablaID)
		if err != nil {
			return err
		}
		for _, col := range campos {
			mermaid += "\t" + col.TipoSql + " " + col.NombreColumna + "\n"
			if col.ReferenciaCampo != nil {
				ref, err := s.ddd.GetCampo(*col.ReferenciaCampo)
				if err != nil {
					return err
				}
				tRef, err := s.ddd.GetTabla(ref.TablaID)
				if err != nil {
					return err
				}
				relaciones += tbl.NombreRepo + " --> " + tRef.NombreRepo + " : " + ref.NombreColumna + "\n"
			}
		}
		mermaid += "}\n"
	}
	mermaid += relaciones
	return c.StatusOk(mermaid)
}
