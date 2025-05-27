package appdominio

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
			ZeroIsNull:      c.ZeroIsNull,
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
