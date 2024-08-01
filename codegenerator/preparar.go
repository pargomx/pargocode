package codegenerator

import (
	"embed"
	"monorepo/tmplutils"
)

//go:embed plantillas
var plantillasFS embed.FS

type Generador struct {
	renderer *tmplutils.Renderer
}

func NuevoGeneradorDeCodigo() (*Generador, error) {
	renderer, err := tmplutils.NuevoRenderer(plantillasFS, "plantillas")
	if err != nil {
		return nil, err
	}
	return &Generador{
		renderer: renderer,
	}, nil

}
