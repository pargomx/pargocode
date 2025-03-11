package gkfmt

import (
	"regexp"
	"strings"

	"github.com/pargomx/gecko/gko"
)

type token struct {
	Txt    string // puede ser '<tag', 'attr="x"'.
	tipo   tipo
	Indent int // nivel de indentación.
}

type tipo int

const (
	tipoInnerHtml  tipo = iota // Contenido que puede ser texto y ya no contiene nignún "<"
	tipoOpenTagBeg             // <button
	tipoAtributo               // hx-get="/x", required, class="{{ if not .Show }}hidden {{ end }}"
	tipoOpenTagEnd             // >
	tipoClosingTag             // </button>
	tipoComentario             // <!-- Comentario -->
	tipoScript                 // <script> ... </script>
	tipoGoTemplate             // {{ .Something }}
)

func (t token) Tipo() string {
	return [...]string{
		"_text_",
		"tagIni",
		"atribu",
		"tagEnd",
		"tagClo",
		"coment",
		"script",
		"gotmpl",
	}[t.tipo]
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

type Formatter struct {
	inputHTML string
	tokens    []token
	level     int // recursion level

	indentActual  int  // current indent level
	inVoidElement bool // current self-closing tag
}

func Extract(html string) []token {
	f := Formatter{
		inputHTML: html,
		tokens:    []token{},
	}
	f.ExtractRecursive()
	return f.tokens
}

func (f *Formatter) ExtractRecursive() {
	f.level++
	if f.level > 9999999 {
		gko.FatalExit("Runaway recursion")
	}

	f.inputHTML = strings.TrimSpace(f.inputHTML)
	if f.inputHTML == "" {
		return
	}

	lastTokenTipo := tipoInnerHtml
	if len(f.tokens) > 0 {
		lastTokenTipo = f.tokens[len(f.tokens)-1].tipo
	}

	token := ExtractPrefixToken(f.inputHTML, lastTokenTipo)
	token.Txt = strings.TrimSpace(token.Txt)

	if len(token.Txt) > 1001 && token.tipo == tipoInnerHtml {
		gko.FatalExit(token.Txt[:1000])
	}

	es := compararTipo(token.tipo)

	// Indentación.
	if token.tipo == tipoOpenTagBeg {
		token.Indent = f.indentActual
		f.indentActual++
		if token.EsSelfClosingTag() {
			f.inVoidElement = true
		}

	} else if f.inVoidElement && (es(tipoClosingTag) || es(tipoOpenTagEnd)) {
		token.Indent = f.indentActual
		f.indentActual--
		f.inVoidElement = false

	} else if token.tipo == tipoClosingTag {
		f.indentActual--
		if f.inVoidElement {
			f.indentActual--
			f.inVoidElement = false
		}
		token.Indent = f.indentActual

	} else {
		token.Indent = f.indentActual
	}

	// gko.LogDebugf("%03d %v: %v'%v'", f.level, token.Tipo(), strings.Repeat("  ", token.Indent), token.Txt)

	// Guardar token
	f.tokens = append(f.tokens, token)
	f.inputHTML = strings.TrimPrefix(f.inputHTML, token.Txt)

	f.ExtractRecursive()
}

func comienzaPor(loc []int) bool {
	return loc != nil && loc[0] == 0
}
func tieneDespués(loc []int) bool {
	return loc != nil && loc[0] > 0
}

// Extrae del html el match encontrado.
func extraerFound(html string, loc []int, tipo tipo) token {
	return token{Txt: html[loc[0]:loc[1]], tipo: tipo}
}

// Extrae del html lo anterior al match encontrado.
func extraerBeforeFound(html string, loc []int, tipo tipo) token {
	return token{Txt: html[:loc[0]], tipo: tipo}
}

func compararTipo(base tipo) func(tipo) bool {
	return func(t tipo) bool {
		return base == t
	}
}

func ExtractPrefixToken(html string, anterior tipo) token {
	luegoDe := compararTipo(anterior)
	OpenTagBeg := regOpenTagBeg.FindStringIndex(html)
	Atributo := regAtributo.FindStringIndex(html)
	OpenTagEnd := regOpenTagEnd.FindStringIndex(html)
	ClosingTag := regClosingTag.FindStringIndex(html)
	AnyTagBeg := regAnyTagBeg.FindStringIndex(html)
	Comentario := regComentario.FindStringIndex(html)
	GoTemplate := regGoTemplate.FindStringIndex(html)
	InlineScript := regScript.FindStringIndex(html)

	if comienzaPor(InlineScript) {
		return extraerFound(html, InlineScript, tipoScript)
	}

	if luegoDe(tipoInnerHtml) {
		if comienzaPor(OpenTagBeg) {
			return extraerFound(html, OpenTagBeg, tipoOpenTagBeg)
		}
	}

	if luegoDe(tipoOpenTagBeg) || luegoDe(tipoAtributo) {
		if comienzaPor(Atributo) {
			return extraerFound(html, Atributo, tipoAtributo)
		}
		if comienzaPor(OpenTagEnd) {
			return extraerFound(html, OpenTagEnd, tipoOpenTagEnd)
		}
		if comienzaPor(GoTemplate) {
			return extraerFound(html, GoTemplate, tipoGoTemplate)
		}
	}

	if luegoDe(tipoGoTemplate) {
		if comienzaPor(OpenTagBeg) {
			return extraerFound(html, OpenTagBeg, tipoOpenTagBeg)
		}
		if comienzaPor(ClosingTag) {
			return extraerFound(html, ClosingTag, tipoClosingTag)
		}
		if comienzaPor(Atributo) {
			return extraerFound(html, Atributo, tipoAtributo)
		}
		if comienzaPor(OpenTagEnd) {
			return extraerFound(html, OpenTagEnd, tipoOpenTagEnd)
		}
		if comienzaPor(GoTemplate) {
			return extraerFound(html, GoTemplate, tipoGoTemplate)
		}
	}

	if luegoDe(tipoComentario) {
		if comienzaPor(OpenTagBeg) {
			return extraerFound(html, OpenTagBeg, tipoOpenTagBeg)
		}
		if comienzaPor(ClosingTag) {
			return extraerFound(html, ClosingTag, tipoClosingTag)
		}
		if comienzaPor(GoTemplate) {
			return extraerFound(html, GoTemplate, tipoGoTemplate)
		}
	}

	if luegoDe(tipoOpenTagEnd) {
		if comienzaPor(OpenTagBeg) {
			return extraerFound(html, OpenTagBeg, tipoOpenTagBeg)
		}
		if comienzaPor(ClosingTag) {
			return extraerFound(html, ClosingTag, tipoClosingTag)
		}
		if tieneDespués(AnyTagBeg) {
			return extraerBeforeFound(html, AnyTagBeg, tipoInnerHtml)
		}
	}

	if luegoDe(tipoClosingTag) {
		if comienzaPor(OpenTagBeg) {
			return extraerFound(html, OpenTagBeg, tipoOpenTagBeg)
		}
		if comienzaPor(ClosingTag) {
			return extraerFound(html, ClosingTag, tipoClosingTag)
		}
		if tieneDespués(AnyTagBeg) {
			return extraerBeforeFound(html, AnyTagBeg, tipoInnerHtml)
		}
	}

	if comienzaPor(ClosingTag) {
		return extraerFound(html, ClosingTag, tipoClosingTag)
	}

	if comienzaPor(Comentario) {
		return extraerFound(html, Comentario, tipoComentario)
	}

	if comienzaPor(GoTemplate) {
		return extraerFound(html, GoTemplate, tipoGoTemplate)
	}

	if tieneDespués(AnyTagBeg) {
		return extraerBeforeFound(html, AnyTagBeg, tipoInnerHtml)
	}

	if comienzaPor(OpenTagBeg) {
		return extraerFound(html, OpenTagBeg, tipoOpenTagBeg)
	}

	// If none of the above, then its just content at the end.
	return token{Txt: html, tipo: tipoInnerHtml}
}
