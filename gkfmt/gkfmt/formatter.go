package gkfmt

type element struct {
	tag         string   // "button", "a", "span", "script", etc.
	attr        []string // href="..." or required or {{ if .Something }}required{{ end }}
	textContent string   // inner text, or <!-- comment -->, or {{ .Something }}
	indent      int      // nivel de indentación
}

type Element struct {
	tag         string   // "button", "a", "span", "script", etc.
	attr        []string // href="..." or required or {{ if .Something }}required{{ end }}
	textContent string   // inner text, or <!-- comment -->, or {{ .Something }}
	indent      int      // nivel de indentación
}

func (f *parser) Elements() {
}
