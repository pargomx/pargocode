package dpaquete

import "github.com/pargomx/gecko"

func GetAgregadoConsulta(consultaID int, repo Repositorio) (*Consulta, error) {
	ctx := gecko.NewOp("GetAgregadoConsulta").Ctx("consultaID", consultaID)
	consulta, err := repo.GetConsulta(consultaID)
	if err != nil {
		return nil, ctx.Err(err)
	}
	paquete, err := repo.GetPaquete(consulta.PaqueteID)
	if err != nil {
		return nil, ctx.Err(err).Op("GetPaqueteDeConsulta")
	}
	tablaFrom, err := GetTabla(consulta.TablaID, repo)
	if err != nil {
		return nil, ctx.Err(err).Op("GetTablaDeOrigen")
	}

	relaciones, err := repo.ListConsultaRelacionesByConsultaID(consulta.ConsultaID)
	if err != nil {
		return nil, ctx.Err(err).Op("GetRelacionesDeConsulta")
	}
	campos, err := repo.ListConsultaCamposByConsultaID(consulta.ConsultaID)
	if err != nil {
		return nil, ctx.Err(err).Op("GetCamposDeConsulta")
	}

	agregado := Consulta{
		Paquete:     *paquete,
		Consulta:    *consulta,
		TablaOrigen: tablaFrom.Tabla,
		From:        *tablaFrom,
		// Campos:      camposConOrigen,
		// Relaciones:  relacionesConOrigen,
	}

	for _, campo := range campos {
		campoConOrigen := CampoConsulta{
			ConsultaID:  campo.ConsultaID,
			Posicion:    campo.Posicion,
			CampoID:     campo.CampoID,
			Expresion:   campo.Expresion,
			AliasSql:    campo.AliasSql,
			NombreCampo: campo.NombreCampo,
			TipoGo:      campo.TipoGo,
			Pk:          campo.Pk,
			Filtro:      campo.Filtro,
			GroupBy:     campo.GroupBy,
			Descripcion: campo.Descripcion,
			Consulta:    &agregado,
		}
		// Traer info de origen si aplica
		if campo.CampoID != nil {
			cam, err := repo.GetCampo(*campo.CampoID)
			if err != nil {
				return nil, ctx.Err(err)
			}
			tbl, err := repo.GetTabla(cam.TablaID)
			if err != nil {
				return nil, ctx.Err(err)
			}
			paq, err := repo.GetPaquete(tbl.PaqueteID)
			if err != nil {
				return nil, ctx.Err(err)
			}
			campoConOrigen.OrigenPaquete = paq
			campoConOrigen.OrigenTabla = tbl
			campoConOrigen.OrigenCampo = cam
		}
		agregado.Campos = append(agregado.Campos, campoConOrigen)

	}

	for _, relacion := range relaciones {
		joinTbl, err := GetTabla(relacion.JoinTablaID, repo)
		if err != nil {
			return nil, ctx.Err(err)
		}
		rel := Relacion{
			ConsultaID:  relacion.ConsultaID,
			Posicion:    relacion.Posicion,
			TipoJoin:    relacion.TipoJoin,
			JoinTablaID: relacion.JoinTablaID,
			JoinAs:      relacion.JoinAs,
			JoinOn:      relacion.JoinOn,
			FromTablaID: relacion.FromTablaID,
			Join:        *joinTbl,
		}
		if relacion.FromTablaID != 0 {
			fromTbl, err := GetTabla(relacion.FromTablaID, repo)
			if err != nil {
				return nil, ctx.Err(err)
			}
			rel.From = fromTbl
		}
		agregado.Relaciones = append(agregado.Relaciones, rel)
	}

	return &agregado, nil
}
