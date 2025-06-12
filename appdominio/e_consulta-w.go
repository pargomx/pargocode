package appdominio

import (
	"fmt"
	"monorepo/ddd"
	"monorepo/textutils"
	"strings"

	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

func validarConsulta(con ddd.Consulta, repo Repositorio, op *gko.Error) error {
	// Validar Paquete
	if con.PaqueteID == 0 {
		return op.Msg("se debe definir el paquete al que pertenece la nueva consulta")
	}
	paq, err := repo.GetPaquete(con.PaqueteID)
	if err != nil {
		return op.Msgf("no se encontró el paquete con ID %v", con.PaqueteID).Err(err)
	}

	// Validar Tabla FROM
	if con.TablaID == 0 {
		return op.Msg("se debe definir tabla FROM de la consulta")
	}
	_, err = repo.GetTabla(con.TablaID)
	if err != nil {
		return op.Msgf("no se encontró tabla con ID %v para el FROM", con.TablaID).Err(err)
	}

	// Campos obligatorios
	if con.NombreItem == "" {
		return op.Msg("se debe definir nombre Item")
	}
	if len(con.NombreItem) < 3 {
		return op.Msg("el nombre Item debe tener al menos 3 caracteres")
	}

	if con.NombreItems == "" {
		return op.Msg("se debe definir nombre Items")
	}
	if len(con.NombreItems) < 3 {
		return op.Msg("el nombre Items debe tener al menos 3 caracteres")
	}

	if con.Abrev == "" {
		return op.Msg("se debe definir abreviatura")
	}

	// No se puede tener una consulta con el mismo nombre que otra en el mismo paquete
	consultas, err := repo.ListConsultasByPaqueteID(con.PaqueteID)
	if err != nil {
		return op.Err(err).Op("verificar_nombre_unico")
	}
	for _, c := range consultas {
		if c.ConsultaID == con.ConsultaID {
			continue
		}
		if c.NombreItem == con.NombreItem {
			return op.Msgf("ya existe una consulta con el nombre '%v' en el paquete '%v'", con.NombreItem, paq.Nombre)
		}
		if c.NombreItems == con.NombreItems {
			return op.Msgf("ya existe una consulta con el nombre '%v' en el paquete '%v'", con.NombreItems, paq.Nombre)
		}
		if c.Abrev == con.Abrev {
			return op.Msgf("ya existe una consulta con la abreviatura '%v' en el paquete '%v'", con.Abrev, paq.Nombre)
		}
	}
	return nil
}

func CrearConsulta(con ddd.Consulta, repo Repositorio) error {
	op := gko.Op("CrearConsulta")
	if con.ConsultaID == 0 {
		return op.Msg("se debe proporcionar un nuevo ID para la consulta")
	}
	// Verificar que no exista una consulta con el mismo ID
	_, err := repo.GetConsulta(con.ConsultaID)
	if err == nil {
		return op.Msgf("ya existe una consulta con el ID proporcionado %v", con.ConsultaID)
	}
	err = validarConsulta(con, repo, op)
	if err != nil {
		return err
	}
	err = repo.InsertConsulta(con)
	if err != nil {
		return err
	}
	return nil
}

func ActualizarConsulta(consultaID int, new ddd.Consulta, repo Repositorio) error {
	op := gko.Op("ActualizarConsulta")
	con, err := repo.GetConsulta(consultaID)
	if err != nil {
		return err
	}

	// Contrastar nuevos valores
	if con.ConsultaID != new.ConsultaID {
		return op.Msgf("no se puede cambiar el ID de la consulta")
	}
	if con.PaqueteID != new.PaqueteID {
		op.Op("cambiar_paquete").Ctx("newPaqueteID", new.PaqueteID)
	}
	if con.TablaID != new.TablaID {
		op.Op("cambiar_tabla_from").Ctx("newTablaID", new.TablaID)
		return op.Msg("no se puede cambiar la tabla FROM de la consulta (aún :v)")
	}

	// Set nuevos datos
	con.PaqueteID = new.PaqueteID
	con.NombreItem = new.NombreItem
	con.NombreItems = new.NombreItems
	con.Abrev = new.Abrev
	con.EsFemenino = new.EsFemenino
	con.Descripcion = new.Descripcion
	con.Directrices = new.Directrices

	err = validarConsulta(*con, repo, op)
	if err != nil {
		return err
	}
	err = repo.UpdateConsulta(*con)
	if err != nil {
		return err
	}
	return nil
}

func EliminarConsulta(consultaID int, repo Repositorio) error {
	op := gko.Op("EliminarConsulta")
	con, err := repo.GetConsulta(consultaID)
	if err != nil {
		return op.Err(err)
	}
	err = repo.DeleteConsulta(con.ConsultaID)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

// ================================================================ //
// ========== RELACIONES ========================================== //

func AgregarRelacionConsulta(consultaID int, tipo string, joinTablaID int, fromAbrev string, repo Repositorio) error {
	op := gko.Op("AgregarRelacionConsulta")
	// Consulta a la que se le agrega la relación
	con, err := GetConsulta(consultaID, repo)
	if err != nil {
		return op.Msgf("no se puede cargar la consulta %v", consultaID).Err(err)
	}

	// Tabla para relacionar
	if joinTablaID == 0 {
		return op.Msg("se debe definir la tabla JOIN que se va a relacionar")
	}
	tblJoin, err := GetTabla(joinTablaID, repo)
	if err != nil {
		return op.Msgf("la tabla join %v no se encuentra", joinTablaID).Err(err)
	}

	relacion := ddd.ConsultaRelacion{
		ConsultaID:  consultaID,
		Posicion:    len(con.Relaciones) + 1,
		TipoJoin:    ddd.SetTipoJoinDB(tipo),
		JoinTablaID: tblJoin.Tabla.TablaID,
		JoinAs:      tblJoin.Tabla.Abrev, // default
		JoinOn:      "",
		FromTablaID: 0,
	}

	// Si no se especifica el tipo, se asume LEFT
	if relacion.TipoJoin.EsIndefinido() {
		relacion.TipoJoin = ddd.TipoJoinLeft
	}

	// No usar el mismo alias dos veces para la misma tabla.
	apariciones := 0
	aliasRepetido := false
	if con.TablaOrigen.Abrev == relacion.JoinAs {
		aliasRepetido = true
	}
	if con.TablaOrigen.TablaID == relacion.JoinTablaID {
		apariciones++
	}
	for _, r := range con.Relaciones {
		if r.JoinAs == relacion.JoinAs {
			aliasRepetido = true
		}
		if r.JoinTablaID == relacion.JoinTablaID {
			apariciones++
		}
	}
	if aliasRepetido && apariciones > 0 {
		relacion.JoinAs = fmt.Sprintf("%s%d", relacion.JoinAs, apariciones)
	}

	// A partir de quién (FROM alias)
	if fromAbrev == "" {
		return op.Msg("se debe definir la tabla FROM desde la que parte el JOIN")
	}
	if fromAbrev == con.TablaOrigen.Abrev {
		relacion.FromTablaID = con.TablaOrigen.TablaID
	} else {
		for _, r := range con.Relaciones {
			if r.JoinAs == fromAbrev {
				relacion.FromTablaID = r.JoinTablaID
				break
			}
		}
	}
	if relacion.FromTablaID == 0 {
		return op.Msgf("la tabla FROM con el alias '%v' no forma parte de esta consulta", fromAbrev)
	}

	// Construir ON...
	tblFrom, err := GetTabla(relacion.FromTablaID, repo)
	if err != nil {
		return op.Msgf("la tabla from %v no se existe?", relacion.FromTablaID).Err(err)
	}

	for _, cJoin := range tblJoin.PrimaryKeys() {
		cFrom, err := tblFrom.BuscarCampo(cJoin.NombreColumna) // TODO: no usar agregado sino GetCampo(idCampo)
		if err != nil {
			for _, cJoin2 := range tblJoin.ForeignKeys() {
				if cJoin.CampoID == cJoin2.CampoID {
					continue
				}
				cFrom, err = tblFrom.BuscarCampo(cJoin2.NombreColumna) // TODO: no usar agregado sino GetCampo(idCampo)
				if err != nil {
					break
				}
			}
		}
		if cFrom == nil {
			gko.LogWarn("ignorando '" + tblJoin.Tabla.NombreRepo + "->" + tblFrom.Tabla.NombreRepo + "' porque no comparten PK-FK " + cJoin.NombreColumna)
			continue
		}
		if !cJoin.ForeignKey && !cJoin.PrimaryKey {
			gko.LogWarn("El campo '" + cJoin.NombreColumna + "' de '" + tblJoin.Tabla.NombreRepo + "' no está marcado como FK o PK pero se usa en relación " + tblJoin.Tabla.NombreRepo + "->" + tblFrom.Tabla.NombreRepo)
		}
		if !cFrom.ForeignKey && !cFrom.PrimaryKey {
			gko.LogWarn("El campo '" + cFrom.NombreColumna + "' de '" + tblFrom.Tabla.NombreRepo + "' no está marcado como FK o PK pero se usa en relación " + tblJoin.Tabla.NombreRepo + "->" + tblFrom.Tabla.NombreRepo)
		}
		relacion.JoinOn += fmt.Sprintf("%s.%s = %s.%s AND ",
			relacion.JoinAs, cJoin.NombreColumna,
			fromAbrev, cFrom.NombreColumna,
		)
	}
	relacion.JoinOn = strings.TrimSuffix(relacion.JoinOn, " AND ")

	err = repo.InsertConsultaRelacion(relacion)
	if err != nil {
		return op.Msg("No se pudo insertar la relación").Err(err)
	}
	return nil
}

func EliminarRelacionConsulta(consultaID int, posicion int, repo Repositorio) error {
	op := gko.Op("EliminarRelacionConsulta").Ctx("id", consultaID)
	con, err := GetConsulta(consultaID, repo)
	if err != nil {
		return op.Msgf("no se puede cargar la consulta %v", consultaID).Err(err)
	}

	// Validar que la relación exista
	if posicion < 1 || posicion > len(con.Relaciones) {
		return op.Msgf("la posición %v no es válida para la consulta %v que tiene %v relaciones", posicion, consultaID, len(con.Relaciones))
	}
	rel := con.Relaciones[posicion-1]

	// Que no se pueda eliminar la relación si hay campos que dependen de ella.
	for _, c := range con.Campos {
		if strings.Contains(c.Expresion, rel.JoinAs+".") {
			return op.Msgf("no se puede eliminar la relación '%v' porque parece que el campo '%v' la usa", rel.JoinAs, c.NombreCampo)
		}
	}

	err = repo.DeleteRelacionConsulta(consultaID, rel.Posicion)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func ActualizarRelacionConsulta(ConsultaID int, Posicion int, newTipo string, newAs string, newOn string, repo Repositorio) error {
	op := gko.Op("ActualizarRelacionConsulta").Ctx("id", ConsultaID)
	con, err := GetConsulta(ConsultaID, repo)
	if err != nil {
		return op.Msgf("no se puede cargar la consulta %v", ConsultaID).Err(err)
	}
	// Validar que la relación exista
	if Posicion < 1 || Posicion > len(con.Relaciones) {
		return op.Msgf("no existe el JOIN %v en la consulta %v que tiene %v relaciones", Posicion, ConsultaID, len(con.Relaciones))
	}
	old := con.Relaciones[Posicion-1]

	if newAs == "" {
		return op.Msg("no se puede dejar vacío el alias de la relación")
	}
	if newOn == "" {
		return op.Msg("no se puede dejar vacío el ON de la relación")
	}

	actualizada := ddd.ConsultaRelacion{
		ConsultaID:  old.ConsultaID,
		Posicion:    old.Posicion,
		TipoJoin:    ddd.SetTipoJoinDB(newTipo),
		JoinTablaID: old.JoinTablaID,
		JoinAs:      old.JoinAs,
		JoinOn:      old.JoinOn,
		FromTablaID: old.FromTablaID,
	}

	// Si no se especifica el tipo, se asume LEFT
	if actualizada.TipoJoin.EsIndefinido() {
		actualizada.TipoJoin = ddd.TipoJoinLeft
	}

	if newAs != old.JoinAs {

		actualizada.JoinAs = newAs
		actualizada.JoinOn = strings.ReplaceAll(old.JoinOn, old.JoinAs+".", newAs+".") // sustitución automática

		for _, c := range con.Campos {
			if strings.HasPrefix(c.Expresion, old.JoinAs+".") {
				// Sugerir revisar campos calculados afectados.
				if c.CampoID == nil {
					fmt.Printf("Revise el campo '%v' porque probablemente dependa del viejo alias '%v' de la relación\n", c.NombreCampo, old.JoinAs)
					continue
				}
				// Automáticamente cambiar alias de campos afectados
				campo, err := repo.GetConsultaCampo(c.ConsultaID, c.Posicion)
				if err != nil {
					return op.Err(err).Ctx("campo", c.NombreCampo).Op("cambiarOrigenCampoNewAlias")
				}
				if c.Expresion != campo.Expresion {
					return op.Msgf("el campo '%v' tiene una expresión distinta a la que se encontró en la base de datos!!", c.NombreCampo)
				}
				campo.Expresion = newAs + strings.TrimPrefix(campo.Expresion, old.JoinAs)
				err = repo.UpdateConsultaCampo(*campo)
				if err != nil {
					return op.Err(err).Ctx("campo", c.NombreCampo).Op("updateOrigenCampoNewAlias")
				}
			}
		}
	}

	// Ignorar sustitución automática si se cambió manualmente.
	if newOn != old.JoinOn {
		actualizada.JoinOn = newOn
	}

	err = repo.UpdateConsultaRelacion(actualizada)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

// ================================================================ //
// ========== CAMPOS ============================================== //

func agregarAllCamposDeTabla(con *Consulta, fromAbrev string, repo Repositorio) error {
	op := gko.Op("agregarCamposConsultaAll").Ctx("consulta_id", con.Consulta.ConsultaID).Ctx("from", fromAbrev)

	var tblFrom *Tabla = nil
	if fromAbrev == con.From.Tabla.Abrev {
		tblFrom = &con.From
	} else {
		for _, rel := range con.Relaciones {
			if fromAbrev == rel.JoinAs {
				tblFrom = &rel.Join
				break
			}
		}
	}
	if tblFrom == nil {
		return op.Msgf("La tabla de origen '%v' no se encuentra en el FROM ni los JOINs", fromAbrev)
	}

	for _, cam := range tblFrom.Campos {
		err := AgregarCampoConsulta(con.Consulta.ConsultaID, fromAbrev, cam.NombreColumna, repo)
		if err != nil {
			gko.LogError(err)
		}
	}
	return nil
}

func AgregarCampoConsulta(consultaID int, fromAbrev string, expresion string, repo Repositorio) error {
	op := gko.Op("AgregarCampoConsulta").Ctx("consulta_id", consultaID).Ctx("from", fromAbrev).Ctx("expr", expresion)
	con, err := GetConsulta(consultaID, repo)
	if err != nil {
		return op.Err(err)
	}

	// Si se quiere agregar todos los campos de una tabla.
	if expresion == "*" {
		return agregarAllCamposDeTabla(con, fromAbrev, repo)
	}

	campo := ddd.ConsultaCampo{
		ConsultaID: con.Consulta.ConsultaID,
		Posicion:   len(con.Campos) + 1,
		Expresion:  expresion,
	}

	if fromAbrev != "" {
		// El origen de la columna debe encontrarse en el FROM / JOIN por su alias.
		var tblFrom *Tabla = nil
		if fromAbrev == con.From.Tabla.Abrev {
			tblFrom = &con.From
		} else {
			for _, rel := range con.Relaciones {
				if fromAbrev == rel.JoinAs {
					tblFrom = &rel.Join
					break
				}
			}
		}
		if tblFrom == nil {
			return op.Msgf("La tabla de origen '%v' no se encuentra en el FROM ni los JOINs", fromAbrev)
		}
		if expresion == "" {
			return op.Msgf("No se especificó qué columna de '%v %v' se va a a agregar a la consulta", tblFrom.Tabla.NombreRepo, fromAbrev)
		}
		for _, c := range tblFrom.Campos {
			if c.NombreColumna == expresion {
				campo.CampoID = &c.CampoID
				campo.Expresion = fmt.Sprintf("%s.%s", fromAbrev, c.NombreColumna)
				campo.NombreCampo = c.NombreCampo
				campo.TipoGo = c.TipoGo
				campo.Descripcion = c.Descripcion
				break
			}
		}
		if campo.CampoID == nil {
			return op.Msgf("La columna '%v' no se encuentra en la tabla de origen '%v %v'", expresion, tblFrom.Tabla.NombreRepo, fromAbrev)
		}

	} else {
		// También puede ser una expresión libre que no viene de una sola columna.
		if expresion == "" {
			return op.Msg("No se proporcionó ninguna expresión para el campo")
		}
		campo.CampoID = nil
		campo.Expresion = expresion
		campo.NombreCampo = "Calculado"
		campo.TipoGo = "string"
	}

	// Los nombres de los campos no se deben repetir.
	apariciones := 0
	seRepiteAlias := false
	seRepiteCampo := false
	for _, c := range con.Campos {
		if c.NombreCampo == campo.NombreCampo {
			apariciones++
			seRepiteCampo = true
		}
		if c.Expresion == campo.Expresion {
			apariciones++
		}
		if c.AliasSql == campo.AliasSql {
			seRepiteAlias = true
		}
	}
	if seRepiteAlias && apariciones > 0 && campo.CampoID != nil {
		campo.AliasSql = fmt.Sprintf("%s_%d", expresion, apariciones)
	}
	if seRepiteCampo && apariciones > 0 {
		campo.NombreCampo = fmt.Sprintf("%s%d", campo.NombreCampo, apariciones)
	}

	// Insertar en base de datos
	err = repo.InsertConsultaCampo(campo)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func ReordenarCampoConsulta(consultaID int, oldPosicion, newPosicion int, repo Repositorio) error {
	op := gko.Op("MoverCampoConsulta").Ctx("consulta_id", consultaID).Ctx("old", oldPosicion).Ctx("new", newPosicion)
	con, err := GetConsulta(consultaID, repo)
	if err != nil {
		return op.Err(err)
	}
	if oldPosicion < 1 || oldPosicion > len(con.Campos) { // Validar que exista el campo referido
		return op.Msgf("La consulta %v no tiene campo en posición %v, tiene %v campos", consultaID, oldPosicion, len(con.Campos))
	}
	if newPosicion < 1 || newPosicion > len(con.Campos) { // Validar que la nueva posición esté dentro del número de hermanos.
		return op.Msgf("La nueva posición %v es inválida para la consulta %v que tiene %v campos", newPosicion, consultaID, len(con.Campos))
	}
	if oldPosicion == newPosicion {
		return nil
	}
	err = repo.ReordenarCampoConsulta(consultaID, oldPosicion, newPosicion)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func ActualizarCampoConsulta(nuevo ddd.ConsultaCampo, repo Repositorio) error {
	op := gko.Op("ActualizarCampoConsulta").Ctx("consulta_id", nuevo.ConsultaID).Ctx("posicion", nuevo.Posicion)
	con, err := GetConsulta(nuevo.ConsultaID, repo)
	if err != nil {
		return op.Err(err)
	}
	if nuevo.Posicion < 1 || nuevo.Posicion > len(con.Campos) { // Validar que exista el campo referido
		return op.Msgf("La consulta %v no tiene campo en posición %v, tiene %v campos", nuevo.ConsultaID, nuevo.Posicion, len(con.Campos))
	}

	old := con.Campos[nuevo.Posicion-1]
	if nuevo.ConsultaID != old.ConsultaID {
		return op.Msg("No se puede cambiar el ID de la consulta dueña del campo")
	}
	if nuevo.Posicion != old.Posicion {
		return op.Msg("No se puede cambiar la posición del campo con este método")
	}

	updated := ddd.ConsultaCampo{
		ConsultaID:  old.ConsultaID,
		Posicion:    old.Posicion,
		CampoID:     old.CampoID,
		Expresion:   old.Expresion,
		AliasSql:    nuevo.AliasSql,
		NombreCampo: nuevo.NombreCampo,
		TipoGo:      nuevo.TipoGo,
		Pk:          nuevo.Pk,
		Filtro:      nuevo.Filtro,
		GroupBy:     nuevo.GroupBy,
		Descripcion: nuevo.Descripcion,
	}

	if nuevo.Expresion != old.Expresion {
		if old.CampoID != nil {
			return op.Msg("No se puede cambiar la expresión de un campo que no es calculado")
		}
		updated.Expresion = nuevo.Expresion
	}

	for _, c := range con.Campos {
		if updated.NombreCampo == c.NombreCampo && updated.Posicion != c.Posicion {
			return op.Msgf("El nombre '%v' ya está en uso por otro campo", updated.NombreCampo)
		}
		if updated.AliasSql != "" && updated.AliasSql == c.AliasSql && updated.Posicion != c.Posicion {
			return op.Msgf("El alias '%v' ya está en uso por otro campo", updated.AliasSql)
		}
		if updated.Expresion == c.Expresion && updated.AliasSql == c.AliasSql && updated.Posicion != c.Posicion {
			return op.Msgf("La expresión '%v' ya está en uso por otro campo con el mismo alias '%v'", updated.Expresion, updated.AliasSql)
		}
	}

	err = repo.UpdateConsultaCampo(updated)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

func EliminarCampoConsulta(consultaID int, posicion int, repo Repositorio) error {
	op := gko.Op("EliminarCampoConsulta").Ctx("consulta_id", consultaID)
	con, err := GetConsulta(consultaID, repo)
	if err != nil {
		return op.Err(err)
	}
	if posicion < 1 || posicion > len(con.Campos) {
		return op.Msgf("La posición %v no es válida para la consulta %v que tiene %v campos", posicion, consultaID, len(con.Campos))
	}
	err = repo.DeleteConsultaCampo(consultaID, posicion)
	if err != nil {
		return op.Err(err)
	}
	return nil
}

// ================================================================ //
// ========== MODIFICACIONES UNITARIAS ============================ //

type CampoConsultaModif struct {
	ConsultaID int
	Posicion   int
	Valor      string
}
type CampoConsultaModifBool struct {
	ConsultaID int
	Posicion   int
	Valor      bool
}

func CampoConsultaModifExpresion(modif CampoConsultaModif, repo Repositorio, txt *textutils.Utils) (*ddd.ConsultaCampo, error) {
	op := gko.Op("CampoConsultaModifExpresion").Ctx("consulta_id", modif.ConsultaID).Ctx("posicion", modif.Posicion)
	con, err := GetConsulta(modif.ConsultaID, repo)
	if err != nil {
		return nil, op.Err(err)
	}
	if modif.Posicion < 1 || modif.Posicion > len(con.Campos) { // Validar que exista el campo referido
		return nil, op.Msgf("La consulta %v no tiene campo en posición %v, tiene %v campos", modif.ConsultaID, modif.Posicion, len(con.Campos))
	}
	campo, err := repo.GetConsultaCampo(modif.ConsultaID, modif.Posicion)
	if err != nil {
		return nil, op.Err(err)
	}
	campoOld := con.Campos[campo.Posicion-1]
	oldExpr := campo.Expresion
	nuevaExpr := gkt.SinEspaciosExtra(modif.Valor)
	if oldExpr == nuevaExpr {
		return campo, nil
	}
	if nuevaExpr == "" {
		return nil, op.Msg("No se puede dejar sin expresión al campo")
	}
	// Si es un campo calculado simplemente actualizar la expresión
	if campo.CampoID == nil {
		campo.Expresion = nuevaExpr
		gko.LogInfof("Expresión modificada para consulta %v campo %v", con.Consulta.NombreItem, campo.Posicion)

	} else { // Sino buscar la nueva tabla de origen para el campo dado.
		oldSplit := strings.Split(oldExpr, ".")
		newSplit := strings.Split(nuevaExpr, ".")
		if !(len(newSplit) == 2 && len(oldSplit) == 2 && newSplit[1] == oldSplit[1]) {
			return nil, op.Msg("Solo se permite cambiar la tabla origen de un campo no calculado")
		}
		origenSearch := newSplit[0]
		origenFound := ""
		nombreColumna := campoOld.OrigenCampo.NombreColumna
		for _, r := range con.Relaciones {
			if (origenSearch == r.JoinAs || origenSearch == r.Join.Tabla.NombreRepo) && origenFound == "" {
				for _, c := range r.Join.Campos {
					if c.NombreColumna == nombreColumna {
						campo.CampoID = &c.CampoID
						origenFound = r.JoinAs
						break
					}
				}
			}
		}
		if origenFound == "" {
			for _, c := range con.From.Campos {
				if c.NombreColumna != nombreColumna {
					continue
				}
				campo.CampoID = &c.CampoID
				origenFound = con.From.Tabla.Abrev
				break
			}
		}
		if origenFound == "" {
			return nil, op.Msgf("El origen '%v' para el campo '%v' no está en las tablas relacionadas", origenSearch, nombreColumna)
		}
		campo.Expresion = origenFound + "." + nombreColumna
		gko.LogInfof("Origen modificado para consulta %v campo %v", con.Consulta.NombreItem, campo.Posicion)
	}
	err = repo.UpdateConsultaCampo(*campo)
	if err != nil {
		return nil, op.Err(err)
	}
	return campo, nil
}

func CampoConsultaModifAlias(modif CampoConsultaModif, repo Repositorio, txt *textutils.Utils) (*ddd.ConsultaCampo, error) {
	op := gko.Op("CampoConsultaModifAlias").Ctx("consulta_id", modif.ConsultaID).Ctx("posicion", modif.Posicion)
	con, err := GetConsulta(modif.ConsultaID, repo)
	if err != nil {
		return nil, op.Err(err)
	}
	if modif.Posicion < 1 || modif.Posicion > len(con.Campos) { // Validar que exista el campo referido
		return nil, op.Msgf("La consulta %v no tiene campo en posición %v, tiene %v campos", modif.ConsultaID, modif.Posicion, len(con.Campos))
	}
	campo, err := repo.GetConsultaCampo(modif.ConsultaID, modif.Posicion)
	if err != nil {
		return nil, op.Err(err)
	}
	oldAlias := campo.AliasSql
	nuevoAlias := gkt.ToLower(modif.Valor)
	nuevoAlias = strings.ReplaceAll(nuevoAlias, " ", "_")
	nuevoAlias = strings.ReplaceAll(nuevoAlias, "-", "_")
	nuevoAlias = strings.ReplaceAll(nuevoAlias, "__", "_")
	nuevoAlias = textutils.QuitarAcentos(nuevoAlias)
	nuevoAlias = txt.RemoveNonAlphanumeric(nuevoAlias)
	if oldAlias == nuevoAlias {
		return campo, nil
	}
	if nuevoAlias != "" {
		for _, c := range con.Campos {
			if nuevoAlias == c.AliasSql && c.AliasSql != "" {
				return nil, op.Msgf("El alias '%v' ya está en uso por el campo %v", nuevoAlias, c.Posicion)
			}
		}
	}
	campo.AliasSql = nuevoAlias
	err = repo.UpdateConsultaCampo(*campo)
	if err != nil {
		return nil, op.Err(err)
	}
	if nuevoAlias == "" {
		gko.LogInfof("Alias eliminado para consulta %v campo %v", con.Consulta.NombreItem, campo.Posicion)
	} else if oldAlias == "" {
		gko.LogInfof("Alias agregado para consulta %v campo %v", con.Consulta.NombreItem, campo.Posicion)
	} else {
		gko.LogInfof("Alias modificado para consulta %v campo %v", con.Consulta.NombreItem, campo.Posicion)
	}
	return campo, nil
}

func CampoConsultaModifNombre(modif CampoConsultaModif, repo Repositorio, txt *textutils.Utils) (*ddd.ConsultaCampo, error) {
	op := gko.Op("CampoConsultaModifNombre").Ctx("consulta_id", modif.ConsultaID).Ctx("posicion", modif.Posicion)
	con, err := GetConsulta(modif.ConsultaID, repo)
	if err != nil {
		return nil, op.Err(err)
	}
	if modif.Posicion < 1 || modif.Posicion > len(con.Campos) { // Validar que exista el campo referido
		return nil, op.Msgf("La consulta %v no tiene campo en posición %v, tiene %v campos", modif.ConsultaID, modif.Posicion, len(con.Campos))
	}
	campo, err := repo.GetConsultaCampo(modif.ConsultaID, modif.Posicion)
	if err != nil {
		return nil, op.Err(err)
	}
	oldNombre := campo.NombreCampo
	newNombre := gkt.SinEspaciosExtra(modif.Valor)
	newNombre = strings.ReplaceAll(newNombre, " ", "")
	newNombre = textutils.PrimeraMayusc(newNombre)
	newNombre = txt.RemoveNonAlphanumeric(newNombre)
	if newNombre == "" {
		return nil, op.E(gko.ErrDatoIndef).Msg("El nombre del campo Go no puede estar vacío")
	}
	if oldNombre == newNombre {
		return campo, nil
	}
	for _, c := range con.Campos {
		if newNombre == c.NombreCampo {
			return nil, op.E(gko.ErrDatoInvalido).Msgf("El nombre '%v' ya está en uso por el campo #%v", newNombre, c.Posicion)
		}
	}
	campo.NombreCampo = newNombre
	err = repo.UpdateConsultaCampo(*campo)
	if err != nil {
		return nil, op.Err(err)
	}
	return campo, nil
}

func CampoConsultaModifTipo(modif CampoConsultaModif, repo Repositorio, txt *textutils.Utils) (*ddd.ConsultaCampo, error) {
	op := gko.Op("CampoConsultaModifTipo").Ctx("consulta_id", modif.ConsultaID).Ctx("posicion", modif.Posicion)
	con, err := GetConsulta(modif.ConsultaID, repo)
	if err != nil {
		return nil, op.Err(err)
	}
	if modif.Posicion < 1 || modif.Posicion > len(con.Campos) { // Validar que exista el campo referido
		return nil, op.Msgf("La consulta %v no tiene campo en posición %v, tiene %v campos", modif.ConsultaID, modif.Posicion, len(con.Campos))
	}
	campo, err := repo.GetConsultaCampo(modif.ConsultaID, modif.Posicion)
	if err != nil {
		return nil, op.Err(err)
	}
	oldTipo := campo.TipoGo
	newTipo := gkt.SinEspaciosExtra(modif.Valor)
	newTipo = strings.ReplaceAll(newTipo, " ", "")
	if newTipo == "" {
		return nil, op.E(gko.ErrDatoIndef).Msg("El tipo de dato Go no puede estar vacío")
	}
	if oldTipo == newTipo {
		return campo, nil
	}
	campo.TipoGo = newTipo
	err = repo.UpdateConsultaCampo(*campo)
	if err != nil {
		return nil, op.Err(err)
	}
	return campo, nil
}

func CampoConsultaModifDesc(modif CampoConsultaModif, repo Repositorio, txt *textutils.Utils) (*ddd.ConsultaCampo, error) {
	op := gko.Op("CampoConsultaModifDesc").Ctx("consulta_id", modif.ConsultaID).Ctx("posicion", modif.Posicion)
	con, err := GetConsulta(modif.ConsultaID, repo)
	if err != nil {
		return nil, op.Err(err)
	}
	if modif.Posicion < 1 || modif.Posicion > len(con.Campos) { // Validar que exista el campo referido
		return nil, op.Msgf("La consulta %v no tiene campo en posición %v, tiene %v campos", modif.ConsultaID, modif.Posicion, len(con.Campos))
	}
	campo, err := repo.GetConsultaCampo(modif.ConsultaID, modif.Posicion)
	if err != nil {
		return nil, op.Err(err)
	}
	oldDesc := campo.Descripcion
	newDesc := gkt.SinEspaciosExtra(modif.Valor)
	if oldDesc == newDesc {
		return campo, nil
	}
	campo.Descripcion = newDesc
	err = repo.UpdateConsultaCampo(*campo)
	if err != nil {
		return nil, op.Err(err)
	}
	return campo, nil
}

func CampoConsultaModifPK(modif CampoConsultaModifBool, repo Repositorio, txt *textutils.Utils) (*ddd.ConsultaCampo, error) {
	op := gko.Op("CampoConsultaModifPk").Ctx("consulta_id", modif.ConsultaID).Ctx("posicion", modif.Posicion)
	con, err := GetConsulta(modif.ConsultaID, repo)
	if err != nil {
		return nil, op.Err(err)
	}
	if modif.Posicion < 1 || modif.Posicion > len(con.Campos) { // Validar que exista el campo referido
		return nil, op.Msgf("La consulta %v no tiene campo en posición %v, tiene %v campos", modif.ConsultaID, modif.Posicion, len(con.Campos))
	}
	campo, err := repo.GetConsultaCampo(modif.ConsultaID, modif.Posicion)
	if err != nil {
		return nil, op.Err(err)
	}
	if campo.Pk == modif.Valor {
		return campo, nil
	}
	campo.Pk = modif.Valor
	err = repo.UpdateConsultaCampo(*campo)
	if err != nil {
		return nil, op.Err(err)
	}
	return campo, nil
}
func CampoConsultaModifFiltro(modif CampoConsultaModifBool, repo Repositorio, txt *textutils.Utils) (*ddd.ConsultaCampo, error) {
	op := gko.Op("CampoConsultaModifPk").Ctx("consulta_id", modif.ConsultaID).Ctx("posicion", modif.Posicion)
	con, err := GetConsulta(modif.ConsultaID, repo)
	if err != nil {
		return nil, op.Err(err)
	}
	if modif.Posicion < 1 || modif.Posicion > len(con.Campos) { // Validar que exista el campo referido
		return nil, op.Msgf("La consulta %v no tiene campo en posición %v, tiene %v campos", modif.ConsultaID, modif.Posicion, len(con.Campos))
	}
	campo, err := repo.GetConsultaCampo(modif.ConsultaID, modif.Posicion)
	if err != nil {
		return nil, op.Err(err)
	}
	if campo.Filtro == modif.Valor {
		return campo, nil
	}
	campo.Filtro = modif.Valor
	err = repo.UpdateConsultaCampo(*campo)
	if err != nil {
		return nil, op.Err(err)
	}
	return campo, nil
}
func CampoConsultaModifGroup(modif CampoConsultaModifBool, repo Repositorio, txt *textutils.Utils) (*ddd.ConsultaCampo, error) {
	op := gko.Op("CampoConsultaModifPk").Ctx("consulta_id", modif.ConsultaID).Ctx("posicion", modif.Posicion)
	con, err := GetConsulta(modif.ConsultaID, repo)
	if err != nil {
		return nil, op.Err(err)
	}
	if modif.Posicion < 1 || modif.Posicion > len(con.Campos) { // Validar que exista el campo referido
		return nil, op.Msgf("La consulta %v no tiene campo en posición %v, tiene %v campos", modif.ConsultaID, modif.Posicion, len(con.Campos))
	}
	campo, err := repo.GetConsultaCampo(modif.ConsultaID, modif.Posicion)
	if err != nil {
		return nil, op.Err(err)
	}
	if campo.GroupBy == modif.Valor {
		return campo, nil
	}
	campo.GroupBy = modif.Valor
	err = repo.UpdateConsultaCampo(*campo)
	if err != nil {
		return nil, op.Err(err)
	}
	return campo, nil
}
