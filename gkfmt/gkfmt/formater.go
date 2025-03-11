package gkfmt

func (t token) EsSelfClosingTag() bool {
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

func (f *Formatter) Format() {
	indent := 0
	for i := range f.tokens {
		switch f.tokens[i].tipo {
		case tipoOpenTagBeg:
			f.tokens[i].Indent = indent
			indent++

		case tipoAtributo, tipoOpenTagEnd, tipoGoTemplate:
			f.tokens[i].Indent = indent

		default:
			f.tokens[i].Indent = indent
		}
	}
}
