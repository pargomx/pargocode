package gkfmt

import (
	"regexp"
	"strings"
)

// ================================================================ //
// ========== Extractor =========================================== //

type extractor struct {
	html    string          // el token se extrae del inicio del html.
	luegoDe func(tipo) bool // comparar con tipo extraído anteriormente.
}

// Si regex hace match al inicio del html entonces sí comienza por ese tipo de token.
func (e *extractor) comienzaPor(loc []int) bool {
	return loc != nil && loc[0] == 0
}

// Si regex hace match pero no al principio entonces es otro tipo de token al principio.
func (e *extractor) tieneDespués(loc []int) bool {
	return loc != nil && loc[0] > 0
}

// Retorna del html el match encontrado al principo como token del tipo especificado.
func (e *extractor) tokenEncontrado(loc []int, tipo tipo) token {
	return token{
		Txt:  strings.TrimSpace(e.html[loc[0]:loc[1]]),
		tipo: tipo}
}

// Retorna del html lo anterior al match encontrado como token del tipo especificado.
func (e *extractor) tokenAnteriorAlEncontrado(loc []int, tipo tipo) token {
	return token{
		Txt:  strings.TrimSpace(e.html[:loc[0]]),
		tipo: tipo}
}

var (
	regAnyTagBeg  *regexp.Regexp = regexp.MustCompile(`<`)
	regOpenTagBeg *regexp.Regexp = regexp.MustCompile(`<\w+`)
	regAtributo   *regexp.Regexp = regexp.MustCompile(`([a-zA-Z\-:]+)(\s*=\s*("[^"]*"|'[^']*'|[^>\s]+))?`)
	regOpenTagEnd *regexp.Regexp = regexp.MustCompile(`>`)
	regClosingTag *regexp.Regexp = regexp.MustCompile(`<\/[a-zA-Z][a-zA-Z0-9]*\s*>`)
	regComentario *regexp.Regexp = regexp.MustCompile(`(<!--[\s\S]*?-->)`)
	regScript     *regexp.Regexp = regexp.MustCompile(`<script[\s\S]*?</script>`)
	regGoTemplate *regexp.Regexp = regexp.MustCompile(`({{[\s\S]*?}})`)
)

func (s *parser) IdentificarSiguienteToken(html string, anterior tipo) token {
	e := s.extractor
	e.html = html
	e.luegoDe = newComparador(anterior)

	AnyTagBeg := regAnyTagBeg.FindStringIndex(e.html)
	OpenTagBeg := regOpenTagBeg.FindStringIndex(e.html)
	Atributo := regAtributo.FindStringIndex(e.html)
	OpenTagEnd := regOpenTagEnd.FindStringIndex(e.html)
	ClosingTag := regClosingTag.FindStringIndex(e.html)
	Comentario := regComentario.FindStringIndex(e.html)
	Script := regScript.FindStringIndex(e.html)
	GoTemplate := regGoTemplate.FindStringIndex(e.html)

	if e.comienzaPor(Script) {
		return e.tokenEncontrado(Script, tipoScript)
	}

	if e.luegoDe(tipoInnerHtml) {
		if e.comienzaPor(OpenTagBeg) {
			return e.tokenEncontrado(OpenTagBeg, tipoOpenTagBeg)
		}
	}

	if e.luegoDe(tipoOpenTagBeg) || e.luegoDe(tipoAtributo) {
		if e.comienzaPor(Atributo) {
			return e.tokenEncontrado(Atributo, tipoAtributo)
		}
		if e.comienzaPor(OpenTagEnd) {
			return e.tokenEncontrado(OpenTagEnd, tipoOpenTagEnd)
		}
		if e.comienzaPor(GoTemplate) {
			return e.tokenEncontrado(GoTemplate, tipoGoTemplate)
		}
	}

	if e.luegoDe(tipoGoTemplate) {
		if e.comienzaPor(OpenTagBeg) {
			return e.tokenEncontrado(OpenTagBeg, tipoOpenTagBeg)
		}
		if e.comienzaPor(ClosingTag) {
			return e.tokenEncontrado(ClosingTag, tipoClosingTag)
		}
		if e.comienzaPor(Atributo) {
			return e.tokenEncontrado(Atributo, tipoAtributo)
		}
		if e.comienzaPor(OpenTagEnd) {
			return e.tokenEncontrado(OpenTagEnd, tipoOpenTagEnd)
		}
		if e.comienzaPor(GoTemplate) {
			return e.tokenEncontrado(GoTemplate, tipoGoTemplate)
		}
	}

	if e.luegoDe(tipoComentario) {
		if e.comienzaPor(OpenTagBeg) {
			return e.tokenEncontrado(OpenTagBeg, tipoOpenTagBeg)
		}
		if e.comienzaPor(ClosingTag) {
			return e.tokenEncontrado(ClosingTag, tipoClosingTag)
		}
		if e.comienzaPor(GoTemplate) {
			return e.tokenEncontrado(GoTemplate, tipoGoTemplate)
		}
	}

	if e.luegoDe(tipoOpenTagEnd) {
		if e.comienzaPor(OpenTagBeg) {
			return e.tokenEncontrado(OpenTagBeg, tipoOpenTagBeg)
		}
		if e.comienzaPor(ClosingTag) {
			return e.tokenEncontrado(ClosingTag, tipoClosingTag)
		}
		if e.tieneDespués(AnyTagBeg) {
			return e.tokenAnteriorAlEncontrado(AnyTagBeg, tipoInnerHtml)
		}
	}

	if e.luegoDe(tipoClosingTag) {
		if e.comienzaPor(OpenTagBeg) {
			return e.tokenEncontrado(OpenTagBeg, tipoOpenTagBeg)
		}
		if e.comienzaPor(ClosingTag) {
			return e.tokenEncontrado(ClosingTag, tipoClosingTag)
		}
		if e.tieneDespués(AnyTagBeg) {
			return e.tokenAnteriorAlEncontrado(AnyTagBeg, tipoInnerHtml)
		}
	}

	if e.comienzaPor(ClosingTag) {
		return e.tokenEncontrado(ClosingTag, tipoClosingTag)
	}

	if e.comienzaPor(Comentario) {
		return e.tokenEncontrado(Comentario, tipoComentario)
	}

	if e.comienzaPor(GoTemplate) {
		return e.tokenEncontrado(GoTemplate, tipoGoTemplate)
	}

	if e.tieneDespués(AnyTagBeg) {
		return e.tokenAnteriorAlEncontrado(AnyTagBeg, tipoInnerHtml)
	}

	if e.comienzaPor(OpenTagBeg) {
		return e.tokenEncontrado(OpenTagBeg, tipoOpenTagBeg)
	}

	// If none of the above, then its just content at the end.
	return token{Txt: e.html, tipo: tipoInnerHtml}
}
