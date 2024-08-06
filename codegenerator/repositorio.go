package codegenerator

import (
	"fmt"
	"monorepo/ddd"

	"github.com/pargomx/gecko/gko"
)

type Repositorio interface {
	GetPaquete(paqueteID int) (*ddd.Paquete, error)
	ListTablasByPaqueteID(PaqueteID int) ([]ddd.Tabla, error)
	ListConsultasByPaqueteID(PaqueteID int) ([]ddd.Consulta, error)

	GetTabla(tablaID int) (*ddd.Tabla, error)
	ListCamposByTablaID(tablaID int) ([]ddd.Campo, error)
	GetCampo(campoID int) (*ddd.Campo, error)
	GetValoresEnum(CampoID int) ([]ddd.ValorEnum, error)

	GetConsulta(consultaID int) (*ddd.Consulta, error)
	ListConsultaRelacionesByConsultaID(consultaID int) ([]ddd.ConsultaRelacion, error)
	ListConsultaCamposByConsultaID(consultaID int) ([]ddd.ConsultaCampo, error)
}

// ================================================================ //
// ========== TABLA =============================================== //

func (gen *Generador) getTabla(tablaID int) (*tabla, error) {

	op := gko.Op("GetAgregadoTabla").Ctx("tablaID", tablaID)
	tbl, err := gen.db.GetTabla(tablaID)
	if err != nil {
		return nil, op.Err(err)
	}
	paquete, err := gen.db.GetPaquete(tbl.PaqueteID)
	if err != nil {
		return nil, op.Err(err)
	}
	agregado := tabla{
		Tabla:   *tbl,
		Paquete: *paquete,
	}
	campos, err := gen.db.ListCamposByTablaID(tbl.TablaID)
	if err != nil {
		return nil, op.Err(err)
	}
	for _, c := range campos {
		campo := CampoTabla{
			Paquete: paquete,
			Tabla:   tbl,
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
			campo.ValoresPosibles, err = gen.db.GetValoresEnum(c.CampoID)
			if err != nil {
				return nil, op.Err(err)
			}
		}

		if c.ReferenciaCampo != nil {
			campo.CampoFK, err = gen.db.GetCampo(*c.ReferenciaCampo)
			if err != nil {
				return nil, op.Err(err).Op("get_fk")
			}
			campo.TablaFK, err = gen.db.GetTabla(campo.CampoFK.TablaID)
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

// ================================================================ //
// ========== CONSULTA ============================================ //

func (gen *Generador) getConsulta(consultaID int) (*consulta, error) {
	ctx := gko.Op("getConsulta").Ctx("consultaID", consultaID)
	con, err := gen.db.GetConsulta(consultaID)
	if err != nil {
		return nil, ctx.Err(err)
	}
	paquete, err := gen.db.GetPaquete(con.PaqueteID)
	if err != nil {
		return nil, ctx.Err(err).Op("GetPaqueteDeConsulta")
	}
	tablaFrom, err := gen.getTabla(con.TablaID)
	if err != nil {
		return nil, ctx.Err(err).Op("GetTablaDeOrigen")
	}

	relaciones, err := gen.db.ListConsultaRelacionesByConsultaID(con.ConsultaID)
	if err != nil {
		return nil, ctx.Err(err).Op("GetRelacionesDeConsulta")
	}
	campos, err := gen.db.ListConsultaCamposByConsultaID(con.ConsultaID)
	if err != nil {
		return nil, ctx.Err(err).Op("GetCamposDeConsulta")
	}

	agregado := consulta{
		Paquete:     *paquete,
		Consulta:    *con,
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
			cam, err := gen.db.GetCampo(*campo.CampoID)
			if err != nil {
				return nil, ctx.Err(err)
			}
			tbl, err := gen.db.GetTabla(cam.TablaID)
			if err != nil {
				return nil, ctx.Err(err)
			}
			paq, err := gen.db.GetPaquete(tbl.PaqueteID)
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
		joinTbl, err := gen.getTabla(relacion.JoinTablaID)
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
			fromTbl, err := gen.getTabla(relacion.FromTablaID)
			if err != nil {
				return nil, ctx.Err(err)
			}
			rel.From = fromTbl
		}
		agregado.Relaciones = append(agregado.Relaciones, rel)
	}

	return &agregado, nil
}
