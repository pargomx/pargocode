SELECT
  {{ .Consulta.CamposAsSnakeList ",\n  " }}
{{ .Consulta.SqlFromClause "\n"}};