package dpaquete

import (
	"fmt"
	"strings"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ================================================================ //

func (cam CampoConsulta) EsNullable() bool {
	return cam.EsPointer()
}

func (cam CampoConsulta) NotNull() bool {
	return !cam.EsNullable()
}

func (cam CampoConsulta) EsPointer() bool {
	return strings.HasPrefix(cam.TipoGo, "*")
}

func (cam CampoConsulta) EsString() bool {
	return strings.Contains(cam.TipoGo, "string")
}

func (cam CampoConsulta) EsNumero() bool {
	return strings.Contains(cam.TipoGo, "int") ||
		strings.Contains(cam.TipoGo, "float") ||
		cam.EsBool()
}

func (cam CampoConsulta) EsNumeroPositivo() bool {
	return strings.Contains(cam.TipoGo, "uint")
}

func (cam CampoConsulta) EsBool() bool {
	return strings.Contains(cam.TipoGo, "bool")
}

func (con CampoConsulta) EsTiempo() bool {
	return strings.Contains(con.TipoGo, "time.")
}

func (cam CampoConsulta) EsPropiedadExtendida() bool {
	if cam.OrigenCampo == nil {
		return false
	}
	return cam.OrigenCampo.Especial
}

// ================================================================ //
// ================================================================ //

// Ej: "fechaModif", "usuarioID"
func (cam CampoConsulta) Variable() string {
	return strings.ToLower(cam.NombreCampo[:1]) + cam.NombreCampo[1:]
}

func (cam CampoConsulta) IfZeroReturnErr(razón string, nombreVariable string) string {
	return cam.ifZeroReturnErr(razón, nombreVariable, false)
}

func (cam CampoConsulta) IfZeroReturnNilAndErr(razón string, nombreVariable string) string {
	return cam.ifZeroReturnErr(razón, nombreVariable, true)
}

// Crea un snippet de código que verifica que el valor de un campo requerido no sea el zero value.
//
// Ejemplo de resultado:
//
//	if enc.OrganizacionID == 0 {
//		return nil, gecko.NewErr(http.StatusBadRequest).Msg("OrganizacionID sin especificar").Ctx(op, "pk_indefinida")
//	}
//
// razón que se da como contexto al error. Ejemplos: "pk_indefinida" "fk requerida" "campo requerido"
//
// returnNilErr nil cuando la funciona retorna dos valores y el error es el segundo, por ejemplo:
//
//	return nil, errors.New("valor indefinido")
//
// nombreVariable cuando el campo es parte de una struct:
//
//	if item.Campo = "" { ... }
//
// si nombreVariable se deja vacío entonces se usa el campo como variable:
//
//	if Campo = "" { ... }
func (cam CampoConsulta) ifZeroReturnErr(razón string, nombreVariable string, returnNilErr bool) string {

	comparacion := "\tif "
	if nombreVariable != "" { // ej. "if [org].OrganizacionID"
		comparacion += nombreVariable + "."
	}
	comparacion += cam.NombreCampo // ej. "if org.[OrganizacionID]"

	switch {

	case cam.EsNullable():
		comparacion += " == nil "

	case cam.EsString():
		comparacion += ` == "" `

	case cam.EsNumero():
		comparacion += " == 0 "

	case cam.EsTiempo():
		comparacion += ".IsZero() "

	case cam.EsPropiedadExtendida():
		comparacion += ".Es" + cam.NombreCampo + "Indefinido() "

	default:
		gecko.LogWarnf("cam.NombreCampo: No se verificará que %v no sea Zero value", cam.NombreCampo)
		return `\\` + " TODO: verificar que " + cam.NombreCampo + " no esté indefinido"
	}

	comparacion += " {\n" + "\t\t" + "return "

	if returnNilErr {
		comparacion += "nil, "
	}

	comparacion += fmt.Sprintf(`gecko.NewErr(http.StatusBadRequest).Msg("%v sin especificar").Ctx(op, "%v")`, cam.NombreCampo, razón)

	comparacion += "\n}\n"

	return comparacion
}

// ================================================================ //
// ================================================================ //

// SQL: Campos como WHERE clause. Por ejemplo:
//
//	"WHERE usu.usuario_id = ?"
//	"WHERE ins.programa_id = ? AND prg.calendario_id = ?"
func CamposAsSqlWhere(campos []CampoConsulta, incluirTabla bool) string {
	where := "WHERE "
	for _, campo := range campos {
		where += campo.Expresion + " = ? AND "
	}
	return strings.TrimSuffix(where, " AND ")
}

// Go code: PKs como argumentos recibidos en la función. Por ejemplo:
//
//	"UsuarioID int, ProgramaID string"
func CamposAsFuncParams(campos []CampoConsulta) (s string) {
	for _, campo := range campos {
		s += campo.NombreCampo + " " + campo.TipoGo + ", "
	}
	return strings.TrimSuffix(s, ", ")
}

// Como argumentos para llamar una función. Por ejemplo:
//
//	"UsuarioID, ProgramaID"
//
// Si nombreVariable está definida, se usa como prefijo:
//
//	"apr.UsuarioID, apr.ProgramaID"
func CamposAsArguments(campos []CampoConsulta, nombreVariable string) (s string) {
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

// ================================================================ //

// 1. Variables para que rows.Scan() pueda colocar los valores obtenidos de la base de datos.
func ScanTempVars(campos []CampoConsulta) string {
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
			gecko.LogWarnf("el campo " + campo.NombreCampo + " no puede ser " + campo.TipoGo + " para generar SQL")
		}
	}
	return strings.TrimPrefix(res, "\n\t")
}

// ================================================================ //

// 2. Argumentos para llamar a rows.Scan() o row.Scan() en forma de pointers.
//
// El itemVar es el nombre de la variable de la estructura. Ej. "usu" para resultar en &usu.UsuarioID, &usu.Nombre
func ScanArgs(campos []CampoConsulta, itemVar string) string {
	if itemVar == "" {
		gecko.LogWarnf("itemVar indefinida para ScanArgs")
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
			gecko.LogWarnf("el campo " + campo.NombreCampo + " no puede ser " + campo.TipoGo + " para generar SQL")

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
func ScanSetters(campos []CampoConsulta, itemVar string) string {
	if itemVar == "" {
		gecko.LogWarnf("itemVar indefinida para ScanSetters")
	}
	var res string
	for _, c := range campos {
		res += "\n\t"
		switch {

		case c.EsPropiedadExtendida(): // ej. usu.TipoGo = usuario.SetTipoDB(tipo)

			res += itemVar + "." + c.NombreCampo + " = " + c.OrigenPaquete.Nombre + ".Set" + c.TipoGo + "DB(" + c.Variable() + ")"

			// ================================================================ //

		// case c.TipoImportado && c.TipoSetter != "": // ej. usu.TipoImportado = importado.SetTipoDB(tipo)
		// gecko.LogWarnf("Usando TipoImportado no implementado")
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
			// 	gecko.LogWarnf("el campo " + c.NombreCampo + " es time.Time pero no se sabe si timestamp|datetime|date|time")
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
					"if %v.Int64 < 0{\n gecko.LogWarnf(fmt.Sprint(\"el campo %v espera número positivo pero obtuvo \",%v.Int64)) \n}\n"+
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
			gecko.LogWarnf("el campo " + c.NombreCampo + " no puede ser " + c.TipoGo + " para generar SQL")
			// ================================================================ //

		default:
			res = strings.TrimSuffix(res, "\n\t")
			// ================================================================ //
		}
	}
	return strings.TrimPrefix(res, "\n\t")
}
