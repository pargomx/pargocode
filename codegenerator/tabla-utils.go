package codegenerator

import (
	"errors"
	"fmt"
	"monorepo/ddd"
	"strings"

	"github.com/pargomx/gecko/gko"
)

// ================================================================ //
// ================================================================ //

func (tbl *Tabla) NombreRepo() string {
	return tbl.Tabla.NombreRepo
}
func (tbl *Tabla) NombreItem() string {
	return tbl.Tabla.NombreItem
}
func (tbl *Tabla) NombreItems() string {
	return tbl.Tabla.NombreItems
}
func (tbl *Tabla) NombreAbrev() string {
	return tbl.Tabla.Abrev
}
func (tbl *Tabla) NombreNominativo() string {
	if tbl.Tabla.EsFemenino {
		return "la " + strings.ToLower(tbl.Tabla.Humano)
	}
	return "el " + strings.ToLower(tbl.Tabla.Humano)
}

func (tbl *Tabla) Directrices() []Directriz {
	return ToDirectrices(tbl.Tabla.Directrices)
}

func (tbl *Tabla) TablaOrigen() ddd.Tabla {
	return tbl.Tabla
}

// Devuelve una compia del campo si la tabla lo tiene.
// Evalúa nombre Kebab, Snake, Camel, Humano.
func (tbl *Tabla) BuscarCampo(nombre string) (*CampoTabla, error) {
	if nombre == "" {
		return nil, errors.New("no se puede buscar campo sin especificar nombre")
	}
	campo := CampoTabla{}
	for _, c := range tbl.Campos {
		if nombre == c.Campo.NombreCampo ||
			nombre == c.Campo.NombreColumna ||
			nombre == c.Campo.NombreHumano {
			campo = c
			return &campo, nil
		}
	}
	return nil, errors.New("la tabla '" + tbl.Tabla.NombreRepo + "' no tiene campo '" + nombre + "'")
}

func (tbl *Tabla) PrimerCampo() CampoTabla {
	return tbl.Campos[0]
}

func (tbl *Tabla) TieneCamposFiltro() bool {
	for _, campo := range tbl.Campos {
		if campo.Campo.Filtro {
			return true
		}
	}
	return false
}

func (tbl *Tabla) CamposFiltro() []CampoTabla {
	campos := []CampoTabla{}
	for _, campo := range tbl.Campos {
		if campo.Campo.Filtro {
			campos = append(campos, campo)
		}
	}
	return campos
}

// Retorna una copia de todos los campos PK.
func (tbl *Tabla) PrimaryKeys() []CampoTabla {
	res := []CampoTabla{}
	for _, c := range tbl.Campos {
		if c.Campo.PrimaryKey {
			res = append(res, c)
		}
	}
	return res
}

// Retorna una copia de todos los campos FK.
func (tbl *Tabla) ForeignKeys() []CampoTabla {
	res := []CampoTabla{}
	for _, c := range tbl.Campos {
		if c.Campo.ForeignKey {
			res = append(res, c)
		}
	}
	return res
}

func (tbl *Tabla) UniqueKeys() []CampoTabla {
	res := []CampoTabla{}
	for _, c := range tbl.Campos {
		if c.Campo.Unique() {
			res = append(res, c)
		}
	}
	return res
}

func (tbl *Tabla) CamposEspeciales() []CampoTabla {
	res := []CampoTabla{}
	for _, c := range tbl.Campos {
		if c.Campo.Especial {
			res = append(res, c)
		}
	}
	return res
}

func (tbl *Tabla) CamposEditables() []CampoTabla {
	res := []CampoTabla{}
	for _, c := range tbl.Campos {
		if !c.Campo.ReadOnly() {
			res = append(res, c)
		}
	}
	return res
}

// Retorna los campos que son requeridos para escribir (required y PKs).
func (tbl *Tabla) CamposRequeridosOrPK() []CampoTabla {
	res := []CampoTabla{}
	for _, c := range tbl.Campos {
		if c.Required() || c.PrimaryKey {
			res = append(res, c)
		}
	}
	return res
}

// ================================================================ //
// ========== LISTA DE CAMPOS COMO STRING ========================= //

// Campos en forma snake separados por coma. Ejemplo:
//
//	"usuario_id, programa_id, estatus"
func TablaCamposAsSnakeList(campos []CampoTabla, separador string) (s string) {
	for _, campo := range campos {
		s += campo.NombreColumna + separador
	}
	return strings.TrimSuffix(s, separador)
}

func (tbl Tabla) CamposAsSnakeList(separador string) (lista string) {
	return TablaCamposAsSnakeList(tbl.Campos, separador)
}
func (tbl Tabla) CamposEditablesAsSnakeList(separador string) (lista string) {
	return TablaCamposAsSnakeList(tbl.CamposEditables(), separador)
}

// ================================================================ //

// Campos en forma snake = ? separados por coma. Ejemplo:
//
//	"usuario_id = ?, programa_id = ?, estatus = ?"
func CamposAsSnakeEqPlaceholder(campos []CampoTabla) (s string) {
	for _, campo := range campos {
		s += campo.NombreColumna + "=?, "
	}
	return strings.TrimSuffix(s, ", ")
}

func (tbl Tabla) CamposEditablesAsSnakeEqPlaceholder() (lista string) {
	return CamposAsSnakeEqPlaceholder(tbl.CamposEditables())
}

// ================================================================ //

// Par cada campo pone un placeholder "?" separado por comas. Por ejemplo:
//
//	"?, ?, ?"
func CamposAsPlaceholders(campos []CampoTabla) (s string) {
	for range campos {
		s += "?, "
	}
	return strings.TrimSuffix(s, ", ")
}

func (tbl Tabla) CamposEditablesAsPlaceholders() (lista string) {
	return CamposAsPlaceholders(tbl.CamposEditables())
}

// ================================================================ //

// SQL: Campos como WHERE clause. Por ejemplo:
//
//	"WHERE usuario_id = ?"
//	"WHERE programa_id = ? AND calendario_id = ?"
//
// Si se quiere inculir la tabla:
//
//	"WHERE usu.usuario_id = ?"
//	"WHERE insc.programa_id = ? AND insc.calendario_id = ?"
func CamposTablaAsSqlWhere(campos []CampoTabla, incluirTabla bool) (s string) {
	for _, campo := range campos {
		if incluirTabla {
			if s == "" {
				s = "WHERE " + campo.Tabla.Abrev + "." + campo.NombreColumna + " = ?"
			} else {
				s = s + " AND " + campo.Tabla.Abrev + "." + campo.NombreColumna + " = ?"
			}
		} else {
			if s == "" {
				s = "WHERE " + campo.NombreColumna + " = ?"
			} else {
				s = s + " AND " + campo.NombreColumna + " = ?"
			}
		}
	}
	return s
}

func (tbl Tabla) PrimaryKeysAsSqlWhere() (QryWhere string) {
	return CamposTablaAsSqlWhere(tbl.PrimaryKeys(), false)
}
func (tbl Tabla) CamposSeleccionadosAsSqlWhere() (ArgsFunc string) {
	return CamposTablaAsSqlWhere(tbl.CamposSeleccionados, false)
}

// ================================================================ //

// Go code: PKs como argumentos recibidos en la función. Por ejemplo:
//
//	"UsuarioID int, ProgramaID string"
func CamposTablaAsFuncParams(campos []CampoTabla) (s string) {
	for _, campo := range campos {
		s += campo.NombreCampo + " " + campo.TipoGo + ", "
	}
	return strings.TrimSuffix(s, ", ")
}

func (tbl Tabla) PrimaryKeysAsFuncParams() (ArgsFunc string) {
	return CamposTablaAsFuncParams(tbl.PrimaryKeys())
}
func (tbl Tabla) CamposSeleccionadosAsFuncParams() (ArgsFunc string) {
	return CamposTablaAsFuncParams(tbl.CamposSeleccionados)
}

// ================================================================ //

// Como argumentos para llamar una función. Por ejemplo:
//
//	"UsuarioID, ProgramaID"
//
// Si nombreVariable está definida, se usa como prefijo:
//
//	"apr.UsuarioID, apr.ProgramaID"
func CamposTablaAsArguments(campos []CampoTabla, nombreVariable string) (s string) {
	if nombreVariable != "" {
		for _, campo := range campos {
			if campo.EsPropiedadExtendida() {
				s += nombreVariable + "." + campo.NombreCampo + ".String, "
			} else {
				s += nombreVariable + "." + campo.NombreCampo + ", "
			}
		}
	} else {
		for _, campo := range campos {
			if campo.EsPropiedadExtendida() {
				s += campo.NombreCampo + ".String, "
			} else {
				s += campo.NombreCampo + ", "
			}
		}
	}
	return strings.TrimSuffix(s, ", ")
}

func (tbl Tabla) PrimaryKeysAsArguments(nombreVariable string) (ArgsWhere string) {
	return CamposTablaAsArguments(tbl.PrimaryKeys(), nombreVariable)
}
func (tbl Tabla) CamposEditablesAsArguments(nombreVariable string) (lista string) {
	return CamposTablaAsArguments(tbl.CamposEditables(), nombreVariable)
}
func (tbl Tabla) CamposSeleccionadosAsArguments(nombreVariable string) (lista string) {
	return CamposTablaAsArguments(tbl.CamposSeleccionados, nombreVariable)
}

// ================================================================ //
// ================================================================ //

func (tbl *Tabla) SqlFromClause(separador string) string {
	return "FROM " + tbl.Tabla.NombreRepo + " "
}

// compatibilidad con mysql-scan
func (tbl *Tabla) SqlGroupClause(separador string) string {
	return ""
}

// ================================================================ //
// ========== Scan ================================================ //

func (tbl *Tabla) ScanTempVars() string {
	return ScanTempVarsTabla(tbl.Campos)
}
func (tbl *Tabla) ScanArgs() string {
	return ScanArgsTabla(tbl.Campos, tbl.Tabla.Abrev)
}
func (tbl *Tabla) ScanSetters() string {
	return ScanSettersTabla(tbl.Campos, tbl.Tabla.Abrev)
}

// 1. Variables para que rows.Scan() pueda colocar los valores obtenidos de la base de datos.
func ScanTempVarsTabla(campos []CampoTabla) string {
	res := ""
	for _, campo := range campos {
		switch {
		case campo.EsPropiedadExtendida(): //* Propiedad extendida
			res += "\n\tvar " + campo.Variable() + " string" // ej. var tipo string

		// case campo.TipoImportado && campo.TipoSetter != "":
		// res += "\n\tvar " + campo.Variable() + " string" // ej. var tipoImportado string

		case campo.TipoGo == "*time.Time" && campo.EsPointer():
			res += "\n\tvar " + campo.Variable() + " sql.NullTime" // ej. var fechaModif sql.NullTime

		case campo.TipoGo == "time.Duration":
			res += "\n\tvar " + campo.Variable() + " sql.NullString" // ej. var duracion sql.NullString

		case campo.EsNumeroPositivo() && campo.EsPointer():
			res += "\n\tvar " + campo.Variable() + " sql.NullInt64" // ej. var calificacion sql.NullInt64

		case campo.EsNumero() && campo.EsPointer():
			res += "\n\tvar " + campo.Variable() + " sql.NullInt64" // ej. var calificacion sql.NullInt64

		case campo.EsPointer() && campo.EsString() && campo.EsNullable():
			res += "\n\tvar " + campo.Variable() + " sql.NullString" // ej. var matricula sql.NullString

		case campo.EsPointer(): //* No reconocido
			res += "\n\tvar " + "invalid string // No reconocido"
			gko.LogWarnf("el campo " + campo.NombreCampo + " no puede ser " + campo.TipoGo + " para generar SQL")
		}
	}
	return strings.TrimPrefix(res, "\n\t")
}

// ================================================================ //

// 2. Argumentos para llamar a rows.Scan() o row.Scan() en forma de pointers.
//
// El itemVar es el nombre de la variable de la estructura. Ej. "usu" para resultar en &usu.UsuarioID, &usu.Nombre
func ScanArgsTabla(campos []CampoTabla, itemVar string) string {
	if itemVar == "" {
		gko.LogWarnf("itemVar indefinida para ScanArgs")
	}
	var args string
	for _, campo := range campos {
		args += "&"
		switch {
		case campo.EsPropiedadExtendida():
			args += campo.Variable() // ej. &tipo

		// case campo.TipoImportado && campo.TipoSetter != "":
		// args += campo.Variable() // ej. &tipoImportado

		case campo.TipoGo == "time.Time" && !campo.EsPointer():
			args += itemVar + "." + campo.NombreCampo // ej. &usu.FechaRegistro. // time.Time es soportado por el driver si se pone "?parseTime=true"

		case campo.TipoGo == "*time.Time" && campo.EsPointer():
			args += campo.Variable() // ej. &fechaModif

		case campo.TipoGo == "time.Duration":
			args += campo.Variable() // ej. &duracion

		// case campo.EsNumeroPositivo() && campo.EsPointer():
		//	args += campo.Variable() // ej. &calificacion

		case campo.EsNumero() && campo.EsPointer():
			args += campo.Variable() // ej. &calificacion

		case campo.EsPointer() && campo.EsString() && campo.EsNullable():
			args += campo.Variable() // ej. &matricula

		case campo.EsPointer() && campo.EsString(): // Si viene null será "".
			args += itemVar + "." + campo.NombreCampo

		case campo.EsPointer():
			args += "Invalid" // No reconocido
			gko.LogWarnf("el campo " + campo.NombreCampo + " no puede ser " + campo.TipoGo + " para generar SQL")

		default:
			args += itemVar + "." + campo.NombreCampo // ej. &usu.UsuarioID, &usu.Nombre... (int, string)
		}
		args += ", "
	}
	return strings.TrimSuffix(args, ", ")
}

// ================================================================ //

// 3. Tasnformar variables temporales que usa row.Scan para ponerlas en el Item.
//
// El itemVar es el nombre de la variable de la estructura. Ej. "usu" para resultar en usu.SetEstatusDB(...), usu.FechaBaja = (...)
func ScanSettersTabla(campos []CampoTabla, itemVar string) string {
	if itemVar == "" {
		gko.LogWarnf("itemVar indefinida para ScanSetters")
	}
	var res string
	for _, c := range campos {
		res += "\n\t"
		switch {

		case c.EsPropiedadExtendida(): // ej. usu.TipoGo = usuario.SetTipoDB(tipo)

			res += itemVar + "." + c.NombreCampo + " = " + c.Paquete.Nombre + ".Set" + c.TipoGo + "DB(" + c.Variable() + ")"

			// ================================================================ //

		// case c.TipoImportado && c.TipoSetter != "": // ej. usu.TipoImportado = importado.SetTipoDB(tipo)
		// gko.LogWarnf("Usando TipoImportado no implementado")
		// res += itemVar + "." + c.NombreCampo + " = " + strings.ReplaceAll(c.TipoSetter, "?", c.Variable())
		// ================================================================ //

		case c.TipoGo == "*time.Time" && c.EsPointer(): // ej. if fechaBaja.Valid { apr.FechaBaja = &fechaBaja.Time }

			// switch c.TimeTipo() {
			// case "datetime", "timestamp", "date", "time":
			// 	res += fmt.Sprintf(
			// 		"if %v.Valid {\n\t%v.%v = &%v.Time\n}",
			// 		c.Variable(), itemVar, c.NombreCampo, c.Variable(),
			// 	)
			// default:
			// 	res += "invalid"
			// 	gko.LogWarnf("el campo " + c.NombreCampo + " es time.Time pero no se sabe si timestamp|datetime|date|time")
			// }
			// ================================================================ //

		case c.TipoGo == "time.Duration":
			x := c.Variable() + ".String"
			res += fmt.Sprintf(
				`// time.Duration
	switch len(%v) {
	case 0: // ej. "" (NULL)
		%v = "0"
	case 8: // ej. 02:15:59
		%v = %v[0:2] + "h" + %v[3:5] + "m" + %v[6:8] + "s"
	case 9: // ej. 126:15:59
		%v = %v[0:3] + "h" + %v[4:6] + "m" + %v[7:9] + "s"
	}
	%v.%v, err = time.ParseDuration(%v)
	if err != nil {
		fmt.Println(err)
	}`, x, x, x, x, x, x, x, x, x, x, itemVar, c.NombreCampo, x)
			// ================================================================ //

		case c.EsNumeroPositivo() && c.EsPointer():

			res += fmt.Sprintf(
				"\n if %v.Valid{ \n\t\t"+
					"if %v.Int64 < 0{\n gko.LogWarnf(fmt.Sprint(\"el campo %v espera número positivo pero obtuvo \",%v.Int64)) \n}\n"+
					"num := %v(%v.Int64) \n\t\t"+ // ej. if calificacion.Valid {
					"%v.%v = &num \n}", // 				num := int(calificacion.Int64)
				c.Variable(), // 					apr.Calificacion = &num
				c.Variable(), c.NombreCampo, c.Variable(),
				c.TipoGo[1:], c.Variable(),
				itemVar, c.NombreCampo,
			)
			// ================================================================ //

		case c.EsNumero() && c.EsPointer():

			res += fmt.Sprintf(
				"\n if %v.Valid{ \n\t\t"+
					"num := %v(%v.Int64) \n\t\t"+
					"%v.%v = &num \n}",
				c.Variable(),               // ej. if calificacion.Valid {
				c.TipoGo[1:], c.Variable(), // 		num := int(calificacion.Int64)
				itemVar, c.NombreCampo, //	apr.Calificacion = &num
			)
			// ================================================================ //

		case c.EsPointer() && c.EsString() && c.EsNullable():

			res += fmt.Sprintf(
				"\n if %v.Valid{ \n\t\t"+ //       ej. if matricula.Valid {
					"%v.%v = &%v.String \n}", // 	      apr.Matricula = &matricula.String
				c.Variable(),
				itemVar, c.NombreCampo, c.Variable(),
			)
			// ================================================================ //

		case c.EsPointer():

			res += "invalid"
			gko.LogWarnf("el campo " + c.NombreCampo + " no puede ser " + c.TipoGo + " para generar SQL")
			// ================================================================ //

		default:
			res = strings.TrimSuffix(res, "\n\t")
			// ================================================================ //
		}
	}
	return strings.TrimPrefix(res, "\n\t")
}
