package migraciones

import "embed"

//go:embed **/*.sql
var MigracionesFS embed.FS
