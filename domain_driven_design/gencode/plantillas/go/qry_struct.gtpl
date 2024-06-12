package {{ $.Consulta.Paquete.Nombre }}

// {{ $.Consulta.NombreItem }} corresponde a una consulta de solo lectura.
type {{ $.Consulta.NombreItem }} struct {
{{- range $.Consulta.Campos }}
	//  `{{ .Expresion }}`
	{{ .NombreCampo }} {{ .TipoGo }}
{{- end }}
}