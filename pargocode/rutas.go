package main

import (
	"github.com/pargomx/pargocode/assets"

	"github.com/pargomx/gecko"
)

// ================================================================ //
// ========== RUTAS =============================================== //

func (s *servidor) registrarRutas() {

	s.gecko.StaticFS("/assets", assets.AssetsFS)
	s.gecko.FileFS("/favicon.ico", "img/favicon.ico", assets.AssetsFS)

	// ================================================================ //
	// ================================================================ //

	s.gecko.GET("/", s.getPaquetes)

	s.gecko.GET("/mapa", s.getMapaEntidadRelacion)
	s.gecko.GET("/paquetes", s.getPaquetes)
	s.gecko.POS("/paquetes", s.agregarPaquete)
	s.gecko.PUT("/paquetes/{paquete_id}", s.actualizarPaquete)
	s.gecko.DEL("/paquetes/{paquete_id}", s.eliminarPaquete)
	s.gecko.GET("/paquetes/{paquete_id}", s.getMapaEntidadRelacionPaquete)
	s.gecko.GET("/paquetes/{paquete_id}/generar", s.generarDePaqueteArchivos)

	s.gecko.GET("/tablas", s.getPaquetes)              // 1. Tablas en el proyecto
	s.gecko.GET("/tablas/nueva", s.getTablaNueva)      // 2. Formulario para nueva tabla
	s.gecko.POS("/tablas/nueva", s.postTablaNueva)     // 3. Crear nueva tabla
	s.gecko.GET("/tablas/{tabla_id}", s.getTabla)      // 4. Dashboard para tabla
	s.gecko.PUT("/tablas/{tabla_id}", s.putTabla)      // 5. Actualizar datos de la tabla
	s.gecko.DEL("/tablas/{tabla_id}", s.eliminarTabla) // 6. Eliminar tabla
	s.gecko.POS("/tablas/{tabla_id}/campos", s.postCampo)
	s.gecko.PUT("/tablas/{tabla_id}/campos", s.putCampo)
	s.gecko.GET("/tablas/{tabla_id}/generar", s.generarDeTabla)
	s.gecko.POS("/tablas/{tabla_id}/campos_ordenar", s.fixOrdenDeCampos)

	s.gecko.PUT("/campos/{campo_id}", s.updateCampo)
	s.gecko.DEL("/campos/{campo_id}", s.deleteCampo)
	s.gecko.GET("/campos/{campo_id}/form", s.formCampo)
	s.gecko.PUT("/campos/{campo_id}/reordenar", s.reordenarCampo)
	s.gecko.GET("/campos/{campo_id}/enum", s.getEnumCampo)
	s.gecko.POS("/campos/{campo_id}/enum", s.postEnumCampo)

	s.gecko.POS("/consultas", s.crearConsulta)
	s.gecko.GET("/consultas/nueva", s.formNuevaConsulta)
	s.gecko.GET("/consultas/{consulta_id}", s.getConsulta)
	s.gecko.DEL("/consultas/{consulta_id}", s.deleteConsulta)
	s.gecko.PUT("/consultas/{consulta_id}", s.actualizarConsulta)
	s.gecko.GET("/consultas/{consulta_id}/generar", s.generarDeConsulta)
	s.gecko.POS("/consultas/{consulta_id}/relaciones", s.postRelacionConsulta)
	s.gecko.PUT("/consultas/{consulta_id}/relaciones/{posicion}", s.actualizarRelacionConsulta)
	s.gecko.DEL("/consultas/{consulta_id}/relaciones/{posicion}", s.eliminarRelacionConsulta)

	s.gecko.POS("/consultas/{consulta_id}/campos", s.postCampoConsulta)
	s.gecko.PUT("/consultas/{consulta_id}/reordenar-campo", s.reordenarCampoConsulta)
	s.gecko.DEL("/consultas/{consulta_id}/campos/{posicion}", s.eliminarCampoConsulta)
	s.gecko.PUT("/consultas/{consulta_id}/campos/{posicion}", s.actualizarCampoConsulta)
	// atómico
	s.gecko.PUT("/consultas/{consulta_id}/campos/{posicion}/expresion", s.campoConsModifExpresion)
	s.gecko.PUT("/consultas/{consulta_id}/campos/{posicion}/alias_sql", s.campoConsModifAlias)
	s.gecko.PUT("/consultas/{consulta_id}/campos/{posicion}/nombre_campo", s.campoConsModifNombre)
	s.gecko.PUT("/consultas/{consulta_id}/campos/{posicion}/tipo_go", s.campoConsModifTipo)
	s.gecko.PUT("/consultas/{consulta_id}/campos/{posicion}/pk", s.campoConsModifPK)
	s.gecko.PUT("/consultas/{consulta_id}/campos/{posicion}/filtro", s.campoConsModifFiltro)
	s.gecko.PUT("/consultas/{consulta_id}/campos/{posicion}/group_by", s.campoConsModifGroup)
	s.gecko.PUT("/consultas/{consulta_id}/campos/{posicion}/descripcion", s.campoConsModifDesc)

	// LOG SQLITE
	s.gecko.GET("/log", func(c *gecko.Context) error { s.db.ToggleLog(); return c.StatusOk("Log toggled") })

}
