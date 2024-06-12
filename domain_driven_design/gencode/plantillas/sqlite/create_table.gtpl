CREATE TABLE {{ .Tabla.NombreRepo }} (
{{- range .Tabla.Campos }}
  {{ .NombreColumna }}
  	{{- if .EsString }} TEXT{{ else if .EsNumero }} INT{{ end }}
	{{- if .Null}}{{ else }} NOT NULL{{ end }},
{{- end }}
  PRIMARY KEY ( {{- $lenPKs := len .Tabla.PrimaryKeys -}}
	{{- range $i, $v := .Tabla.PrimaryKeys -}}
		{{ .NombreColumna }}
		{{- $i = suma $i 1 -}}
		{{- if ne $i $lenPKs -}}
			,
		{{- end -}}
	{{- end -}} )

{{- range .Tabla.UniqueKeys }},
  UNIQUE ({{ .NombreColumna }})
{{- end }}

{{- range .Tabla.ForeignKeys }},
  FOREIGN KEY ({{ .NombreColumna }}) REFERENCES {{ .TablaFK.NombreRepo }} ({{ .CampoFK.NombreColumna }}) ON DELETE RESTRICT ON UPDATE CASCADE
{{- end }}
);