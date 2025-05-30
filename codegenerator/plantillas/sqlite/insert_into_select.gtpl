{{ with .Tabla -}}
INSERT INTO main.{{ .NombreRepo }} (
  {{ .CamposAsSnakeList ",\n  " }}
) SELECT
  {{ .CamposAsSnakeList ",\n  " }}
  FROM old_schema.{{ .Tabla.NombreRepo }}
;
{{ end }}