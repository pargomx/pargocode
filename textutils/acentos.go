package textutils

import (
	"regexp"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func QuitarAcentos(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

type Utils struct {
	alfanumericos *regexp.Regexp
}

func NewTextUtils() *Utils {
	return &Utils{
		alfanumericos: regexp.MustCompile(`[^a-zA-Z0-9_]`),
	}
}

func (s *Utils) RemoveNonAlphanumeric(str string) string {
	return s.alfanumericos.ReplaceAllString(str, "")
}
