package gecko

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ================================================================ //
// ========== PROCESAMIENTO ======================================= //

// txtSanitizar corta los espacios extra y hace
// las siguientes sustituciones de caracteres:
//   - `<` `>` por `-`
//   - `;` por `,`
//   - `"` por `'`
//
// Se espera que ayude a mitigar algún tipo de ataque básico.
func txtSanitizar(txt string) string {
	txt = txtQuitarEspacios(txt)
	txt = strings.ReplaceAll(txt, ">", "-")
	txt = strings.ReplaceAll(txt, "<", "-")
	txt = strings.ReplaceAll(txt, ";", ",")
	txt = strings.ReplaceAll(txt, "\"", "'")
	return txt
}

// Remueve todos los espacios al inicio, al final y dobles.
// También sustituye los saltos de línea y tabuladores por
// espacios simples.
func txtQuitarEspacios(txt string) string {
	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(txt, " "))
}

// Retorna el valor sanitizado en mayúsculas.
func txtUpper(txt string) string {
	return strings.ToUpper(txtSanitizar(txt))
}

// Retorna el valor sanitizado en minúsculas.
func txtLower(txt string) string {
	return strings.ToLower(txtSanitizar(txt))
}

// Retorna false a menos de que el valor sea:
// "on", "true", "1".
func txtBool(txt string) bool {
	str := txtLower(txt)
	return str == "on" || str == "true" || str == "1"
}

// Retorna el valor en tipo entero.
// Retorna error si no es un número válido.
func txtInt(txt string) (int, error) {
	return strconv.Atoi(txtSanitizar(txt))
}

// Retorna el valor en tipo entero positivo de 8 bytes.
// Valor máximo aceptado: 18446744073709551615.
func txtUint64(txt string) (uint64, error) {
	return strconv.ParseUint(txt, 10, 64)
}

// Devuelve los centavos a partir de un string de dinero
// que puede ser "$200.00", "200", "200.0" por ejemplo.
//
// El valor recibido debe estar en unidades de pesos.
// Se puede incluir centavos pero deben estar como decimales.
func txtCentavos(txt string) (int, error) {
	str := txtSanitizar(txt)
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, ",", "")
	str = strings.TrimLeft(str, "$")

	partes := strings.Split(str, ".")
	if len(partes) > 2 {
		return 0, errors.New("más de 1 punto para centavos: " + str)
	}
	pesos := partes[0]
	centavos := ""
	if len(partes) == 2 {
		centavos = partes[1]
	}

	switch len(centavos) {
	case 0:
		str = pesos + "00"
	case 1:
		str = pesos + centavos + "0"
	case 2:
		str = pesos + centavos
	default:
		return 0, errors.New("no se admite más que centavos luego del punto: [" + centavos + "]")
	}

	res, err := strconv.Atoi(str)
	if err != nil {
		return 0, errors.New("no es un número válido: [" + str + "]")
	}

	return res, nil
}

// ====== Fechas y tiempo ======= //

// Convierte el valor del formulario con la clave "key" en
// tiempo utilizando el layout dado en la zona horaria de
// México central.
//
// Si no hay un valor para la fecha, se considera error.
//
// Utiliza time.ParseInLocation porque sino MySQL recibe la
// fecha en UTC cuando la espera en hora local, y esto hace que
// la fecha guardada sea un día anterior (6h) al esperado.
//
// Por ejemplo: key="fecha_inicio" layout="2006-01-02".
func txtTime(txt string, layout string) (time.Time, error) {
	txt = txtSanitizar(txt)
	if txt == "" {
		return time.Time{}, errors.New("valor indefinido")
	}
	fechaTxt := strings.ReplaceAll(txt, "/", "-") // Permitir "2006/01/02"
	// ParseInLocation porque sino MySQL recibe la fecha en UTC cuando la espera en hora local.
	tz, err := time.LoadLocation("America/Mexico_City")
	if err != nil {
		return time.Time{}, err
	}
	tiempo, err := time.ParseInLocation(layout, fechaTxt, tz)
	if err != nil {
		return time.Time{}, err
	}
	return tiempo, nil
}

// Convierte el valor del formulario con la clave "key" en
// tiempo utilizando el layout dado en la zona horaria de
// México central.
//
// Si no hay valor para la fecha, se devuelve nil sin error.
//
// Utiliza time.ParseInLocation porque sino MySQL recibe la
// fecha en UTC cuando la espera en hora local, y esto hace que
// la fecha guardada sea un día anterior (6h) al esperado.
//
// Por ejemplo: key="fecha_inicio" layout="2006-01-02".
func txtTimeNullable(txt, layout string) (*time.Time, error) {
	txt = txtSanitizar(txt)
	if txt == "" {
		return nil, nil
	}
	fechaTxt := strings.ReplaceAll(txt, "/", "-") // Permitir "2006/01/02"
	// ParseInLocation porque sino MySQL recibe la fecha en UTC cuando la espera en hora local.
	tz, err := time.LoadLocation("America/Mexico_City")
	if err != nil {
		return nil, err
	}
	tiempo, err := time.ParseInLocation(layout, fechaTxt, tz)
	if err != nil {
		return nil, err
	}
	return &tiempo, nil
}

// Convierte el valor del formulario con la clave "key" en
// tiempo utilizando el layout dado en la zona horaria de
// México central.
//
// Si no hay un valor para la fecha, se considera error.
//
// Utiliza time.ParseInLocation porque sino MySQL recibe la
// fecha en UTC cuando la espera en hora local, y esto hace que
// la fecha guardada sea un día anterior (6h) al esperado.
//
// Por ejemplo: key="fecha_inicio" layout="2006-01-02".
// Acepta fechas así: "2006-01-02" "2006/01/02" "28-01-2006" "28/01/2006".
func txtFecha(txt string) (time.Time, error) {
	txt = txtSanitizar(txt)
	txt = strings.ReplaceAll(txt, " ", "") // Quitar espacios
	if txt == "" {
		return time.Time{}, errors.New("valor indefinido")
	}
	if len(txt) != 10 {
		return time.Time{}, errors.New("debe tener el año completo y separarse por guiones. Por ejemplo: 30/01/2023")
	}
	txt = strings.ReplaceAll(txt, "/", "-") // Permitir "2006/01/02"
	if txt[2:3] == "-" {                    // Permitir fecha volteada: "28-01-2006"
		txt = txt[6:] + "-" + txt[3:5] + "-" + txt[:2]
	}
	// ParseInLocation porque sino MySQL recibe la fecha en UTC cuando la espera en hora local.
	tz, err := time.LoadLocation("America/Mexico_City")
	if err != nil {
		return time.Time{}, err
	}
	tiempo, err := time.ParseInLocation("2006-01-02", txt, tz)
	if err != nil {
		return time.Time{}, err
	}
	return tiempo, nil
}

// Convierte el valor del formulario con la clave "key" en
// fecha utilizando el layout dado en la zona horaria de
// México central.
//
// Si no hay valor para la fecha, se devuelve nil sin error.
//
// Utiliza time.ParseInLocation porque sino MySQL recibe la
// fecha en UTC cuando la espera en hora local, y esto hace que
// la fecha guardada sea un día anterior (6h) al esperado.
//
// Por ejemplo: key="fecha_inicio" layout="2006-01-02".
//
// Acepta fechas así: "2006-01-02" "2006/01/02" "28-01-2006" "28/01/2006".
func txtFechaNullable(txt string) (*time.Time, error) {
	txt = txtSanitizar(txt)
	txt = strings.ReplaceAll(txt, " ", "") // Quitar espacios
	if txt == "" {
		return nil, nil
	}
	if len(txt) != 10 {
		return nil, errors.New("debe tener el año completo y separarse por guiones. Por ejemplo: 30/01/2023")
	}
	txt = strings.ReplaceAll(txt, "/", "-") // Permitir "2006/01/02"
	if txt[2:3] == "-" {                    // Permitir fecha volteada: "28-01-2006"
		txt = txt[6:] + "-" + txt[3:5] + "-" + txt[:2]
	}

	// ParseInLocation porque sino MySQL recibe la fecha en UTC cuando la espera en hora local.
	tz, err := time.LoadLocation("America/Mexico_City")
	if err != nil {
		return nil, err
	}
	tiempo, err := time.ParseInLocation("2006-01-02", txt, tz)
	if err != nil {
		return nil, err
	}
	return &tiempo, nil
}
