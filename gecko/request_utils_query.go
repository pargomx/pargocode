package gecko

import (
	"errors"
	"time"
)

// ================================================================ //
// ========== QUERY PARAMS ======================================== //

// Valor del query sin sanitizar.
func (c *Context) QueryTalCual(name string) string {
	return c.QueryParam(name)
}

// Valor del query sanitizado.
func (c *Context) QueryVal(name string) string {
	return txtSanitizar(c.QueryParam(name))
}

// Valor del query sanitizado en mayúsculas.
func (c *Context) QueryUpper(name string) string {
	return txtUpper(c.QueryParam(name))
}

// Valor del query sanitizado en minúsculas.
func (c *Context) QueryLower(name string) string {
	return txtLower(c.QueryParam(name))
}

// Valor del query convertido a bool.
// Retorna false a menos de que el valor sea: "on", "true", "1".
func (c *Context) QueryBool(name string) bool {
	return txtBool(c.QueryParam(name))
}

// Valor del query convertido a entero.
func (c *Context) QueryIntMust(name string) (int, error) {
	return txtInt(c.QueryParam(name))
}

// Valor del query convertido a entero sin verificar error (default 0).
func (c *Context) QueryInt(name string) int {
	num, _ := txtInt(c.QueryParam(name))
	return num
}

// Valor del query convertido a uint64.
func (c *Context) QueryUintMust(name string) (uint64, error) {
	return txtUint64(c.QueryParam(name))
}

// Valor del query convertido a uint64 sin verificar error (default 0).
func (c *Context) QueryUint(name string) uint64 {
	num, _ := txtUint64(c.QueryParam(name))
	return num
}

// Valor del query convertido a centavos.
func (c *Context) QueryCentavos(name string) (int, error) {
	return txtCentavos(c.QueryParam(name))
}

// Valor del query convertido a time.
func (c *Context) QueryTime(name string, layout string) (time.Time, error) {
	return txtTime(c.QueryParam(name), layout)
}

// Valor del query convertido a time, que puede estar indefinido.
func (c *Context) QueryTimeNullable(name string, layout string) (*time.Time, error) {
	return txtTimeNullable(c.QueryParam(name), layout)
}

// Valor del query convertido a time desde una fecha 28/08/2022 o 2022-02-13.
func (c *Context) QueryFecha(name string, layout string) (time.Time, error) {
	return txtFecha(c.QueryParam(name))
}

// Valor del path formato fecha convertido a time, que puede estar indefinido.
func (c *Context) QueryFechaNullable(name string) (*time.Time, error) {
	return txtFechaNullable(c.QueryParam(name))
}

// ================================================================ //

// Múltiples valores sin sanitizar obtenidos del query.
func (c *Context) QueryValues(name string) []string {
	if c.query == nil {
		c.query = c.request.URL.Query()
	}
	return c.query[name]
}

// Múltiples valores sanitizados obtenidos del query.
func (c *Context) MultiQueryVal(name string) []string {
	res := []string{}
	for _, v := range c.QueryValues(name) {
		res = append(res, txtSanitizar(v))
	}
	return res
}

// Múltiples valores sanitizados en mayúsculas obtenidos del query.
func (c *Context) MultiQueryUpper(name string) []string {
	res := []string{}
	for _, v := range c.QueryValues(name) {
		res = append(res, txtUpper(v))
	}
	return res
}

// Múltiples valores sanitizados en minúsculas obtenidos del query.
func (c *Context) MultiQueryLower(name string) []string {
	res := []string{}
	for _, v := range c.QueryValues(name) {
		res = append(res, txtLower(v))
	}
	return res
}

// Múltiples valores convertidos a enteros obtenidos del query.
// No se agregan los valores que tengan errores en la conversión.
func (c *Context) MultiQueryInt(name string) []int {
	res := []int{}
	for _, v := range c.QueryValues(name) {
		n, err := txtInt(v)
		if err != nil {
			continue
		}
		res = append(res, n)
	}
	return res
}

// Múltiples valores convertidos a enteros obtenidos del query.
// Los valores deben ser números válidos todos.
func (c *Context) MultiQueryIntMust(name string) ([]int, error) {
	res := []int{}
	for _, v := range c.QueryValues(name) {
		n, err := txtInt(v)
		if err != nil {
			return nil, errors.New("el valor [" + v + "] no es un número válido para [" + name + "]")
		}
		res = append(res, n)
	}
	return res, nil
}

// Múltiples valores convertidos a enteros obtenidos del query.
// No se agregan los valores que tengan errores en la conversión.
func (c *Context) MultiQueryUint(name string) []uint64 {
	res := []uint64{}
	for _, v := range c.QueryValues(name) {
		n, err := txtUint64(v)
		if err != nil {
			continue
		}
		res = append(res, n)
	}
	return res
}

// Múltiples valores convertidos a enteros obtenidos del query.
// Los valores deben ser números válidos todos.
func (c *Context) MultiQueryUintMust(name string) ([]uint64, error) {
	res := []uint64{}
	for _, v := range c.QueryValues(name) {
		n, err := txtUint64(v)
		if err != nil {
			return nil, errors.New("el valor [" + v + "] no es un número válido para [" + name + "]")
		}
		res = append(res, n)
	}
	return res, nil
}

// ================================================================ //

// QueryParamDefault returns the query param value or default
// value for the provided name. Note: it does not distinguish
// if form had no value by that name or value was empty string
func (c *Context) QueryParamDefault(name, defaultValue string) string {
	if c.query == nil {
		c.query = c.request.URL.Query()
	}
	value := c.query.Get(name)
	if value == "" {
		return defaultValue
	}
	return value
}
