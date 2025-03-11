package gkfmt

import (
	"fmt"
	"strings"

	"github.com/pargomx/gecko/gko"
)

// ================================================================ //
// ========== PARSE =============================================== //

const RECURSION_LIMIT = 999999
const MAX_TOKEN_LENGHT = 2000 // caracteres

type parser struct {
	inputHTML string

	extractor         *extractor
	tokens            []token
	tipoTokenAnterior tipo

	level int // recursion level

	indentActual  int  // current indent level
	inVoidElement bool // current self-closing tag

	stackOfTags []string
	openTag     *openTag // open tag actual si la hay.
	nodos       []nodo
}

func FormatGeckoTemplate(html string, builder *strings.Builder, fmtConTokens bool) []token {
	s := parser{
		inputHTML:         html,
		extractor:         &extractor{},
		tokens:            []token{},
		tipoTokenAnterior: tipoInnerHtml, // Estado inicial para extractor.
	}
	s.parseRecursive()
	s.tokensToElements()

	// Ambas formas de dar formato
	if fmtConTokens {
		for _, token := range s.tokens {
			builder.WriteString(fmt.Sprintf("%v%v\n", strings.Repeat("\t", token.Indent), token.Txt))
		}

	} else {
		for _, nodo := range s.nodos {
			builder.WriteString(nodo.String() + "\n")
		}
	}

	return s.tokens
}

func (s *parser) parseRecursive() {
	// Final inesperado de la recursión.
	s.level++
	if s.level > RECURSION_LIMIT {
		gko.FatalExit("Recursión descontrolada :0")
	}

	// Final normal de la recursión.
	s.inputHTML = strings.TrimSpace(s.inputHTML)
	if s.inputHTML == "" {
		return
	}

	// Extraer token.
	token := s.IdentificarSiguienteToken(s.inputHTML, s.tipoTokenAnterior)
	s.tipoTokenAnterior = token.tipo
	s.inputHTML = strings.TrimPrefix(s.inputHTML, token.Txt)

	// Error si es algo que probablemente no cachó el extractor.
	if len(token.Txt) > MAX_TOKEN_LENGHT && token.tipo == tipoInnerHtml {
		gko.FatalExit(token.Txt[:MAX_TOKEN_LENGHT-1])
	}

	// Indentación.
	es := newComparador(token.tipo)
	if token.tipo == tipoOpenTagBeg {
		token.Indent = s.indentActual
		s.indentActual++
		if token.esSelfClosingTag() {
			s.inVoidElement = true
		}

	} else if s.inVoidElement && (es(tipoClosingTag) || es(tipoOpenTagEnd)) {
		token.Indent = s.indentActual
		s.indentActual--
		s.inVoidElement = false

	} else if token.tipo == tipoClosingTag {
		s.indentActual--
		if s.inVoidElement {
			s.indentActual--
			s.inVoidElement = false
		}
		token.Indent = s.indentActual

	} else {
		token.Indent = s.indentActual
	}

	// gko.LogDebugf("%03d %v: %v'%v'", f.level, token.Tipo(), strings.Repeat("  ", token.Indent), token.Txt)

	// Guardar token
	s.tokens = append(s.tokens, token)

	// Continuar recursión.
	s.parseRecursive()
}
