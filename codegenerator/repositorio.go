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
// ========== PAQUETE ============================================= //

func GetTablasYConsultas(paqueteID int, repo Repositorio) ([]Tabla, []Consulta, error) {
	op := gko.Op("GetTablasYConsultas")
	tablas, err := repo.ListTablasByPaqueteID(paqueteID)
	if err != nil {
		return nil, nil, op.Err(err)
	}
	Tablas := []Tabla{}
	for _, t := range tablas {
		tbl, err := getTabla(t.TablaID, repo)
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

// ================================================================ //
// ========== TABLA =============================================== //

func getTabla(tablaID int, repo Repositorio) (*Tabla, error) {

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

// ================================================================ //
// ========== CONSULTA ============================================ //

func GetConsulta(consultaID int, repo Repositorio) (*Consulta, error) {
	ctx := gko.Op("GetAgregadoConsulta").Ctx("consultaID", consultaID)
	consulta, err := repo.GetConsulta(consultaID)
	if err != nil {
		return nil, ctx.Err(err)
	}
	paquete, err := repo.GetPaquete(consulta.PaqueteID)
	if err != nil {
		return nil, ctx.Err(err).Op("GetPaqueteDeConsulta")
	}
	tablaFrom, err := getTabla(consulta.TablaID, repo)
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
		joinTbl, err := getTabla(relacion.JoinTablaID, repo)
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
			fromTbl, err := getTabla(relacion.FromTablaID, repo)
			if err != nil {
				return nil, ctx.Err(err)
			}
			rel.From = fromTbl
		}
		agregado.Relaciones = append(agregado.Relaciones, rel)
	}

	return &agregado, nil
}
