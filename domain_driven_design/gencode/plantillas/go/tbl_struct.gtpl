package {{ $.Tabla.Paquete.Nombre }}

// {{ $.Tabla.NombreItem }} corresponde a un elemento de la tabla '{{ $.Tabla.NombreRepo }}'.
type {{ $.Tabla.NombreItem }} struct {
{{- range $.Tabla.Campos }}
	{{ .NombreCampo  }} {{ .TipoGo -}}
	{{/* comentarios */}} // `{{ $.Tabla.NombreRepo }}.{{ .NombreColumna }}` {{ if .Descripcion }} {{ .Descripcion }}{{ end }}
{{- end }}
}