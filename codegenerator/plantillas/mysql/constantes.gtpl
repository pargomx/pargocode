{{ with $.TablaOrConsulta -}}
// Lista de columnas separadas por coma para usar en consulta SELECT
// en conjunto con scanRow o scanRows, ya que las columnas coinciden
// con los campos escaneados.
const columnas{{ .NombreItem }} string = "{{ .CamposAsSnakeList ", " }}"

// Origen de los datos de {{ .Paquete.Nombre }}.{{ .NombreItem }}
//
// {{ .SqlFromClause "\n// "}}
const from{{ .NombreItem }} string = "{{ .SqlFromClause " " }}"

{{ if .SqlGroupClause " " }}
// Agregación de los resultados por estas columnas.
const group{{ .NombreItem }} string = "{{ .SqlGroupClause " " }}"
{{ end }}

{{- end }}
