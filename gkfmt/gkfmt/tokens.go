package gkfmt

// ================================================================ //
// ========== Tokens ============================================== //

type token struct {
	Txt    string // puede ser '<tag', 'attr="x"'.
	tipo   tipo
	Indent int // nivel de indentación.
}

type tipo int

func newComparador(base tipo) func(tipo) bool {
	return func(t tipo) bool {
		return base == t
	}
}

const (
	tipoTextContent  tipo = iota // Contenido que puede ser texto y ya no contiene nignún "<"
	tipoOpenTagBeg               // <button
	tipoAtributo                 // hx-get="/x", required, class="{{ if not .Show }}hidden {{ end }}"
	tipoOpenTagEnd               // >
	tipoClosingTag               // </button>
	tipoComentario               // <!-- Comentario -->
	tipoScript                   // <script> ... </script>
	tipoGoTemplate               // {{ .Something }}
	tipoExtraNewLine             // Espacio vertical intencional.
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
		"__br__",
	}[t.tipo]
}

func (t token) esSelfClosingTag() bool {
	if t.tipo != tipoOpenTagBeg {
		return false
	}
	switch t.Txt {
	case "<area",
		"<base",
		"<br",
		"<col",
		"<embed",
		"<hr",
		"<img",
		"<input",
		"<link",
		"<meta",
		"<param",
		"<source",
		"<track",
		"<wbr":
		return true
	}
	return false
}
