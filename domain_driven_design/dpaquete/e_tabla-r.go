package dpaquete

import (
	"fmt"

	"github.com/pargomx/gecko/gko"
)

func GetTabla(tablaID int, repo Repositorio) (*Tabla, error) {

	op := gko.Op("GetAgregadoTabla").Ctx("tablaID", tablaID)
	tabla, err := repo.GetTabla(tablaID)
	if err != nil {
		return nil, op.Err(err)
	}
	paquete, err := repo.GetPaquete(tabla.PaqueteID)
	if err != nil {
		return nil, op.Err(err)
	}
	agregado := Tabla{
		Tabla:   *tabla,
		Paquete: *paquete,
	}
	campos, err := repo.ListCamposByTablaID(tabla.TablaID)
	if err != nil {
		return nil, op.Err(err)
	}
	for _, c := range campos {
		campo := CampoTabla{
			Paquete: paquete,
			Tabla:   tabla,
			Campo:   c,

			CampoID:         c.CampoID,
			TablaID:         c.TablaID,
			NombreCampo:     c.NombreCampo,
			NombreColumna:   c.NombreColumna,
			NombreHumano:    c.NombreHumano,
			TipoGo:          c.TipoGo,
			TipoSql:         c.TipoSql,
			Setter:          c.Setter,
			Importado:       c.Importado,
			PrimaryKey:      c.PrimaryKey,
			ForeignKey:      c.ForeignKey,
			Uq:              c.Uq,
			Req:             c.Req,
			Ro:              c.Ro,
			Filtro:          c.Filtro,
			Nullable:        c.Nullable,
			MaxLenght:       c.MaxLenght,
			Uns:             c.Uns,
			DefaultSql:      c.DefaultSql,
			Especial:        c.Especial,
			ReferenciaCampo: c.ReferenciaCampo,
			Expresion:       c.Expresion,
			EsFemenino:      c.EsFemenino,
			Descripcion:     c.Descripcion,
			Posicion:        c.Posicion,
		}
		if c.Especial {
			campo.ValoresPosibles, err = repo.GetValoresEnum(c.CampoID)
			if err != nil {
				return nil, op.Err(err)
			}
		}

		if c.ReferenciaCampo != nil {
			campo.CampoFK, err = repo.GetCampo(*c.ReferenciaCampo)
			if err != nil {
				return nil, op.Err(err).Op("get_fk")
			}
			campo.TablaFK, err = repo.GetTabla(campo.CampoFK.TablaID)
			if err != nil {
				return nil, op.Err(err).Op("get_fk")
			}
		}
		agregado.Campos = append(agregado.Campos, campo)
	}

	for _, fk := range agregado.ForeignKeys() {
		if fk.CampoFK == nil {
			fmt.Println("campo FK no tiene referencia", fk.NombreColumna)
		}
	}

	return &agregado, nil

}

func GetTablas(repo Repositorio) ([]Tabla, error) {
	ctx := gko.Op("GetAgregadosTabla")
	tablas, err := repo.ListTablas()
	if err != nil {
		return nil, ctx.Err(err)
	}
	items := []Tabla{}
	for _, t := range tablas {
		tbl, err := GetTabla(t.TablaID, repo)
		if err != nil {
			return nil, ctx.Err(err)
		}
		items = append(items, *tbl)
	}
	return items, nil
}

func GetTablasYConsultas(paqueteID int, repo Repositorio) ([]Tabla, []Consulta, error) {
	op := gko.Op("GetTablasYConsultas")
	tablas, err := repo.ListTablasByPaqueteID(paqueteID)
	if err != nil {
		return nil, nil, op.Err(err)
	}
	Tablas := []Tabla{}
	for _, t := range tablas {
		tbl, err := GetTabla(t.TablaID, repo)
		if err != nil {
			return nil, nil, op.Err(err)
		}
		Tablas = append(Tablas, *tbl)
	}

	consultas, err := repo.ListConsultasByPaqueteID(paqueteID)
	if err != nil {
		return nil, nil, op.Err(err)
	}
	Consultas := []Consulta{}
	for _, c := range consultas {
		consulta, err := GetConsulta(c.ConsultaID, repo)
		if err != nil {
			return nil, nil, op.Err(err)
		}
		Consultas = append(Consultas, *consulta)
	}
	return Tablas, Consultas, nil
}
