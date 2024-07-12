package main

import (
	"github.com/pargomx/gecko"
)

func getInicio(c *gecko.Context) error {
	data := map[string]any{
		"Titulo": "Pargo 🐟",
	}
	return c.RenderOk("app/inicio", data)
}
