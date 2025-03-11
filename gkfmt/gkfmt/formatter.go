package gkfmt

import (
	"fmt"
	"strings"

	"github.com/pargomx/gecko/gko"
)

type nodo struct {
	openTag    *openTag
	closingTag string   // ejemplo: </div>
	comment    string   // comentario html
	contenido  []string // Texto normal. inner text, or <!-- comment -->, or {{ .Something }}

	indent int // nivel de indentación
}

type openTag struct {
	tag    string   // "button", "a", "span", "script", etc.
	attr   []string // href="..." or required or {{ if .Something }}required{{ end }}
	indent int      // nivel de indentación
}

// ================================================================ //

func (n *nodo) String() string {
	indent := strings.Repeat("\t", n.indent)

	if n.openTag != nil {
		return n.openTag.String()
	}
	if n.closingTag != "" {
		return fmt.Sprintf("%s</%s>", indent, n.closingTag)
	}
	if n.comment != "" {
		return indent + n.comment
	}
	if len(n.contenido) > 0 {
		return indent + strings.Join(n.contenido, "\n"+indent)
	}

	return "ERROR"
}

func (o *openTag) String() string {
	// TODO: ordenar atributos y aplicar reglas aquí.
	indent1 := strings.Repeat("\t", o.indent)
	indent2 := strings.Repeat("\t", o.indent+1)
	return fmt.Sprintf(
		"%s<%v\n%s%v\n%v>",
		indent1,
		o.tag,
		indent2,
		strings.Join(o.attr, "\n"+indent2),
		indent2,
	)
}

// ================================================================ //

func (s *parser) tokensToElements() {
	s.indentActual = 0
	for _, token := range s.tokens {
		s.tokenToElement(token)
	}
}

func (s *parser) tokenToElement(token token) {
	switch token.tipo {

	case tipoOpenTagBeg:
		s.comenzarOpenTag(token.Txt, token.Indent)
	case tipoAtributo:
		s.openTag.attr = append(s.openTag.attr, token.Txt)
	case tipoOpenTagEnd:
		s.terminarOpenTag(s.openTag.tag)

	case tipoClosingTag:
		s.cerrarEtiqueta(token.Txt, token.Indent)

	case tipoComentario:
		s.nodos = append(s.nodos, nodo{
			comment: token.Txt,
			indent:  token.Indent,
		})

	case tipoGoTemplate:
		if s.openTag != nil {
			s.openTag.attr = append(s.openTag.attr, token.Txt)
		} else {
			s.nodos = append(s.nodos, nodo{
				contenido: strings.Split(token.Txt, "\n"),
				indent:    token.Indent,
			})
		}

	case tipoScript, tipoInnerHtml:
		s.nodos = append(s.nodos, nodo{
			contenido: []string{token.Txt},
			indent:    token.Indent,
		})

	}
}

// Declara el insideOpenTag y comienza nuevo nodo actual.
func (s *parser) comenzarOpenTag(tag string, indent int) {
	s.openTag = &openTag{
		tag:    strings.TrimPrefix(tag, "<"),
		indent: indent,
	}
}

// Termina insideOpenTag, agrega el nodo, pushes tag to stack si aplica.
func (s *parser) terminarOpenTag(tag string) {
	s.nodos = append(s.nodos, nodo{
		openTag: s.openTag,
		indent:  s.openTag.indent,
	})
	s.openTag = nil

	// No agregar self-closing tags al stack.
	switch tag {
	case "area", "base", "br", "col", "embed", "hr", "img",
		"input", "link", "meta", "param", "source", "track", "wbr":
		return
	}
	s.stackOfTags = append(s.stackOfTags, tag)
}

// Pop tag from stack. Porque fue cerrada.
func (s *parser) cerrarEtiqueta(tag string, indent int) {
	tag = strings.TrimPrefix(tag, "</")
	tag = strings.TrimSuffix(tag, ">")
	tag = strings.TrimSpace(tag)

	if len(s.stackOfTags) == 0 {
		gko.FatalExitf("bad html: extra closing tag %v", tag)
	}
	lastTag := s.stackOfTags[len(s.stackOfTags)-1]
	if lastTag != tag {
		gko.FatalExitf("bad html: found closing tag %v before %v", tag, lastTag)
	}
	s.stackOfTags = s.stackOfTags[:len(s.stackOfTags)-1]

	s.nodos = append(s.nodos, nodo{
		closingTag: tag,
		indent:     indent,
	})
}
