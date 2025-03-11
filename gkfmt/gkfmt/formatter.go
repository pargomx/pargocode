package gkfmt

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/pargomx/gecko/gko"
)

type nodo struct {
	openTag    *openTag
	closingTag string   // ejemplo: </div>
	comment    string   // comentario html
	contenido  []string // Texto normal. inner text, or <!-- comment -->, or {{ .Something }}
	extraSpace bool     // Si solamente es un salto de línea adicional

	indent int // nivel de indentación
}

type openTag struct {
	tag    string   // "button", "a", "span", "script", etc.
	attr   []string // href="..." or required or {{ if .Something }}required{{ end }}
	indent int      // nivel de indentación

	attrCondicionales []string
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
	if n.extraSpace {
		return ""
	}

	return "ERROR"
}

// Ordena los atributos y aplica reglas de poner en una o varias líneas.
func (o *openTag) String() string {

	sort.SliceStable(o.attr, func(i, j int) bool {
		pi, pj := getAttrPriority(o.attr[i]), getAttrPriority(o.attr[j])
		if pi == pj { // si son de la misma prioridad entonces alfabéticamente.
			return o.attr[i] < o.attr[j]
		}
		return pi < pj
	})

	indent1 := strings.Repeat("\t", o.indent)
	indent2 := strings.Repeat("\t", o.indent+1)

	primeraLinea := o.tag

	// Si tiene un atributo cualquiera, poner inline.
	if len(o.attr) == 1 {
		return fmt.Sprintf("%s<%s %s>", indent1, primeraLinea, o.attr[0])
	}

	// El id, tipo y class siempre van inline en ese orden.
	if atrib := o.sacarAtributoEq("id"); atrib != "" {
		primeraLinea += " " + atrib
	}
	if atrib := o.sacarAtributoEq("tipo"); atrib != "" {
		primeraLinea += " " + atrib
	}
	if atrib := o.sacarAtributoEq("class"); atrib != "" {
		primeraLinea += " " + atrib
	}

	// Si no hay ningún atributo más allá de id, tipo, class dejar inline.
	if len(o.attr) == 0 {
		return fmt.Sprintf("%s<%s>", indent1, primeraLinea)
	}

	return fmt.Sprintf(
		"%s<%s\n%s%s\n%s>",
		indent1,
		primeraLinea,
		indent2,
		strings.Join(o.attr, "\n"+indent2),
		indent2,
	)
}

// Para el orden de los atributos.
func getAttrPriority(attr string) int {
	switch {
	case strings.HasPrefix(attr, "id="):
		return 1
	case strings.HasPrefix(attr, "tipo="):
		return 2
	case strings.HasPrefix(attr, "class="):
		return 3
	case strings.HasPrefix(attr, "hx-"):
		return 4
	case strings.HasPrefix(attr, "type="):
		return 5
	default:
		return 6 // será alfabético
	}
}

// Si el atributo se encuentra antes de un "=" se devuelve
// y se elmina de la lista de atributos.
func (o *openTag) sacarAtributoEq(attrName string) string {
	for i, atrib := range o.attr {
		if strings.HasPrefix(atrib, attrName+"=") {
			o.attr = slices.Delete(o.attr, i, i+1)
			return atrib
		}
	}
	return ""
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
		if len(s.openTag.attrCondicionales) > 0 {
			s.openTag.attrCondicionales = append(s.openTag.attrCondicionales, token.Txt)
		} else {
			s.openTag.attr = append(s.openTag.attr, token.Txt)
		}

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
		s.agregarGoSintax(token)

	case tipoExtraNewLine:
		s.nodos = append(s.nodos, nodo{
			extraSpace: true,
		})

	case tipoScript:
		s.nodos = append(s.nodos, nodo{
			contenido: []string{token.Txt},
			indent:    token.Indent,
		})

	case tipoTextContent:
		// Permitir line breaks extra dentro del contenido.
		split := strings.Split(token.Txt, "\n")
		result := []string{}

		afterEmptyLine := false

		for _, line := range split {
			line := strings.TrimSpace(line)

			// Ignorar empty lines si solo es una
			if line == "" && !afterEmptyLine {
				afterEmptyLine = true
				continue

			} else { // podría ser emptyline afterEmptyLine
				afterEmptyLine = false
				result = append(result, line)
			}
		}
		s.nodos = append(s.nodos, nodo{
			contenido: result,
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

func (s *parser) agregarGoSintax(token token) {
	// Puede ser simple contenido.
	if s.openTag == nil {
		s.nodos = append(s.nodos, nodo{
			contenido: strings.Split(token.Txt, "\n"),
			indent:    token.Indent,
		})
		return
	}

	// O ser parte de un atributo condicional.
	if len(s.openTag.attrCondicionales) == 0 || s.openTag.attrCondicionales == nil {
		s.openTag.attrCondicionales = append(s.openTag.attrCondicionales, token.Txt)
		return
	}

	// Puede ser el último de las sentencias condicionales.
	if strings.Contains(token.Txt, "end") {
		s.openTag.attrCondicionales = append(s.openTag.attrCondicionales, token.Txt)

		// Si es una condición simple entonces en una línea.
		if len(s.openTag.attrCondicionales) == 3 {
			s.openTag.attr = append(s.openTag.attr,
				strings.Join(s.openTag.attrCondicionales, ""))
		} else {
			// Si son más de 3 poner en varias líneas con indentación.
			s.openTag.attr = append(s.openTag.attr,
				strings.Join(s.openTag.attrCondicionales, "\n"+
					strings.Repeat("\t", token.Indent)))
		}
		s.openTag.attrCondicionales = nil

	} else {
		s.openTag.attrCondicionales = append(s.openTag.attrCondicionales, token.Txt)
	}
}
