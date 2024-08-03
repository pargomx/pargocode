package {{ if .TablaOrConsulta.Sqlite }}sqlite{{ else }}mysql{{ end }}{{ .TablaOrConsulta.Paquete.Nombre }}
