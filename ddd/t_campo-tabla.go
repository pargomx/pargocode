package ddd

import (
	"errors"
	"strings"
)

// Campo corresponde a un elemento de la tabla 'campos'.
type Campo struct {
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

var (
	ErrCampoNotFound      error = errors.New("el campo no se encuentra")
	ErrCampoAlreadyExists error = errors.New("el campo ya existe")
)

func (cam *Campo) Validar() error {

	return nil
}

func (c Campo) DefaultSQL() string {
	return c.DefaultSql
}

func (c Campo) TipoSQL() string {
	return c.TipoSql
}

func (c Campo) Unsigned() bool {
	return c.Uns
}

func (c Campo) Unique() bool {
	return c.Uq
}

func (c Campo) Required() bool {
	return c.Req
}

func (c Campo) ReadOnly() bool {
	return c.Ro
}

func (c Campo) Null() bool {
	return c.Nullable
}

func (c Campo) NotNull() bool {
	return !c.Nullable
}

func (c Campo) EsCalculado() bool {
	return c.Expresion != ""
}

func (c Campo) EsSqlChar() bool {
	return strings.ToUpper(c.TipoSql) == "CHAR"
}

func (c Campo) EsSqlVarchar() bool {
	return strings.ToUpper(c.TipoSql) == "VARCHAR"
}

func (c Campo) EsSqlText() bool {
	return strings.Contains(strings.ToUpper(c.TipoSql), "TEXT")
}

func (c Campo) EsSqlInt() bool {
	switch c.TipoSql {
	case "tinyint", "smallint", "mediumint", "int", "bigint":
		return true
	default:
		return false
	}
}
