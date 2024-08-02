package codegenerator

import (
	"embed"
	"monorepo/tmplutils"

	"github.com/pargomx/gecko/gko"
)

//go:embed plantillas
var plantillasFS embed.FS

type Generador struct {
	renderer *tmplutils.Renderer
	db       Repositorio
}

func NuevoGeneradorDeCodigo(repo Repositorio) (*Generador, error) {
	renderer, err := tmplutils.NuevoRenderer(plantillasFS, "plantillas")
	if err != nil {
		return nil, gko.Err(err).Op("NuevoGeneradorDeCodigo")
	}
	if repo == nil {
		return nil, gko.ErrInesperado().Str("repo es nil").Op("NuevoGeneradorDeCodigo")
	}
	return &Generador{
		renderer: renderer,
		db:       repo,
	}, nil

}
