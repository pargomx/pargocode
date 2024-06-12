package ddd

import (
	"errors"
)

// ConsultaRelacion corresponde a un elemento de la tabla 'consulta_relaciones'.
type ConsultaRelacion struct {
	ConsultaID  int      // `consulta_relaciones.consulta_id`
	Posicion    int      // `consulta_relaciones.posicion`
	TipoJoin    TipoJoin // `consulta_relaciones.tipo_join`
	JoinTablaID int      // `consulta_relaciones.join_tabla_id`
	JoinAs      string   // `consulta_relaciones.join_as`
	JoinOn      string   // `consulta_relaciones.join_on`
	FromTablaID int      // `consulta_relaciones.from_tabla_id`
}

var (
	ErrConsultaRelacionNotFound      error = errors.New("la consulta relación no se encuentra")
	ErrConsultaRelacionAlreadyExists error = errors.New("la consulta relación ya existe")
)

func (rel *ConsultaRelacion) Validar() error {

	if rel.TipoJoin.EsTodos() {
		return errors.New("ddd.ConsultaRelacion no admite propiedad TipoJoinTodos")
	}

	return nil
}

//  ================================================================  //
//  ========== Tipo de join ========================================  //

// Enumeración
type TipoJoin struct {
	ID          int
	String      string
	Filtro      string
	Descripcion string
}

var (
	// TipoJoinTodos solo se utiliza como filtro.
	TipoJoinTodos = TipoJoin{
		ID:          -1,
		String:      "",
		Filtro:      "todos",
		Descripcion: "Todos los ",
	}
	// Indica explícitamente que la propiedad está indefinida.
	TipoJoinIndefinido = TipoJoin{
		ID:          0,
		String:      "",
		Filtro:      "sin_tipo_join",
		Descripcion: "Indefinido",
	}

	// Inner
	TipoJoinInner = TipoJoin{
		ID:          1,
		String:      "INNER",
		Filtro:      "inner",
		Descripcion: "Inner",
	}
	// Left
	TipoJoinLeft = TipoJoin{
		ID:          2,
		String:      "LEFT",
		Filtro:      "left",
		Descripcion: "Left",
	}
	// Right
	TipoJoinRight = TipoJoin{
		ID:          3,
		String:      "RIGHT",
		Filtro:      "right",
		Descripcion: "Right",
	}
)

// Enumeración excluyendo TipoJoinTodos
var ListaTipoJoin = []TipoJoin{
	TipoJoinIndefinido,

	TipoJoinInner,
	TipoJoinLeft,
	TipoJoinRight,
}

// Enumeración incluyendo TipoJoinTodos
var ListaFiltroTipoJoin = []TipoJoin{
	TipoJoinTodos,
	TipoJoinIndefinido,

	TipoJoinInner,
	TipoJoinLeft,
	TipoJoinRight,
}

// Comparar un Tipo de join con otro.
func (a TipoJoin) Es(e TipoJoin) bool {
	return a.ID == e.ID
}

func (e TipoJoin) EsTodos() bool {
	return e.ID == TipoJoinTodos.ID
}
func (e TipoJoin) EsIndefinido() bool {
	return e.ID == TipoJoinIndefinido.ID
}
func (e TipoJoin) EsInner() bool {
	return e.ID == TipoJoinInner.ID
}
func (e TipoJoin) EsLeft() bool {
	return e.ID == TipoJoinLeft.ID
}
func (e TipoJoin) EsRight() bool {
	return e.ID == TipoJoinRight.ID
}

func (i *ConsultaRelacion) EsTipoJoinTodos() bool {
	return i.TipoJoin.Es(TipoJoinTodos)
}
func (i *ConsultaRelacion) EsTipoJoinIndefinido() bool {
	return i.TipoJoin.Es(TipoJoinIndefinido)
}
func (i *ConsultaRelacion) EsTipoJoinInner() bool {
	return i.TipoJoin.Es(TipoJoinInner)
}
func (i *ConsultaRelacion) EsTipoJoinLeft() bool {
	return i.TipoJoin.Es(TipoJoinLeft)
}
func (i *ConsultaRelacion) EsTipoJoinRight() bool {
	return i.TipoJoin.Es(TipoJoinRight)
}

// Recibe la forma .String
func SetTipoJoinDB(str string) TipoJoin {
	for _, e := range ListaTipoJoin {
		if e.String == str {
			return e
		}
	}
	if str == TipoJoinTodos.String {
		return TipoJoinIndefinido
	}
	return TipoJoinIndefinido
}

// Recibe la forma .Filtro
func SetTipoJoinFiltro(str string) TipoJoin {
	if str == "" || str == TipoJoinTodos.Filtro {
		return TipoJoinTodos
	}
	for _, e := range ListaTipoJoin {
		if e.Filtro == str {
			return e
		}
	}
	return TipoJoinIndefinido
}

// Recibe la forma .String o .Filtro
func (i *ConsultaRelacion) SetTipoJoin(str string) {
	for _, e := range ListaTipoJoin {
		if e.String == str {
			i.TipoJoin = e
			return
		}
		if e.Filtro == str {
			i.TipoJoin = e
			return
		}
	}
	if str == TipoJoinTodos.String {
		i.TipoJoin = TipoJoinIndefinido
	}
	i.TipoJoin = TipoJoinIndefinido
}
