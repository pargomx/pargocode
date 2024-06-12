package gk{{$.Tabla.Nombre.Clave}}

import (

)

// ================================================================ //
// ====== C R E A R =============================================== //

func (s *handlersGecko) Crear(c *gecko.Context) error {

	{{$.Tabla.Nombre.Camel}} := {{$.Tabla.Nombre.Clave}}.{{$.Tabla.NombreItem}}{}

	err := s.PrepararItemDesdePostForm(c, &{{$.Tabla.Nombre.Camel}})
	if err != nil {
		return c.ErrBadRequest(err)
	}

	err = s.escritor.InsertItem({{$.Tabla.Nombre.Camel}})
	if err != nil {
		return c.ErrBadRequest(err)
	}

	{{ if len $.Tabla.UniqueKeys -}}
	{{ $clave := index $.Tabla.UniqueKeys 0 -}}
	avelog.Evento("{{$.Tabla.Nombre.Camel}} nuev{{$.Tabla.Nombre.LetraGen}}: " + {{$.Tabla.Nombre.Camel}}.{{$clave.Nombre.Camel}})
	return c.Redir("/{{$.Tabla.NombreColumnaPlural}}/%v", {{$.Tabla.Nombre.Camel}}.{{$clave.Nombre.Camel}})
	{{ else -}}
	avelog.Evento("{{$.Tabla.Nombre.Camel}} nuev{{$.Tabla.Nombre.LetraGen}}: ")
	return c.Redir("/{{$.Tabla.NombreColumnaPlural}}/%v", {{$.Tabla.Nombre.Camel}}.{{(index $.Tabla.Campos 0).Nombre.Camel }})
	{{ end -}}
}

// ================================================================ //
// ====== A C T U A L I Z A R ===================================== //

func (s *handlersGecko) Actualizar(c *gecko.Context) error {
	{{- range $.Tabla.PrimaryKeys }}
	{{ $eFK := (index $.TablasFK .Nombre.Camel) -}}
	{{ if len $eFK.Modelos -}}
	{{ $mFK := $eFK.PrimerModelo -}}
	{{ $cFK := index $mFK.UniqueKeys 0 -}}
	{{$eFK.Nombre.Camel}}, err := s.{{$eFK.Nombre.ClavePlural}}.FetchItemBy{{$cFK.Nombre.Camel}}(c.FormVal("{{$eFK.Nombre.Clave}}"))
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{- end }}
	{{- end }}
	
	{{$FirstPK := index $.Tabla.PrimaryKeys 0}}
	{{/* Usar su clave única si la tiene */}}
	{{ if len $.Tabla.UniqueKeys }}
	{{ $cClave := index $.Tabla.UniqueKeys 0 -}}
	{{$.Tabla.Nombre.Camel}}, err := s.lector.Fetch{{$.Tabla.NombreItem}}By{{$cClave.Nombre.Camel}}(c.FormVal("{{$.Tabla.Nombre.Clave}}"))
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{/* O si solo tiene una sola PK y es numérica ... */}}
	{{ else if and (eq (len $.Tabla.PrimaryKeys) 1) ($FirstPK.EsNumero) }}
	{{$.ParamsIDs}}, err := c.FormIntMust("{{$FirstPK.NombreColumna}}")
	if err != nil {
		return c.ErrBadRequest(err)
	}
	{{$.Tabla.Nombre.Camel}}, err := s.lector.Fetch{{$.Tabla.NombreItem}}ByID({{$.ParamsIDs}})
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{/* O si solo tiene una sola PK, probablemente string... */}}
	{{ else if (eq (len $.Tabla.PrimaryKeys) 1) }}
	{{$.Tabla.Nombre.Camel}}, err := s.lector.Fetch{{$.Tabla.NombreItem}}ByID(c.FormVal("{{$FirstPK.NombreColumna}}"))
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{/* O si sabemos qué usar, dejar parámetros fáciles de popular. */}}
	{{ else }}
	{{$.Tabla.Nombre.Camel}}, err := s.lector.Fetch{{$.Tabla.NombreItem}}ByID{{if gt (len $.Tabla.PrimaryKeys) 1}}s{{end}}({{$.ParamsIDs}})
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{ end }}{{/* termina FetchItemBy... */}}
	
	err = s.PrepararItemDesdePostForm(c, {{$.Tabla.Nombre.Camel}})
	if err != nil {
		return c.ErrBadRequest(err)
	}

	err = s.escritor.UpdateItem(*{{$.Tabla.Nombre.Camel}})
	if err != nil {
		return c.ErrBadRequest(err)
	}

	{{ if len $.Tabla.UniqueKeys -}}
	{{ $clave := index $.Tabla.UniqueKeys 0 -}}
	avelog.Evento("{{$.Tabla.Nombre.Camel}} actualizad{{$.Tabla.Nombre.LetraGen}}: " + {{$.Tabla.Nombre.Camel}}.{{$clave.Nombre.Camel}})
	return c.Redir("/{{$.Tabla.NombreColumnaPlural}}/%v", {{$.Tabla.Nombre.Camel}}.{{$clave.Nombre.Camel}})
	{{ else -}}
	avelog.Evento("{{$.Tabla.Nombre.Camel}} actualizad{{$.Tabla.Nombre.LetraGen}}: ")
	return c.Redir("/{{$.Tabla.NombreColumnaPlural}}/%v", {{$.Tabla.Nombre.Camel}}.{{(index $.Tabla.Campos 0).Nombre.Camel }})
	{{ end -}}
}


// ================================================================ //
// ================================================================ //

// PrepararItemDesdePostForm popula un item
// con los datos proporcionados en el post form.
func (s *handlersGecko) PrepararItemDesdePostForm(c *gecko.Context, {{$.Tabla.NombreAbrev}} *{{$.Tabla.Nombre.Clave}}.{{$.Tabla.NombreItem}}) (err error) {

	{{ range $.Tabla.CamposEditables -}}
	
	{{ if eq .Tipo "string" -}}
	{{$.Tabla.NombreAbrev}}.{{.Nombre.Camel}} = c.FormVal("{{.NombreColumna}}")

	{{ else if eq .Tipo "int" -}}
	{{$.Tabla.NombreAbrev}}.{{.Nombre.Camel}}, err = c.FormIntMust("{{.NombreColumna}}")
	if err != nil {
		return err
	}

	{{ else if eq .Tipo "*time.Time" -}}
	{{$.Tabla.NombreAbrev}}.{{.Nombre.Camel}}, err = c.FormTime{{if .Null}}Nullable{{end}}("{{.NombreColumna}}", {{if eq .TimeTipo "date"}}"2006-01-02"{{else}}"????"{{end}})
	if err != nil {
		return averr.E(err, "{{.NombreColumna}} debe tener formato {{if eq .TimeTipo "date"}}AAAA-MM-DD{{else}}????{{end}}")
	}

	{{ else if .Especial -}}
	{{$.Tabla.NombreAbrev}}.{{.Nombre.Camel}} = {{$.Tabla.Nombre.Clave}}.Set{{.Nombre.Camel}}DB(c.FormVal("{{.NombreColumna}}"))
	

	{{ else }}
	// {{$.Tabla.NombreAbrev}}.{{.Nombre.Camel}} = c.FormVal("{{.NombreColumna}}") ???

	{{ end -}}
	{{ end }}{{/* range $.Tabla.CamposEditables */}}
	return nil
}

// ========================================================================== //
// ====== E L I M I N A R =================================================== //

// Eliminar borra el registro permanentemente.
func (s *handlersGecko) Eliminar(c *gecko.Context) error {
	{{- range $.Tabla.PrimaryKeys }}
	{{ $eFK := (index $.TablasFK .Nombre.Camel) -}}
	{{ if len $eFK.Modelos -}}
	{{ $mFK := $eFK.PrimerModelo -}}
	{{ $cFK := index $mFK.UniqueKeys 0 -}}
	{{$eFK.Nombre.Camel}}, err := s.{{$eFK.Nombre.ClavePlural}}.FetchItemBy{{$cFK.Nombre.Camel}}(c.FormVal("{{$eFK.Nombre.Clave}}"))
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{- end }}
	{{- end }}

	{{$FirstPK := index $.Tabla.PrimaryKeys 0}}
	{{/* Usar su clave única si la tiene */}}
	{{ if len $.Tabla.UniqueKeys }}
	{{ $cClave := index $.Tabla.UniqueKeys 0 -}}
	{{$.Tabla.Nombre.Camel}}, err := s.lector.Fetch{{$.Tabla.NombreItem}}By{{$cClave.Nombre.Camel}}(c.FormVal("{{$.Tabla.Nombre.Clave}}"))
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{/* O si solo tiene una sola PK y es numérica ... */}}
	{{ else if and (eq (len $.Tabla.PrimaryKeys) 1) ($FirstPK.EsNumero) }}
	{{$.ParamsIDs}}, err := c.FormIntMust("{{$FirstPK.NombreColumna}}")
	if err != nil {
		return c.ErrBadRequest(err)
	}
	{{$.Tabla.Nombre.Camel}}, err := s.lector.Fetch{{$.Tabla.NombreItem}}ByID({{$.ParamsIDs}})
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{/* O si solo tiene una sola PK, probablemente string... */}}
	{{ else if (eq (len $.Tabla.PrimaryKeys) 1) }}
	{{$.Tabla.Nombre.Camel}}, err := s.lector.Fetch{{$.Tabla.NombreItem}}ByID(c.FormVal("{{$FirstPK.NombreColumna}}"))
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{/* O si sabemos qué usar, dejar parámetros fáciles de popular. */}}
	{{ else }}
	{{$.Tabla.Nombre.Camel}}, err := s.lector.Fetch{{$.Tabla.NombreItem}}ByID{{if gt (len $.Tabla.PrimaryKeys) 1}}s{{end}}({{$.ParamsIDs}})
	if err != nil {
		return c.ErrNotFound(err)
	}
	{{ end }}{{/* termina FetchItemBy... */}}

	err = s.escritor.DeleteItem({{range $.Tabla.PrimaryKeys}}{{$.Tabla.Nombre.Camel}}.{{.Nombre.Camel}},{{end}})
	if err != nil {
		return c.ErrBadRequest(err)
	}

	{{ if len $.Tabla.UniqueKeys -}}
	{{ $clave := index $.Tabla.UniqueKeys 0 -}}
	avelog.Evento("{{$.Tabla.Nombre.Camel}} eliminad{{$.Tabla.Nombre.LetraGen}}: " + {{$.Tabla.Nombre.Camel}}.{{$clave.Nombre.Camel}})
	{{ else -}}
	avelog.Evento("{{$.Tabla.Nombre.Camel}} eliminad{{$.Tabla.Nombre.LetraGen}}: ")
	{{ end -}}
	return c.Redir("/{{$.Tabla.Nombre.KebabPlural}}")
}