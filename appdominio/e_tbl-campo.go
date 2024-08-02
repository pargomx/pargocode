package appdominio

import (
	"monorepo/ddd"
	"strings"
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

func (c CampoTabla) EsSqlChar() bool {
	return strings.ToUpper(c.TipoSql) == "CHAR"
}

func (c CampoTabla) EsSqlVarchar() bool {
	return strings.ToUpper(c.TipoSql) == "VARCHAR"
}

func (c CampoTabla) EsSqlInt() bool {
	switch c.TipoSql {
	case "tinyint", "smallint", "mediumint", "int", "bigint":
		return true
	default:
		return false
	}
}
