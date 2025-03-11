package gkfmt

import (
	"regexp"
)

// ================================================================ //
// ========== Extractor =========================================== //

type extractor struct {
	s       *parser
	luegoDe func(tipo) bool // comparar con tipo extraído anteriormente.
}

// Si regex hace match al inicio del html entonces sí comienza por ese tipo de token.
// si es true se puede usar log[1] para el index donde termina el match.
func (e *extractor) comienzaPor(loc []int) bool {
	return loc != nil && loc[0] == 0
}

// Si regex hace match pero no al principio entonces es otro tipo de token al principio.
// si es true se puede usar log[0] y log[1] para el index donde comienza y termina el match.
func (e *extractor) tieneDespués(loc []int) bool {
	return loc != nil && loc[0] > 0
}

// Retorna del html el match encontrado al principo como token del tipo especificado.
func (e *extractor) tokenEncontrado(loc []int, tipo tipo) token {
	token := token{
		Txt:  e.s.html[loc[0]:loc[1]],
		tipo: tipo}
	// Quitar token extraído y preparar para siguiente vuelta.
	e.s.html = e.s.html[loc[1]:]
	e.s.tipoTokenAnterior = token.tipo
	return token
}

// Retorna del html lo anterior al match encontrado como token del tipo especificado.
func (e *extractor) tokenAnteriorAlEncontrado(loc []int, tipo tipo) token {
	token := token{
		Txt:  e.s.html[:loc[0]],
		tipo: tipo}
	// Quitar token extraído y preparar para siguiente vuelta.
	e.s.html = e.s.html[loc[0]:]
	e.s.tipoTokenAnterior = token.tipo
	return token
}

var (
	regExtraNewLine = regexp.MustCompile(`[\s]*\n[\s]*\n[\s]*`) // line breaks intencionales.
	regWhitespace   = regexp.MustCompile(`\s+`)                 // espacio a ignorar.
	regAnyTagBeg    = regexp.MustCompile(`<`)
	regOpenTagBeg   = regexp.MustCompile(`<\w+`)
	regAtributo     = regexp.MustCompile(`([a-zA-Z\-:]+)(\s*=\s*("[^"]*"|'[^']*'|[^>\s]+))?`)
	regOpenTagEnd   = regexp.MustCompile(`>`)
	regClosingTag   = regexp.MustCompile(`<\/[a-zA-Z][a-zA-Z0-9]*\s*>`)
	regComentario   = regexp.MustCompile(`(<!--[\s\S]*?-->)`)
	regScript       = regexp.MustCompile(`<script[\s\S]*?</script>`)
	regGoTemplate   = regexp.MustCompile(`({{[\s\S]*?}})`)
)

func (s *parser) IdentificarSiguienteToken(anterior tipo) token {
	s.extractor.luegoDe = newComparador(anterior)
	e := s.extractor

	Whitespace := regWhitespace.FindStringIndex(s.html)
	ExtraNewLine := regExtraNewLine.FindStringIndex(s.html)
	AnyTagBeg := regAnyTagBeg.FindStringIndex(s.html)
	OpenTagBeg := regOpenTagBeg.FindStringIndex(s.html)
	Atributo := regAtributo.FindStringIndex(s.html)
	OpenTagEnd := regOpenTagEnd.FindStringIndex(s.html)
	ClosingTag := regClosingTag.FindStringIndex(s.html)
	Comentario := regComentario.FindStringIndex(s.html)
	Script := regScript.FindStringIndex(s.html)
	GoTemplate := regGoTemplate.FindStringIndex(s.html)

	if e.comienzaPor(ExtraNewLine) {
		s.html = s.html[Whitespace[1]:]
		return token{
			tipo:   tipoExtraNewLine,
			Txt:    "",
			Indent: 0,
		}
	}

	if e.comienzaPor(Whitespace) {
		s.html = s.html[Whitespace[1]:] // trim whitespace
		return s.IdentificarSiguienteToken(anterior)
	}

	if e.comienzaPor(Script) {
		return e.tokenEncontrado(Script, tipoScript)
	}

	if e.luegoDe(tipoTextContent) {
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
			return e.tokenAnteriorAlEncontrado(AnyTagBeg, tipoTextContent)
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
			return e.tokenAnteriorAlEncontrado(AnyTagBeg, tipoTextContent)
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
		return e.tokenAnteriorAlEncontrado(AnyTagBeg, tipoTextContent)
	}

	if e.comienzaPor(OpenTagBeg) {
		return e.tokenEncontrado(OpenTagBeg, tipoOpenTagBeg)
	}

	// If none of the above, then its just content at the end.
	return token{Txt: s.html, tipo: tipoTextContent}
}
