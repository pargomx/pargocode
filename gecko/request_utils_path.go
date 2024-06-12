package gecko

import (
	"time"
)

func (c *Context) PathParam(name string) string {
	return c.Param(name)
}

// ================================================================ //
// ========== PATH PARAMS ========================================= //

// Valor del path sin sanitizar.
func (c *Context) PathTalCual(name string) string {
	return c.PathParam(name)
}

// Valor del path sanitizado.
func (c *Context) PathVal(name string) string {
	return txtSanitizar(c.PathParam(name))
}

// Valor del path sanitizado en mayúsculas.
func (c *Context) PathUpper(name string) string {
	return txtUpper(c.PathParam(name))
}

// Valor del path sanitizado en minúsculas.
func (c *Context) PathLower(name string) string {
	return txtLower(c.PathParam(name))
}

// Valor del path convertido a bool.
// Retorna false a menos de que el valor sea: "on", "true", "1".
func (c *Context) PathBool(name string) bool {
	return txtBool(c.PathParam(name))
}

// Valor del path convertido a entero.
func (c *Context) PathIntMust(name string) (int, error) {
	return txtInt(c.PathParam(name))
}

// Valor del path convertido a entero sin verificar error (default 0).
func (c *Context) PathInt(name string) int {
	num, _ := txtInt(c.PathParam(name))
	return num
}

// Valor del path convertido a uint64.
func (c *Context) PathUintMust(name string) (uint64, error) {
	return txtUint64(c.PathParam(name))
}

// Valor del path convertido a uint64 sin verificar error (default 0).
func (c *Context) PathUint(name string) uint64 {
	num, _ := txtUint64(c.PathParam(name))
	return num
}

// Valor del path convertido a centavos.
func (c *Context) PathCentavos(name string) (int, error) {
	return txtCentavos(c.PathParam(name))
}

// Valor del path convertido a time.
func (c *Context) PathTime(name string, layout string) (time.Time, error) {
	return txtTime(c.PathParam(name), layout)
}

// Valor del path convertido a time, que puede estar indefinido.
func (c *Context) PathTimeNullable(name string, layout string) (*time.Time, error) {
	return txtTimeNullable(c.PathParam(name), layout)
}

// Valor del path convertido a time desde una fecha 28/08/2022 o 2022-02-13.
func (c *Context) PathFecha(name string, layout string) (time.Time, error) {
	return txtFecha(c.PathParam(name))
}

// Valor del path formato fecha convertido a time, que puede estar indefinido.
func (c *Context) PathFechaNullable(name string) (*time.Time, error) {
	return txtFechaNullable(c.PathParam(name))
}
