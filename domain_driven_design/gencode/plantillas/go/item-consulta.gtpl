{{ with .AgregadoConsulta -}}
package {{ .Paquete.Nombre }}

// {{ .Consulta.NombreItem }} corresponde a una consulta de solo lectura.
type {{ .Consulta.NombreItem }} struct {
{{- range .Campos }}
	{{ .NombreCampo }} {{ .TipoGo }}
	{{- if .CampoID }} // `{{ .OrigenTabla.NombreRepo }}.{{ .OrigenCampo.NombreColumna }}`
	{{- else }} // `{{ .Expresion }}`
	{{- end }}
	{{- if .Descripcion }} {{ .Descripcion }}{{ end }}
{{- end }}
}
{{ end }}