package textutils

import (
	"unicode"
)

// ========================================================================== //
// ========== Mayúsculas ==================================================== //

func PrimeraMayusc(str string) string {
	if str == "" {
		return ""
	}
	r := []rune(str)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func QuitarPrimeraMayusc(str string) string {
	if str == "" {
		return ""
	}
	r := []rune(str)
	// Omitir si la segunda letra es mayúscula.

	if len(r) > 2 &&
		unicode.ToUpper(r[1]) == r[1] {
		return str
	}
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
