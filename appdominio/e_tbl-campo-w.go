package appdominio

import (
	"errors"
	"monorepo/ddd"
	"monorepo/textutils"
	"strings"

	"github.com/pargomx/gecko/gko"
)

func ReordenarCampo(campoID int, newPosicion int, repo Repositorio) error {
	cam, err := repo.GetCampo(campoID)
	if err != nil {
		return err
	}
	err = repo.ReordenarCampo(cam, newPosicion)
	if err != nil {
		return err
	}
	return nil
}

func validarCampo(cam *ddd.Campo, repo Repositorio) error {
	op := gko.Op("validarCampo")

	// Debe pertenecer a una tabla
	if cam.TablaID == 0 {
		return op.Msg("TablaID sin especificar")
	}
	_, err := repo.GetTabla(cam.TablaID)
	if err != nil {
		return err
	}

	// Campos obligatorios
	if cam.CampoID == 0 {
		return op.Msg("CampoID sin especificar")
	}
	if cam.NombreColumna == "" {
		return op.Msg("NombreColumna sin especificar")
	}
	if cam.TipoSql == "" {
		return op.Msg("TipoSql sin especificar")
	}
	if cam.NombreCampo == "" {
		return op.Msg("NombreCampo sin especificar")
	}
	if cam.TipoGo == "" {
		return op.Msg("TipoGo sin especificar")
	}
	if cam.NombreHumano == "" {
		return op.Msg("NombreHumano sin especificar")
	}

	// Si es clave foránea, traer la referencia.
	if cam.ForeignKey {
		fk, err := repo.GetCampoPrimaryKey(cam.NombreColumna)
		if err != nil {
			return err
		}
		cam.ReferenciaCampo = &fk.CampoID
	} else {
		cam.ReferenciaCampo = nil
	}

	// Si tiene valores enum entonces es un campo especial.
	// if len(cam.ValoresEnum) > 0 {
	// 	cam.Especial = true
	// }

	// Handy defaults
	if cam.DefaultSql == "" && cam.NotNull() && !cam.PrimaryKey && !cam.Required() &&
		(cam.EsSqlChar() || cam.EsSqlVarchar()) {
		cam.DefaultSql = "DEFAULT ''"
	}
	if cam.DefaultSql == "AI" || cam.DefaultSql == "AUTOINCREMENT" {
		cam.DefaultSql = "AUTO_INCREMENT"
	}
	if cam.DefaultSql == "NULL" {
		cam.DefaultSql = "DEFAULT NULL"
	}
	if cam.DefaultSql == "0" || cam.DefaultSql == "FALSE" {
		cam.DefaultSql = "DEFAULT '0'"
	}
	// if cam.Nombre.Kebab == "fecha-registro" {
	// 	cam.DefaultSql = "DEFAULT CURRENT_TIMESTAMP"
	// 	cam.TipoSQL = "timestamp"
	// 	cam.ReadOnly = true
	// }
	// if cam.Nombre.Kebab == "fecha-modif" {
	// 	cam.TipoSQL = "timestamp"
	// 	cam.DefaultSql = "DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"
	// 	cam.ReadOnly = true
	// }
	if cam.DefaultSql == "-" {
		cam.DefaultSql = ""
	}

	// Sólo los enteros pueden ser unsigned
	if !cam.EsSqlInt() {
		cam.Uns = false
	}

	if cam.EsSqlChar() && cam.MaxLenght == 0 {
		cam.MaxLenght = 15
	}
	if cam.EsSqlVarchar() && cam.MaxLenght == 0 {
		cam.MaxLenght = 240
	}

	// TipoGo desde TipoSQL
	if cam.TipoGo == "" {
		switch cam.TipoSql {

		case "tinyint", "smallint", "mediumint", "int":
			cam.TipoGo = "int"

		case "bigint":
			if cam.Unsigned() {
				cam.TipoGo = "uint64"
			} else {
				cam.TipoGo = "int" // int es 8 bytes en sistemas de 64bits.
			}

		case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
			cam.TipoGo = "string"

		case "timestamp", "datetime", "date", "time":
			cam.TipoGo = "time.Time"
			// cam.TimeTipo = cam.TipoSQL

		case "year":
			cam.TipoGo = "string"
		}

		if cam.Nullable {
			cam.TipoGo = "*" + cam.TipoGo
		}
	}

	// Nullable time debe ser default null.
	// if cam.Null && cam.TimeTipo != "" {
	// 	cam.DefaultSql = "DEFAULT NULL"
	// }

	// Validar que no haya otro campo con el mismo nombre en la misma tabla
	campos, err := repo.ListCamposByTablaID(cam.TablaID)
	if err != nil {
		return op.Err(err)
	}
	for _, c := range campos {
		if c.CampoID == cam.CampoID {
			continue
		}
		if c.NombreColumna == cam.NombreColumna {
			return op.Msgf("Ya existe un campo con el nombre '%s' en la tabla", cam.NombreColumna)
		}
		if c.NombreCampo == cam.NombreCampo {
			return op.Msgf("Ya existe un campo con el nombre '%s' en la tabla", cam.NombreCampo)
		}
	}
	return nil
}

// Insertar campo con defaults basado solamente en el nombre del campo.
func InsertarCampoQuick(tablaID int, nombreCol string, repo Repositorio) error {
	op := gko.Op("InsertarCampoQuick")
	if nombreCol == "" {
		return op.Msg("Debe especificar un nombre para el nuevo campo")
	}

	nombreCol = textutils.QuitarAcentos(nombreCol)
	nombreCol = strings.ReplaceAll(nombreCol, " ", "_")
	nombreCol = strings.ReplaceAll(nombreCol, "-", "_")
	nombreCol = strings.ToLower(nombreCol)

	// Default para campos es texto
	cam := ddd.Campo{
		CampoID:       ddd.NewCampoID(),
		TablaID:       tablaID,
		NombreColumna: nombreCol,
		TipoSql:       "varchar",
		NombreCampo:   textutils.SnakeToCamel(nombreCol),
		TipoGo:        "string",
		NombreHumano:  textutils.SnakeToCamel(nombreCol),
	}

	// Default para IDs
	if strings.HasSuffix(nombreCol, "id") {
		cam.TipoSql = "bigint"
		cam.Uns = true
		cam.PrimaryKey = true
		cam.Nullable = false
		cam.TipoGo = "uint64"
	}

	// Default para fechas
	if strings.Contains(nombreCol, "fecha") {
		cam.TipoSql = "timestamp"
		cam.Nullable = false
		cam.TipoGo = "time.Time"
	}

	// Default basado en un campo con el mismo nombre
	similar, err := repo.GetCampoByNombre(nombreCol)
	if err != nil {
		gko.LogError(err)
	} else {
		cam.NombreCampo = similar.NombreCampo
		cam.NombreColumna = similar.NombreColumna
		cam.NombreHumano = similar.NombreHumano
		cam.TipoGo = similar.TipoGo
		cam.TipoSql = similar.TipoSql
		cam.Setter = similar.Setter
		cam.Importado = similar.Importado
		cam.Uq = similar.Uq
		cam.Req = similar.Req
		cam.Ro = similar.Ro
		cam.Filtro = similar.Filtro
		cam.Nullable = similar.Nullable
		cam.MaxLenght = similar.MaxLenght
		cam.Uns = similar.Uns
		cam.DefaultSql = similar.DefaultSql
		cam.Especial = similar.Especial
		cam.ReferenciaCampo = similar.ReferenciaCampo
		cam.Expresion = similar.Expresion
		cam.EsFemenino = similar.EsFemenino
		cam.Descripcion = similar.Descripcion
		if similar.PrimaryKey {
			cam.ForeignKey = true
			cam.ReferenciaCampo = &similar.CampoID
		}
	}

	err = InsertarCampo(cam, repo)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func InsertarCampo(cam ddd.Campo, repo Repositorio) error {
	op := gko.Op("InsertarCampo")
	err := validarCampo(&cam, repo)
	if err != nil {
		return op.Err(err)
	}
	err = repo.InsertCampo(cam)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func ActualizarCampo(campoID int, new ddd.Campo, repo Repositorio) error {
	op := gko.Op("UpdateCampo")
	cam, err := repo.GetCampo(campoID)
	if err != nil {
		return op.Err(err)
	}

	if cam.CampoID != new.CampoID {
		return op.Msgf("CampoID no coincide: old=%v new%v", cam.CampoID, new.CampoID)
	}

	cam.CampoID = new.CampoID
	// cam.TablaID = new.TablaID // No se puede cambiar la tabla
	cam.NombreColumna = new.NombreColumna
	cam.TipoSql = new.TipoSql
	cam.DefaultSql = new.DefaultSql
	cam.NombreCampo = new.NombreCampo
	cam.TipoGo = new.TipoGo
	cam.Setter = new.Setter
	cam.NombreHumano = new.NombreHumano
	cam.Descripcion = new.Descripcion
	cam.Nullable = new.Nullable
	cam.Uns = new.Uns
	cam.MaxLenght = new.MaxLenght
	cam.PrimaryKey = new.PrimaryKey
	cam.ForeignKey = new.ForeignKey
	cam.Uq = new.Uq
	cam.Req = new.Req
	cam.Ro = new.Ro
	cam.Filtro = new.Filtro
	cam.Especial = new.Especial

	err = validarCampo(cam, repo)
	if err != nil {
		return op.Err(err)
	}
	err = repo.UpdateCampo(*cam)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func EliminarCampo(campoID int, repo Repositorio) error {
	cam, err := repo.GetCampo(campoID)
	if err != nil {
		return err
	}

	// TODO: eliminar todo lo que tiene que ver con el campo, como valores_enum, buscar referencias en consultas, etc.

	err = repo.DeleteCampo(cam.CampoID)
	if err != nil {
		return err
	}
	return nil
}

func ActualizarOpcionesDeCampoEnum(campoID int, nuevosValores []ddd.ValorEnum, repo Repositorio) error {
	// Debe existir el campo
	if campoID == 0 {
		return errors.New("CampoID sin especificar")
	}
	_, err := repo.GetCampo(campoID)
	if err != nil {
		return err
	}
	err = repo.GuardarValoresEnum(campoID, nuevosValores)
	if err != nil {
		return err
	}
	return nil
}

func FixPosicionDeCampos(tablaID int, repo Repositorio) error {
	campos, err := repo.ListCamposByTablaID(tablaID)
	if err != nil {
		return err
	}
	for i, c := range campos {
		if c.Posicion != i+1 {
			err := repo.ReordenarCampo(&c, i+1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
