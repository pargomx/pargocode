package gecko

import (
	"fmt"
	"time"
)

// Indica si se imprimirá la timestamp al inicio
// de cada entrada de log al Stdout.
var PrintLogTimestamps bool = true

// 2020-11-25 18:54:32 [DEBUG] Algo interesante sucede cyan.
func (c *Context) LogDebug(a ...any) {
	println(timestamp() + cCyan + "[DEBUG] " + rWhite + fmt.Sprint(a...) + reset)
}

// 2020-11-25 18:54:32 [INFOR] Algo interesante sucede cyan.
func (c *Context) LogInfo(a ...any) {
	println(timestamp() + cCyan + "[INFOR] " + rWhite + fmt.Sprint(a...) + reset)
}

// 2020-11-25 18:54:32 [EVENT] Algo importante sucede cyan bold.
func (c *Context) LogEvento(a ...any) {
	println(timestamp() + cCyan + "[EVENT] " + bold + fmt.Sprint(a...) + reset)

}

// Warn LOG AVISO al fallar en un proceso no escencial
func (c *Context) LogWarn(a ...any) {
	println(timestamp() + cYellow + "[AVISO] " + reset + fmt.Sprint(a...))
}

// Abort LOG ABORT al fallar en un proceso no escencial por error
func (c *Context) LogAbort(a ...any) {
	println(timestamp() + cYellow + "[ABORT] " + reset + fmt.Sprint(a...))
}

// LOG ERROR
func (c *Context) LogError(a ...any) {
	if len(a) == 1 {
		if e, ok := a[0].(*Gkerror); ok {
			LogError(e)
		}
	}
	println(timestamp() + cRed + "[ERROR] " + reset + fmt.Sprint(a...))
}

// Okey LOG LISTO al terminar con éxito una función
func (c *Context) LogOkey(a ...any) {
	println(timestamp() + cGreen + "[LISTO] " + reset + fmt.Sprint(a...))
}

// ================================================================ //
// ========== FMT ================================================= //

// 2020-11-25 18:54:32 [INFOR] Algo interesante sucede cyan.
func (c *Context) LogInfof(format string, a ...any) {
	println(timestamp() + cCyan + "[INFOR] " + rWhite + fmt.Sprintf(format, a...) + reset)
}

// 2020-11-25 18:54:32 [EVENT] Algo importante sucede cyan bold.
func (c *Context) LogEventof(format string, a ...any) {
	println(timestamp() + cCyan + "[EVENT] " + bold + fmt.Sprintf(format, a...) + reset)
}

// Warn LOG AVISO al fallar en un proceso no escencial
func (c *Context) LogWarnf(format string, a ...any) {
	println(timestamp() + cYellow + "[AVISO] " + reset + fmt.Sprintf(format, a...))
}

// Abort LOG ABORT al fallar en un proceso no escencial por error
func (c *Context) LogAbortf(format string, a ...any) {
	println(timestamp() + cYellow + "[ABORT] " + reset + fmt.Sprintf(format, a...))
}

// Okey LOG LISTO al terminar con éxito una función
func (c *Context) LogOkeyf(format string, a ...any) {
	println(timestamp() + cGreen + "[LISTO] " + reset + fmt.Sprintf(format, a...))
}

// ================================================================ //
// ========== STANDALONE ========================================== //

func LogDebug(op *Gkerror, format string, a ...any) {
	println(timestamp() + cCyan + "[DEBUG] " + reset + fmt.Sprintf(op.Error()+" "+format, a...))
}
func LogAlert(op *Gkerror, format string, a ...any) {
	println(timestamp() + cYellow + "[ALERT] " + reset + fmt.Sprintf(op.Error()+" "+format, a...))
}

func LogInfof(format string, a ...any) {
	println(timestamp() + cCyan + "[INFOR] " + rWhite + fmt.Sprintf(format, a...) + reset)
}
func LogEventof(format string, a ...any) {
	println(timestamp() + cCyan + "[EVENT] " + bold + fmt.Sprintf(format, a...) + reset)
}
func LogWarnf(format string, a ...any) {
	println(timestamp() + cYellow + "[AVISO] " + reset + fmt.Sprintf(format, a...))
}
func LogAbortf(format string, a ...any) {
	println(timestamp() + cYellow + "[ABORT] " + reset + fmt.Sprintf(format, a...))
}
func LogOkeyf(format string, a ...any) {
	println(timestamp() + cGreen + "[LISTO] " + reset + fmt.Sprintf(format, a...))
}

// ================================================================ //
// ========== TIME ================================================ //
func timestamp() string {
	if PrintLogTimestamps {
		return time.Now().Format("2006-01-02 15:04:05") + " "
	}
	return ""
}

// ================================================================ //
// ========== ERROR =============================================== //

// Error imprime el error en la consola con formato.
func LogError(err error) {
	if err == nil {
		println(timestamp() + bRed + "[ERROR] " + "\033[33mError nulo.\033[0m")
		return
	}
	e, ok := err.(*Gkerror)
	if !ok { // Cualquier otro tipo de error
		println(timestamp() + bRed + "[ERROR] " + "\033[31m" + err.Error() + "\033[0m")
		return
	}

	msg := timestamp()
	if e.codigo > 0 {
		msg += bRed + fmt.Sprintf("[ERROR] (%d)", e.codigo) + reset
	} else {
		msg += bRed + "[ERROR]" + reset
	}
	if e.operación != "" {
		msg += " " + rYellow + e.operación
	}
	if e.err != nil {
		msg += " " + rRed + e.err.Error()
		if e.mensaje != "" && MostrarMensajeEnErrores {
			msg += ":"
		}
	}
	if e.mensaje != "" && MostrarMensajeEnErrores {
		msg += " " + bRed + e.mensaje + "." + reset
	}
	if e.contexto != "" {
		msg += " " + rPurple + e.contexto
	}
	println(msg + reset)
}

// ================================================================ //
// ========== COLORES ============================================= //

const (
	reset = "\033[0m"

	bold = "\033[1m"
	dim  = "\033[2m"

	// Color
	cRed    = "\033[31m"
	cGreen  = "\033[32m"
	cYellow = "\033[33m"
	cBlue   = "\033[34m"
	cPurple = "\033[35m"
	cCyan   = "\033[36m"
	cGray   = "\033[37m"
	cWhite  = "\033[97m"

	// Reset and then color
	rRed    = "\033[0;31m"
	rGreen  = "\033[0;32m"
	rYellow = "\033[0;33m"
	rBlue   = "\033[0;34m"
	rPurple = "\033[0;35m"
	rCyan   = "\033[0;36m"
	rGray   = "\033[0;37m"
	rWhite  = "\033[0;97m"

	// Bold colors
	bRed    = "\033[1;31m"
	bGreen  = "\033[1;32m"
	bYellow = "\033[1;33m"
	bBlue   = "\033[1;34m"
	bPurple = "\033[1;35m"
	bCyan   = "\033[1;36m"
	bGray   = "\033[1;37m"
	bWhite  = "\033[1;97m"
)
