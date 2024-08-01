package ddd

import "errors"

// ConsultaCampo corresponde a un elemento de la tabla 'consulta_campos'.
type ConsultaCampo struct {
	ConsultaID  int    // `consulta_campos.consulta_id`
	Posicion    int    // `consulta_campos.posicion`
	Expresion   string // `consulta_campos.expresion`  Columna o expresión calculada sin incluir alias Ejs. usu.usuario_id, SUM(pag.monto) (SELECT nombre FROM x WHERE y = z).
	AliasSql    string // `consulta_campos.alias_sql`  Alias que poner a la columna en la consulta por ejemplo para no repetir nombres en el resultset.
	NombreCampo string // `consulta_campos.nombre_campo`  Nombre que poner al campo en la estructura de Go
	TipoGo      string // `consulta_campos.tipo_go`  Tipo Go tal cual aparecerá en la estructura, incluyendo paquete de importación si se necesita.
	CampoID     *int   // `consulta_campos.campo_id`  Referencia al campo de la tabla de donde proviene el campo, lo que indica que no es un campo calculado.
	Pk          bool   // `consulta_campos.pk`  Si es parte de la clave primaria en el contexto de la consulta
	Filtro      bool   // `consulta_campos.filtro`  Si se puede filtrar por este campo
	GroupBy     bool   // `consulta_campos.group_by`  Si se usa este campo para agregar los resultados
	Descripcion string // `consulta_campos.descripcion`
}

var (
	ErrConsultaCampoNotFound      error = errors.New("el campo de consulta no se encuentra")
	ErrConsultaCampoAlreadyExists error = errors.New("el campo de consulta ya existe")
)

func (campcons *ConsultaCampo) Validar() error {

	return nil
}
