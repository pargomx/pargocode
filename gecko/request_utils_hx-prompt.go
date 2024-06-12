package gecko

import (
	"time"
)

// ================================================================ //
// ========== HX-Prompt Header ==================================== //

// Valor del header Hx-Prompt sin sanitizar.
func (c *Context) PromptTalCual() string {
	return c.Request().Header.Get("HX-Prompt")
}

// Valor del header Hx-Prompt sanitizado.
func (c *Context) PromptVal() string {
	return txtSanitizar(c.Request().Header.Get("HX-Prompt"))
}

// Valor del header Hx-Prompt sanitizado en mayúsculas.
func (c *Context) PromptUpper() string {
	return txtUpper(c.Request().Header.Get("HX-Prompt"))
}

// Valor del header Hx-Prompt sanitizado en minúsculas.
func (c *Context) PromptLower() string {
	return txtLower(c.Request().Header.Get("HX-Prompt"))
}

// Valor del header Hx-Prompt convertido a bool.
// Retorna false a menos de que el valor sea: "on", "true", "1".
func (c *Context) PromptBool() bool {
	return txtBool(c.Request().Header.Get("HX-Prompt"))
}

// Valor del header Hx-Prompt convertido a entero.
func (c *Context) PromptIntMust() (int, error) {
	return txtInt(c.Request().Header.Get("HX-Prompt"))
}

// Valor del header Hx-Prompt convertido a entero sin verificar error (default 0).
func (c *Context) PromptInt() int {
	num, _ := txtInt(c.Request().Header.Get("HX-Prompt"))
	return num
}

// Valor del header Hx-Prompt convertido a uint64.
func (c *Context) PromptUintMust() (uint64, error) {
	return txtUint64(c.Request().Header.Get("HX-Prompt"))
}

// Valor del header Hx-Prompt convertido a uint64 sin verificar error (default 0).
func (c *Context) PromptUint() uint64 {
	num, _ := txtUint64(c.Request().Header.Get("HX-Prompt"))
	return num
}

// Valor del header Hx-Prompt convertido a centavos.
func (c *Context) PromptCentavos() (int, error) {
	return txtCentavos(c.Request().Header.Get("HX-Prompt"))
}

// Valor del header Hx-Prompt convertido a time.
func (c *Context) PromptTime(layout string) (time.Time, error) {
	return txtTime(c.Request().Header.Get("HX-Prompt"), layout)
}

// Valor del header Hx-Prompt convertido a time, que puede estar indefinido.
func (c *Context) PromptTimeNullable(layout string) (*time.Time, error) {
	return txtTimeNullable(c.Request().Header.Get("HX-Prompt"), layout)
}

// Valor del header Hx-Prompt convertido a time desde una fecha 28/08/2022 o 2022-02-13.
func (c *Context) PromptFecha(layout string) (time.Time, error) {
	return txtFecha(c.Request().Header.Get("HX-Prompt"))
}

// Valor del path formato fecha convertido a time, que puede estar indefinido.
func (c *Context) PromptFechaNullable() (*time.Time, error) {
	return txtFechaNullable(c.Request().Header.Get("HX-Prompt"))
}
