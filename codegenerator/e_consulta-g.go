package codegenerator

import (
	"fmt"
	"strings"

	"github.com/pargomx/gecko/gko"
)

// ================================================================ //
// ========== GENERAR CÓDIGO ====================================== //

// Lista que incluye la abreviacion de la tabla origen del campo.
// Los campos de tablas relacionadas con JOIN que no son nullable
// se seleccionan con un coalesce para evitar porblemas con NULL.
//
// Si es un campo calculado entonces pone la expresión dada.
//
// Ej. usu.UsuarioID, prg.Titulo
func (consulta *Consulta) camposAsSnakeList(campos []CampoConsulta, separador string) (res string) {
	for _, c := range campos {
		switch {
		case c.Expresion == "":
			res += "'campo_sin_expresion'"

		case c.EsNullable():
			res += c.Expresion

		case c.Consulta != nil && c.OrigenTabla != nil && c.OrigenTabla.TablaID == c.Consulta.TablaOrigen.TablaID:
			res += c.Expresion // No necesita coalesce si es de la tabla origen.

		case c.EsNumero():
			res += fmt.Sprintf("coalesce(%s, 0)", c.Expresion)

		default:
			res += fmt.Sprintf("coalesce(%s, '')", c.Expresion)
		}
		if c.AliasSql != "" {
			res += " AS " + c.AliasSql
		}
		res += separador
	}
	return strings.TrimSuffix(res, separador)
}

func (con *Consulta) CamposAsSnakeList(separador string) string {
	return con.camposAsSnakeList(con.Campos, separador)
}

func (consulta *Consulta) SqlGroupClause(separador string) string {
	if !consulta.TieneCamposGroupBy() {
		return ""
	}
	if !strings.Contains(separador, ",") {
		separador = "," + separador
	}
	return "GROUP BY " + consulta.camposAsSnakeList(consulta.CamposGroupBy(), separador)
}

func (consulta *Consulta) SqlFromClause(sep string) string {

	fromSQL := fmt.Sprintf("FROM %v %v", consulta.TablaOrigen.NombreRepo, consulta.TablaOrigen.Abrev)

	for i, rel := range consulta.Relaciones {
		fromSQL += sep + rel.joinString()
		if rel.From == nil {
			gko.LogWarnf("la relación FromTabla = nil")
			continue
		}
		if i == 0 && rel.FromTablaID != consulta.TablaOrigen.TablaID {
			gko.LogWarnf("la consulta '" + consulta.Consulta.NombreItem + "' debería tener el primer join a partir de la tabla principal '" + consulta.TablaOrigen.NombreRepo + "'")
		}
	}
	return fromSQL + sep
}

// Ej. "LEFT JOIN ofertas ofe ON ofe.calendaio_id = insc.calendario_id AND ofe.programa_id = insc.programa_id"
func (r Relacion) joinString() string {
	join := r.TipoJoin.String + " JOIN "
	join += r.Join.Tabla.NombreRepo + " "
	// Alias de la tabla relacionada.
	if r.JoinAs != "" {
		join += r.JoinAs
	} else {
		join += r.Join.Tabla.Abrev
	}
	join += " ON "
	// Poner ON personalizado si se especificó.
	if r.JoinOn != "" {
		return join + r.JoinOn
	}
	// O bien construir ON a partir de campos clave comunes.
	join += r.onString()

	return strings.TrimSuffix(join, " AND ")
}

// Ej. "ofe.calendaio_id = insc.calendario_id AND ofe.programa_id = insc.programa_id"
func (r Relacion) onString() string {
	comparaciones := ""
	for _, tJoin := range r.Join.PrimaryKeys() {
		tFrom, err := r.From.BuscarCampo(tJoin.NombreColumna) // TODO: no usar agregado sino GetCampo(idCampo)
		if err != nil {
			for _, tJoin = range r.Join.ForeignKeys() {
				tFrom, err = r.From.BuscarCampo(tJoin.NombreColumna) // TODO: no usar agregado sino GetCampo(idCampo)
				if err != nil {
					// gko.LogWarnf("Tablas en relación '" + r.TablaJoin.NombreRepo + "->" + r.TablaFrom.NombreRepo + "' no comparten PK-FK " + tJoin.NombreColumna)
					continue // Esto es normal porque la relación puede ser entre un subset de los PK y FK.
				}
				break
			}
		}
		if tFrom == nil {
			gko.LogWarnf("ignorando '" + r.Join.Tabla.NombreRepo + "->" + r.From.Tabla.NombreRepo + "' porque no comparten PK-FK " + tJoin.NombreColumna)
			continue
		}
		if !tJoin.ForeignKey && !tJoin.PrimaryKey {
			gko.LogWarnf("El campo '" + tJoin.NombreColumna + "' de '" + r.Join.Tabla.NombreRepo + "' no está marcado como FK o PK pero se usa en relación " + r.Join.Tabla.NombreRepo + "->" + r.From.Tabla.NombreRepo)
		}
		if !tFrom.ForeignKey && !tFrom.PrimaryKey {
			gko.LogWarnf("El campo '" + tFrom.NombreColumna + "' de '" + r.From.Tabla.NombreRepo + "' no está marcado como FK o PK pero se usa en relación " + r.Join.Tabla.NombreRepo + "->" + r.From.Tabla.NombreRepo)
		}
		comparaciones += fmt.Sprintf("%s.%s = %s.%s AND ",
			r.From.Tabla.Abrev, tFrom.NombreColumna,
			r.Join.Tabla.Abrev, tJoin.NombreColumna,
		)
	}
	return comparaciones
}

// ================================================================ //

func (consulta Consulta) PrimaryKeysAsSqlWhere() (QryWhere string) {
	return CamposAsSqlWhere(consulta.PrimaryKeys(), true)
}
func (consulta Consulta) CamposSeleccionadosAsSqlWhere() (ArgsFunc string) {
	return CamposAsSqlWhere(consulta.CamposSeleccionados, true)
}

func (consulta Consulta) PrimaryKeysAsFuncParams() (ArgsFunc string) {
	return CamposAsFuncParams(consulta.PrimaryKeys())
}
func (consulta Consulta) CamposSeleccionadosAsFuncParams() (ArgsFunc string) {
	return CamposAsFuncParams(consulta.CamposSeleccionados)
}

func (consulta Consulta) PrimaryKeysAsArguments(nombreVariable string) (ArgsWhere string) {
	return CamposAsArguments(consulta.PrimaryKeys(), nombreVariable)
}
func (consulta Consulta) CamposSeleccionadosAsArguments(nombreVariable string) (lista string) {
	return CamposAsArguments(consulta.CamposSeleccionados, nombreVariable)
}

func (consulta Consulta) ScanTempVars() string {
	return ScanTempVars(consulta.Campos)
}
func (consulta Consulta) ScanArgs() string {
	return ScanArgs(consulta.Campos, consulta.Consulta.Abrev)
}
func (consulta Consulta) ScanSetters() string {
	return ScanSetters(consulta.Campos, consulta.Consulta.Abrev)
}
