package main

import (
	"encoding/json"
	"fmt"
	"io"
	"monorepo/appdominio"
	"monorepo/ddd"
	"monorepo/textutils"
	"strings"

	"github.com/pargomx/gecko"
)

func (s *servidor) putCampo(c *gecko.Context) error {
	err := appdominio.InsertarCampoQuick(c.PathInt("tabla_id"), c.FormVal("nombre_columna"), s.ddd)
	if err != nil {
		return err
	}
	return c.StatusOkf("Campo %v insertado", c.FormVal("campo_nuevo"))
}

func (s *servidor) postCampo(c *gecko.Context) error {
	cam := ddd.Campo{
		CampoID: ddd.NewCampoID(),
		TablaID: c.PathInt("tabla_id"),

		NombreColumna: c.FormValue("nombre_sql"),
		TipoSql:       c.FormValue("tipo_sql"),
		DefaultSql:    c.FormVal("default_sql"),

		NombreCampo: c.FormValue("nombre_go"),
		TipoGo:      c.FormValue("tipo_go"),
		Setter:      c.FormValue("setter"),

		NombreHumano: c.FormValue("nombre_ui"),
		Descripcion:  c.FormValue("descripcion"),

		Nullable:   !c.FormBool("not_null"),
		Uns:        c.FormBool("unsinged"),
		MaxLenght:  c.FormInt("max_lenght"),
		PrimaryKey: c.FormBool("pk"),
		ForeignKey: c.FormBool("fk"),
		Uq:         c.FormBool("unique"),
		Req:        c.FormBool("required"),
		Ro:         c.FormBool("readonly"),
		Filtro:     c.FormBool("filtro"),
		Especial:   c.FormBool("especial"),
	}
	err := appdominio.InsertarCampo(cam, s.ddd)
	if err != nil {
		return err
	}
	return c.StatusOkf("Campo %v insertado", cam.CampoID)
}

func (s *servidor) updateCampo(c *gecko.Context) error {
	cam := ddd.Campo{
		CampoID: c.PathInt("campo_id"),
		TablaID: c.PathInt("tabla_id"),

		NombreColumna: c.FormValue("nombre_sql"),
		TipoSql:       c.FormValue("tipo_sql"),
		DefaultSql:    c.FormVal("default_sql"),

		NombreCampo: c.FormValue("nombre_go"),
		TipoGo:      c.FormValue("tipo_go"),
		Setter:      c.FormValue("setter"),

		NombreHumano: c.FormValue("nombre_ui"),
		Descripcion:  c.FormValue("descripcion"),

		Nullable:   !c.FormBool("not_null"),
		ZeroIsNull: c.FormBool("zero_is_null"),
		Uns:        c.FormBool("unsinged"),
		MaxLenght:  c.FormInt("maxlenght"),
		PrimaryKey: c.FormBool("pk"),
		ForeignKey: c.FormBool("fk"),
		Uq:         c.FormBool("unique"),
		Req:        c.FormBool("required"),
		Ro:         c.FormBool("readonly"),
		Filtro:     c.FormBool("filtro"),
		Especial:   c.FormBool("especial"),
	}
	refCampo := c.FormInt("campo_fk")
	if refCampo > 0 {
		cam.ForeignKey = true
		cam.ReferenciaCampo = &refCampo
	}
	err := appdominio.ActualizarCampo(cam.CampoID, cam, s.ddd)
	if err != nil {
		return err
	}
	return c.StatusOkf("Campo %v actualizado", cam.CampoID)
}

func (s *servidor) reordenarCampo(c *gecko.Context) error {
	err := appdominio.ReordenarCampo(c.PathInt("campo_id"), c.FormInt("newPosicion"), s.ddd)
	if err != nil {
		return err
	}
	return c.StatusOkf("Campo %v reordenado", c.PathInt("campo_id"))
}

func (s *servidor) postEnumCampo(c *gecko.Context) error {
	type RequestCampoEspecial struct {
		Etiqueta    string `json:"etiqueta"`
		Clave       string `json:"clave"`
		Descripcion string `json:"descripcion"`
	}
	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	valores := []RequestCampoEspecial{}
	err = json.Unmarshal(data, &valores)
	if err != nil {
		return err
	}
	fmt.Println("Recibido", valores)
	nuevosValores := []ddd.ValorEnum{}
	for i, v := range valores {
		nuevosValores = append(nuevosValores, ddd.ValorEnum{
			Numero:      i + 1,
			Etiqueta:    textutils.PrimeraMayusc(v.Etiqueta),
			Clave:       textutils.QuitarAcentos(strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(v.Clave, " ", "-"), "-", "_"))),
			Descripcion: textutils.PrimeraMayusc(v.Descripcion),
		})
	}
	err = appdominio.ActualizarOpcionesDeCampoEnum(c.PathInt("campo_id"), nuevosValores, s.ddd)
	if err != nil {
		return err
	}
	return c.StatusOk("Guardado")
}

func (s *servidor) deleteCampo(c *gecko.Context) error {
	err := appdominio.EliminarCampo(c.PathInt("campo_id"), s.ddd)
	if err != nil {
		return err
	}
	return c.StatusOkf("Campo %v eliminado", c.PathInt("campo_id"))
}

func (s *servidor) fixOrdenDeCampos(c *gecko.Context) error {
	err := appdominio.FixPosicionDeCampos(c.PathInt("tabla_id"), s.ddd)
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}
