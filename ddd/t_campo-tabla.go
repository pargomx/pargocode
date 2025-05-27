package ddd

import "strings"

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
