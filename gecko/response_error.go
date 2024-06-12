package gecko

import (
	"fmt"
	"net/http"
	"time"
)

// Responder con el HTTPErrorHandler definido para gecko.
func (c *Context) Error(err error) {
	c.gecko.HTTPErrorHandler(err, c)
}

func errorHandler(err error, c *Context) {
	statusCode := http.StatusInternalServerError // Default 500
	mensaje := ""
	logMsg := c.Request().URL.Path + " "
	op := NewOp(c.Request().URL.Path)

	// Si es un error http de gecko (preferido)
	if hee, ok := err.(*Gkerror); ok {
		statusCode = hee.GetCode()
		mensaje = hee.Error()
		logMsg += hee.Error()

	} else {
		mensaje = err.Error()
		logMsg += err.Error()
	}

	// Log error
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	if statusCode >= 400 {
		LogError(op.Err(err))
	}

	// Preparar respuesta
	data := map[string]any{
		"Mensaje":    mensaje,
		"StatusCode": statusCode,
	}

	if mensaje == "" {
		data["Mensaje"] = "Hubo un error al procesar tu solicitud. Intenta de manera diferente o contacta a soporte."
	}

	switch statusCode {
	case http.StatusNotFound:
		data["Titulo"] = "No encontrado"
		data["Referer"] = c.Request().Referer()
	case http.StatusInternalServerError:
		data["Titulo"] = "Error en el servidor"
	case http.StatusBadRequest:
		data["Titulo"] = "Solicitud no procesada"
	case http.StatusOK:
		data["Titulo"] = "Listo :D"
	default:
		data["Titulo"] = fmt.Sprint("Error ", statusCode)
	}

	// Renderizar el estatus
	if c.EsHTMX() {
		c.String(statusCode, mensaje)
	} else {
		c.Render(statusCode, "error", data)
	}
}

// El handler centralizado para errores de gecko.
//
// NOTA: el error se ignora cuando se genera en un middleware después de
// que un handler ya respondió algo al cliente y no hubo error con él.
func (g *Gecko) GeckoHTTPErrorHandler(err error, c *Context) {
	if err == nil {
		fmt.Println("PELIGRO: se retornó un err nil")
		return
	}
	if c == nil {
		fmt.Println("PELIGRO: context nil en ErrorHandler")
		return
	}
	if c.Response().Committed {
		fmt.Printf("PELIGRO: err generado luego de responder: %s %s %s\n", c.Request().Method, c.Request().URL.String(), err.Error())
		return
	}

	//* PREPARAR MENSAJE
	statusCode := http.StatusInternalServerError
	msgUsuario := ""
	msgLog := ""

	if errGecko, ok := err.(*Gkerror); ok {
		// Si es un error gecko (preferido)
		if errGecko == nil {
			fmt.Println("err nil response_error")
			return
		}
		statusCode = errGecko.codigo
		msgUsuario = errGecko.mensaje
		msgLog += errGecko.Error()

	} else {
		// Si es un error genérico
		msgUsuario = err.Error()
		msgLog += err.Error()
	}

	if msgUsuario == "" {
		msgUsuario = "Hubo un error, por favor intenta de otra manera o contacta a soporte."
	}

	ctx := c.Request().Method + " " + c.Request().URL.String()
	if c.Sesion.Correo != "" {
		ctx += " " + c.Sesion.Correo
	}

	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	//* REGISTRAR EN LOG
	println(time.Now().Format("2006-01-02 15:04:05") + "\033[1;31m" + " [ERROR] " +
		"\033[0;33m" + ctx + "\033[1;31m " + msgLog + "\033[0m")

	//* RESPONDER AL CLIENTE

	// Método HEAD debe responder sin body.
	if c.Request().Method == http.MethodHead {
		err := c.NoContent(statusCode)
		if err != nil {
			fmt.Println("PELIGRO: error al enviar error head: " + err.Error())
		}
		return
	}

	// HTMX solo necesita un string.
	if c.EsHTMX() {
		err = c.String(statusCode, msgUsuario)
		if err != nil {
			fmt.Println("PELIGRO: error al enviar error htmx: " + err.Error())
		}
		return
	}

	// Mandar plantilla con el error.
	// if g.ErrorTemplate != "" {
	// 	data := map[string]any{
	// 		"Mensaje":    msgUsuario,
	// 		"StatusCode": fmt.Sprint(code),
	// 		"Sesion":     c.Sesion,
	// 		"Titulo":     errorTitulo(code),
	// 	}
	// 	err = c.Render(code, e.ErrorTemplate, data)
	// 	if err != nil {
	// 		fmt.Println("PELIGRO: error al enviar error: " + err.Error())
	// 	}
	// } else {
	// c.String(code, msgUsuario)
	// }

	// Preparar respuesta
	data := map[string]any{
		"Mensaje":    msgUsuario,
		"StatusCode": statusCode,
	}
	if msgUsuario == "" {
		data["Mensaje"] = "Hubo un error al procesar tu solicitud. Intenta de manera diferente o contacta a soporte."
	}
	switch statusCode {
	case http.StatusNotFound:
		data["Titulo"] = "No encontrado"
		data["Referer"] = c.Request().Referer()
	case http.StatusInternalServerError:
		data["Titulo"] = "Error en el servidor"
	case http.StatusBadRequest:
		data["Titulo"] = "Solicitud no procesada"
	case http.StatusOK:
		data["Titulo"] = "Listo :D"
	default:
		data["Titulo"] = fmt.Sprint("Error ", statusCode)
	}

	if c.EsHTMX() {
		c.String(statusCode, msgUsuario)
	} else {
		c.Render(statusCode, "error", data)
	}
}
