package dpaquete

import (
	"fmt"
	"monorepo/domain_driven_design/ddd"
	"strings"

	"github.com/pargomx/gecko/gko"
)

type CampoTabla struct {
	Paquete         *ddd.Paquete
	Tabla           *ddd.Tabla
	Campo           ddd.Campo
	ValoresPosibles []ddd.ValorEnum

	TablaFK *ddd.Tabla
	CampoFK *ddd.Campo

	CampoID         int    // `campos.campo_id`
	TablaID         int    // `campos.tabla_id`
	NombreCampo     string // `campos.nombre_campo`  Nombre camel del campo en Go
	NombreColumna   string // `campos.nombre_columna`  Nombre de la columna SQL
	NombreHumano    string // `campos.nombre_humano`  Nombre para poner en etiquetas de la UI
	TipoGo          string // `campos.tipo_go`  Tipo de dato en código. Ej. 'int' | 'Estatus' | '*time.Time'
	TipoSql         string // `campos.tipo_sql`  Tipo de dato en SQL. Ej. 'MEDIUMINT', 'VARCHAR', 'TIMESTAMP'.
	Setter          string // `campos.setter`  ej. `usuario.SetEstatusDB(?)` donde '?' es la variable tipo string hacia la que se hace el scan en mysql.
	Importado       bool   // `campos.importado`  Si el tipo se importa de otro paquete
	PrimaryKey      bool   // `campos.primary_key`  Clave primaria: 'pk'
	ForeignKey      bool   // `campos.foreign_key`  Clave foránea: 'fk'
	Uq              bool   // `campos.uq`  Clave natural: 'unique' (ej. UUID, clave, matrícula, correo...)
	Req             bool   // `campos.req`  Requerida para update e insert: 'required'
	Ro              bool   // `campos.ro`  No se guardará ni actualizará su valor en el repositorio: 'readonly'
	Filtro          bool   // `campos.filtro`  Se debe generar lógica para filtrar por este campo desde el repositorio: 'filtro'
	Nullable        bool   // `campos.nullable`  Declara que el campo puede ser nulo: 'null' (ej *int | *string | *time.Time )
	MaxLenght       int    // `campos.max_lenght`  Longitud máxima permitida en repositorio para strings: 'max' [ej. varchar(12), char(3)].
	Uns             bool   // `campos.uns`  Valor numérico positivo.
	DefaultSql      string // `campos.default_sql`  Valor por default. Puede ser: '', 'NULL', 'AUTO_INCREMENT', 'OTRO'.
	Especial        bool   // `campos.especial`  Propiedad extendida o especial funciona como ENUM dentro del código de la aplicación.
	ReferenciaCampo *int   // `campos.referencia_campo`  Campo al que referencia FK
	Expresion       string // `campos.expresion`  Expresión SQL para calcular el campo
	EsFemenino      bool   // `campos.es_femenino`
	Descripcion     string // `campos.descripcion`
	Posicion        int    // `campos.posicion`
}

func (c CampoTabla) DefaultSQL() string {
	return c.DefaultSql
}

func (c CampoTabla) TipoSQL() string {
	return c.TipoSql
}

func (c CampoTabla) Unsigned() bool {
	return c.Uns
}

func (c CampoTabla) Unique() bool {
	return c.Uq
}

func (c CampoTabla) Required() bool {
	return c.Req
}

func (c CampoTabla) ReadOnly() bool {
	return c.Ro
}

func (c CampoTabla) Null() bool {
	return c.Nullable
}

func (c CampoTabla) EsCalculado() bool {
	return c.Expresion != ""
}

func (c CampoTabla) EsSqlChar() bool {
	return strings.ToUpper(c.TipoSql) == "CHAR"
}

func (c CampoTabla) EsSqlVarchar() bool {
	return strings.ToUpper(c.TipoSql) == "VARCHAR"
}

func (c CampoTabla) EsSqlText() bool {
	return strings.Contains(strings.ToUpper(c.TipoSql), "TEXT")
}

func (c CampoTabla) EsSqlInt() bool {
	switch c.TipoSql {
	case "tinyint", "smallint", "mediumint", "int", "bigint":
		return true
	default:
		return false
	}
}

func (c CampoTabla) EsBool() bool {
	return c.TipoGo == "bool"
}

// Retorna true si el campo termina en _id.
func (c CampoTabla) EsID() bool {
	return strings.HasSuffix(c.NombreColumna, "_id")
}

// Retorna si es una clave de uuid.
func (c CampoTabla) EsUUID() bool {
	return strings.Contains(c.NombreColumna, "uuid")
}

// Retorna si es un número que se puede comparar con 0
// pero que solo sea positivo.
func (c CampoTabla) EsNumeroPositivo() bool {
	return strings.Contains(c.TipoGo, "uint")
}

// Retorna si es un número que se puede comparar con 0.
func (c CampoTabla) EsNumero() bool {
	return strings.Contains(c.TipoGo, "int")
}

// Retorna si es un string que se puede comparar con "".
func (c CampoTabla) EsString() bool {
	return strings.Contains(c.TipoGo, "string")
}

// Retorna si es time.Time
func (c CampoTabla) EsTiempo() bool {
	return c.TipoGo == "time.Time"
}

// Compara si el tipo definido en el modelo
// coincide con el nombre de una propiedad
// extendida declarada en la semilla.csv

func (c CampoTabla) EsPropiedadExtendida() bool {
	return c.Especial
}

// Retorna true si el tipo del campo comienza con "*".
// ej. *int, *string.
func (c CampoTabla) EsPointer() bool {
	if c.TipoGo == "" {
		return false
	}
	return c.TipoGo[:1] == "*"
}

// Verifica si el texto coincide con el
// nombre de una propiedad extendida
// declarada en la semilla.csv
// Útil para discriminar declaraciones de types.
func EsPropiedadExtendida(nombre string) bool {
	gko.FatalExit("EsPropiedadExtendida: deprecated")
	return false
}

func (cam CampoTabla) EsNullable() bool {
	return cam.EsPointer()
}

func (cam CampoTabla) NotNull() bool {
	return !cam.EsNullable()
}

// func (c CampoTabla) NotNull() bool {
// 	return !c.Nullable
// }

// ================================================================ //
// ========== Code gen ============================================ //

// Retorna el nombre en forma de variable de go.
// ej. fechaModif
func (c CampoTabla) Variable() string {
	return strings.ToLower(c.NombreCampo[:1]) + c.NombreCampo[1:]
}

func (c CampoTabla) IfZeroReturnErr(razón string, nombreVariable string) string {
	return c.ifZeroReturnErr(razón, nombreVariable, false)
}

func (c CampoTabla) IfZeroReturnNilAndErr(razón string, nombreVariable string) string {
	return c.ifZeroReturnErr(razón, nombreVariable, true)
}

// Crea un snippet de código que verifica que el valor de un campo requerido no sea el zero value.
//
// Ejemplo de resultado:
//
//	if enc.OrganizacionID == 0 {
//		return nil, gko.ErrDatoInvalido().Msg("OrganizacionID sin especificar").Ctx(op, "pk_indefinida")
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
func (c CampoTabla) ifZeroReturnErr(razón string, nombreVariable string, returnNilErr bool) string {

	comparacion := "\tif "
	if nombreVariable != "" { // ej. "if [org].OrganizacionID"
		comparacion += nombreVariable + "."
	}
	comparacion += c.NombreCampo // ej. "if org.[OrganizacionID]"

	switch {

	case c.Nullable:
		comparacion += " == nil "

	case c.EsString():
		comparacion += ` == "" `

	case c.EsNumero():
		comparacion += " == 0 "

	case c.EsTiempo():
		comparacion += ".IsZero() "

	case c.EsPropiedadExtendida():
		comparacion += ".Es" + c.NombreCampo + "Indefinido() "

	default:
		gko.LogWarnf("No se verificará que %v no sea Zero value", c.NombreCampo)
		return `\\` + " TODO: verificar que " + c.NombreCampo + " no esté indefinido"
	}

	comparacion += " {\n" + "\t\t" + "return "

	if returnNilErr {
		comparacion += "nil, "
	}

	comparacion += fmt.Sprintf(`gko.ErrDatoInvalido().Msg("%v sin especificar").Ctx(op, "%v")`, c.NombreCampo, razón)

	comparacion += "\n}\n"

	return comparacion
}
