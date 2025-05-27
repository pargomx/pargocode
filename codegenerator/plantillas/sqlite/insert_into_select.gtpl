{{ with .Tabla -}}
INSERT INTO main.{{ .NombreRepo }} (
  {{ .CamposAsSnakeList ",\n  " }}
) SELECT
  {{ .CamposAsSnakeList ",\n  " }}
  FROM old.{{ .Tabla.NombreRepo }}
;
{{ end }}