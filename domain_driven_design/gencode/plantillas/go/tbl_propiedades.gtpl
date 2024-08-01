{{ if .PackageDef }}package {{ .Tabla.Paquete.Nombre }}{{br}}{{ end }}

{{- range .Tabla.CamposEspeciales -}}
{{- $tipoGo := .TipoGo }}
{{- $nombreCampo := .NombreCampo }}
{{ separador .NombreHumano }}
// Enumeración
type {{ .TipoGo }} struct {
	ID          int
	String      string
	Filtro      string
	Etiqueta    string
	Descripcion string
}

var (
	// {{ .TipoGo }}Todos solo se utiliza como filtro.
	{{ .TipoGo }}Todos = {{ .TipoGo }}{
		ID:          -1,
		String:      "",
		Filtro:      "todos",
		Etiqueta:    "Todos",
		Descripcion: "Todos los valores posibles para {{ .NombreHumano }}",
	}
	// Indica explícitamente que la propiedad está indefinida.
	{{ .TipoGo }}Indefinido = {{ .TipoGo }}{
		ID:          0,
		String:      "",
		Filtro:      "sin_{{ .NombreColumna }}",
		Etiqueta:    "Indefinido",
		Descripcion: "Indefinido",
	}

{{range $i, $v := .ValoresPosibles}}
	// {{ .Etiqueta }}{{ if .Descripcion }}: {{ .Descripcion }}{{ end }}
	{{ $tipoGo }}{{ .Camel }} = {{ $tipoGo }}{
		ID:          {{ suma $i 1 }},
		String:      "{{ .Clave }}",
		Filtro:      "{{ .ClaveLower }}",
		Etiqueta:    "{{ .Etiqueta }}",
		Descripcion: "{{ if .Descripcion }}{{ .Descripcion }}{{ else }}{{ .Etiqueta }}{{ end }}",
	}
{{- end}}
)

// Enumeración excluyendo {{ .TipoGo }}Todos
var Lista{{ $tipoGo }} = []{{ $tipoGo }}{
	{{ $tipoGo }}Indefinido,
{{ range .ValoresPosibles }}
	{{ $tipoGo }}{{ .Camel }},
{{- end}}
}

// Enumeración incluyendo {{ .TipoGo }}Todos
var ListaFiltro{{ $tipoGo }} = []{{ $tipoGo }}{
	{{ $tipoGo }}Todos,
	{{ $tipoGo }}Indefinido,
{{ range .ValoresPosibles }}
	{{ $tipoGo }}{{ .Camel }},
{{- end}}
}

// Comparar un {{ .NombreHumano }} con otro.
func (a {{ $tipoGo }}) Es(e {{ $tipoGo }}) bool {
	return a.ID == e.ID
}

func (e {{ $tipoGo }}) EsTodos() bool {
	return e.ID == {{ $tipoGo }}Todos.ID
}
func (e {{ $tipoGo }}) EsIndefinido() bool {
	return e.ID == {{ $tipoGo }}Indefinido.ID
}
{{- range .ValoresPosibles }}
func (e {{ $tipoGo }}) Es{{ .Camel }}() bool {
	return e.ID == {{ $tipoGo }}{{ .Camel }}.ID
}
{{- end}}


// Recibe la forma .String
func Set{{ .TipoGo }}DB(str string) {{ .TipoGo }} {
	for _, e := range Lista{{ .TipoGo }} {
		if e.String == str {
			return e
		}
	}
	if str == {{ .TipoGo }}Todos.String {
		gko.LogWarn("{{ $.Tabla.Paquete.Nombre }}.Set{{ .TipoGo }}: {{ .TipoGo }}Todos es inválido en DB")
		return {{ .TipoGo }}Indefinido
	}
	gko.LogWarnf("{{ $.Tabla.Paquete.Nombre }}.Set{{ .TipoGo }} inválido: '%v'", str)
	return {{ .TipoGo }}Indefinido
}

// Recibe la forma .Filtro
func Set{{ .TipoGo }}Filtro(str string) {{ .TipoGo }} {
	if str == "" || str == {{ .TipoGo }}Todos.Filtro {
		return {{ .TipoGo }}Todos
	}
	for _, e := range Lista{{ .TipoGo }} {
		if e.Filtro == str {
			return e
		}
	}
	gko.LogWarnf("{{ $.Tabla.Paquete.Nombre }}.Set{{ .TipoGo }} inválido: '%v'", str)
	return {{ .TipoGo }}Indefinido
}


{{/* --- Métodos Modelo.SetPropiedad() ------------ */}}
// Recibe la forma .String o .Filtro
func (i *{{ $.Tabla.NombreItem }}) Set{{ $nombreCampo }}(str string) {
	for _, e := range Lista{{ .TipoGo }} {
		if e.String == str {
			i.{{ $nombreCampo }} = e
			return
		}
		if e.Filtro == str {
			i.{{ $nombreCampo }} = e
			return
		}
	}
	if str == {{ .TipoGo }}Todos.String {
		gko.LogWarn("{{ $.Tabla.Paquete.Nombre }}.Set{{ .TipoGo }}: {{ .TipoGo }}Todos es inválido en DB")
		i.{{ $nombreCampo }} = {{ .TipoGo }}Indefinido
	}
	gko.LogWarnf("{{ $.Tabla.Paquete.Nombre }}.Set{{ .TipoGo }} inválido: '%v'", str)
	i.{{ $nombreCampo }} = {{ .TipoGo }}Indefinido
}
{{- end -}}{{/* Rango campos de modelo */}}