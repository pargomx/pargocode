package main

import (
	"monorepo/dpaquete"

	"github.com/pargomx/gecko"
)

func (s *servidor) formCampo(c *gecko.Context) error {
	cam, err := s.ddd.GetCampo(c.PathInt("campo_id"))
	if err != nil {
		return err
	}
	return c.Render(200, "ddd/campo-form", cam)
}

func (s *servidor) getEnumCampo(c *gecko.Context) error {
	cam, err := s.ddd.GetCampo(c.PathInt("campo_id"))
	if err != nil {
		return err
	}
	enums, err := s.ddd.GetValoresEnum(cam.CampoID)
	if err != nil {
		return err
	}
	campoConEnums := dpaquete.CampoTabla{
		Campo:           *cam,
		ValoresPosibles: enums,
	}
	return c.Render(200, "ddd/campo-enums", campoConEnums)
}
