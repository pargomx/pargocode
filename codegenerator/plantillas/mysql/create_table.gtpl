CREATE TABLE `{{ .Tabla.NombreRepo }}` (
{{- range .Tabla.Campos }}
  `{{.NombreColumna}}` {{ .TipoSQL }}
	{{- if and .EsSqlInt .Unsigned }} unsigned{{ end }}
	{{- if and (or .EsSqlChar .EsSqlVarchar) (gt .MaxLenght 0) }}({{.MaxLenght}}){{ end }}
	{{- if or .EsSqlChar .EsSqlVarchar .EsSqlText }}{{ if or .PrimaryKey .ForeignKey }} CHARACTER SET ascii COLLATE ascii_general_ci{{ else }} CHARACTER SET utf8mb4 COLLATE utf8mb4_es_0900_ai_ci{{ end }}{{ end }}
	{{- if .Null}}{{ else }} NOT NULL{{ end }}
	{{- if .DefaultSQL }} {{ .DefaultSQL }}{{ end }},
{{- end }}
  PRIMARY KEY ( {{- $lenPKs := len .Tabla.PrimaryKeys -}}
	{{- range $i, $v := .Tabla.PrimaryKeys -}}
		`{{ .NombreColumna }}`
		{{- $i = suma $i 1 -}}
		{{- if ne $i $lenPKs -}}
			,
		{{- end -}}
	{{- end -}} )

{{- range .Tabla.UniqueKeys }},
  UNIQUE KEY `{{ $.Tabla.NombreAbrev }}_{{ .NombreColumna }}_UNIQUE` (`{{ .NombreColumna }}`)
{{- end }}

{{- range .Tabla.ForeignKeys }},
  CONSTRAINT `{{ .TablaFK.Abrev }}_{{ $.Tabla.NombreRepo }}` FOREIGN KEY (`{{ .NombreColumna }}`) REFERENCES `{{ .TablaFK.NombreRepo }}` (`{{ .CampoFK.NombreColumna }}`) ON DELETE RESTRICT ON UPDATE CASCADE
{{- end }}
);