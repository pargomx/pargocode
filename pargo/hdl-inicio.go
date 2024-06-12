package main

import (
	"monorepo/gecko"
)

func getInicio(c *gecko.Context) error {
	data := map[string]any{
		"Titulo": "Pargo 🐟",
	}
	return c.RenderOk("app/inicio", data)
}
