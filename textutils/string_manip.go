package textutils

import "strings"

// ========================================================================== //
// ========================================================================== //

// Para mejorar lectura en comparaciones repetitivas.
type CompareString string

// Verdadero si substr está dentro de CompareString.
func (s CompareString) Contiene(substr string) bool {
	return strings.Contains(string(s), substr)
}

// Verdadero si comienza por prefix
func (s CompareString) HasPrefix(prefix string) bool {
	return strings.HasPrefix(string(s), prefix)
}

// Verdadero si termina por suffix
func (s CompareString) HasSuffix(suffix string) bool {
	return strings.HasSuffix(string(s), suffix)
}

// ExtraerEnmedio retorna el string que resulta
// de quitar lo que está a partir de izq y der.
func (s CompareString) ExtraerEnmedio(izq, der string) string {
	// str="uno %% hola && dos"
	// izq="% "   der=" &"
	str := string(s)
	spli := strings.Split(str, izq)

	if len(spli) < 2 { // Cuando no hay qué cortar a la izquierda.
		return strings.Split(str, der)[0]
	}
	res := spli[1] // res="hola && dos"

	res = strings.Split(res, der)[0] // res="hola"
	return res
}

// ExtraerEnmedio retorna el string que resulta
// de quitar lo que está a partir de izq y der.
func ExtraerEnmedio(str, izq, der string) string {
	// str="uno %% hola && dos"
	// izq="% "   der=" &"

	spli := strings.Split(str, izq)

	if len(spli) < 2 { // Cuando no hay qué cortar a la izquierda.
		return strings.Split(str, der)[0]
	}
	res := spli[1] // res="hola && dos"

	res = strings.Split(res, der)[0] // res="hola"
	return res
}

// ========================================================================== //
// ========================================================================== //

// Toma un nombre de columna SQL en formato user_test_id
// y la transforma en un nombre de campo exportado UserTestID.
func SnakeToCamel(columna string) (field string) {
	palabras := strings.Split(columna, "_")
	for _, palabra := range palabras {
		if palabra == "" {
			continue
		}

		if palabra == "id" || palabra == "uuid" || palabra == "url" {
			field += strings.ToUpper(palabra)
		} else {
			field = field + strings.ToUpper(palabra[0:1]) + palabra[1:]
		}

	}
	return field
}

// ========================================================================== //
// ========================================================================== //

// Toma un string en formato kebab (cosa-de-algo)
// y lo transforma en snake lower (cosa_de_algo).
func KebabToSnake(kebab string) string {
	return strings.ToLower(strings.ReplaceAll(kebab, "-", "_"))
}

// Toma un string en formato kebab (cosa-de-algo)
// y lo transforma en snake upper (COSA_DE_ALGO).
func KebabToSnakeUp(kebab string) string {
	return strings.ToUpper(strings.ReplaceAll(kebab, "-", "_"))
}

// SnakeToKebab transforma "COSA_de_algo" en "cosa-de-algo"
func SnakeToKebab(kebab string) string {
	return strings.ToLower(strings.ReplaceAll(kebab, "_", "-"))
}

// KebabToCamel transforma "tipo-SD" en "TipoSD"
func KebabToCamel(kebab string) string {
	camel := ""
	palabras := strings.Split(kebab, "-")
	for _, palabra := range palabras {
		if palabra == "" {
			continue
		}
		if palabra == "id" || palabra == "uuid" || palabra == "url" {
			camel += strings.ToUpper(palabra)
		} else {
			camel = camel + strings.ToUpper(palabra[0:1]) + palabra[1:]
		}
	}
	return camel
}
