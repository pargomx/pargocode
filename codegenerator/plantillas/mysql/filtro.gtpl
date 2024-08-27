{{ define "filtro" }}{{/* with .Campo */}}

{{ if .EsBool -}}
if filtros.{{ .NombreCampo }} != nil {
	argumentos = append(argumentos, *filtros.{{ .NombreCampo }})
	where += " AND {{ if .EsConsulta }}{{ .Expresion }}{{ else }}{{ .NombreColumna }}{{ end }} = ?"
}

{{- else -}}
// Filtro {{ .NombreCampo }}
if len(filtros.{{ .NombreCampo }}) == 1 {

	{{- if .EsPropiedadExtendida }}
	if !filtros.{{ .NombreCampo }}[0].Es({{ .TipoGo }}Todos) {
	{{- else if .EsNumero }}
	if filtros.{{ .NombreCampo }}[0] != 0 {
	{{- else if .EsString }}
	if filtros.{{ .NombreCampo }}[0] != "" {
	{{- else }}
	if true { // No implementado para este tipo de dato.
	{{- end }}
		argumentos = append(argumentos, filtros.{{ .NombreCampo }}[0]{{ if .EsPropiedadExtendida }}.String{{ end }})
		where += " AND {{ if .EsConsulta }}{{ .Expresion }}{{ else }}{{ .NombreColumna }}{{ end }} = ?"
	}
} else if len(filtros.{{ .NombreCampo }}) != 0 { // Varios filtros para mismo campo...
	var where{{ .NombreCampo }} = " AND ("   // temp por si se anula el filtro.
	var args{{ .NombreCampo }} []interface{} // temp por si se anula el filtro.
	for i, {{ .Variable }} := range filtros.{{ .NombreCampo }} {

	{{- if .EsPropiedadExtendida }}
		if {{ .Variable }}.Es({{ .TipoGo }}Todos) {
	{{- else if .EsNumero }}
		if {{ .Variable }} == 0 {
	{{- else if .EsString }}
		if {{ .Variable }} == "" {
	{{- else }}
		if _FALSE_ { // No implementado para este tipo de dato.
	{{- end }}
			break // Cualquier aparición anula el filtro.
		}
		args{{ .NombreCampo }} = append(args{{ .NombreCampo }}, {{ .Variable }}{{ if .EsPropiedadExtendida }}.String{{ end }}) // Argumento para Query
		where{{ .NombreCampo }} += "{{ if .EsConsulta }}{{ .Expresion }}{{ else }}{{ .NombreColumna }}{{ end }} = ?"
		if i != len(filtros.{{ .NombreCampo }})-1 { // No es el último.
			where{{ .NombreCampo }} += " OR "
			continue // Seguimos en el mismo campo.
		}
		// Terminamos con este campo.
		argumentos = append(argumentos, args{{ .NombreCampo }}...)
		where += where{{ .NombreCampo }} + ")"
	}
}
{{- end -}}
{{- end -}}