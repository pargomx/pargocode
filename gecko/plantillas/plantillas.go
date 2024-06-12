package plantillas

import (
	"html/template"
	"os"

	"monorepo/gecko"
)

// TemplateResponder prepara plantillas html y las ejecuta con
// gracia a discreción.
type TemplateResponder struct {

	// t representa a una o varias plantillas identificadas por
	// un nombre único. Las plantillas pueden o no estar
	// asociadas entre sí (anidadas unas en otras).
	t *template.Template

	// Carpeta desde donde obtener las plantillas.
	carpeta string

	// IDEA: Plantillas desde archivos adicionales fuera de la carpeta.
	// archivos []string

	// Volver a leer plantillas antes de ejecutarlas.
	// Utilizar solamente durante el desarrollo.
	reparse bool
}

// NuevoServicioPlantillas prepara todas las plantillas dentro
// del directorio especificado (incluyendo subdirectorios) y
// las mantiene en memoria para ser ejecutadas a discreción.
//
// Si no se pasan archivos, se usan todos los .html de la
// carpeta.
//
// Si se pasan archivos, se usan solo los especificados.
//
// Cuando haya error no debe usarse la plantilla si no se
// quiere panic por nil pointer.
//
// Si el env LISTAR_PLANTILLAS=true se imprime las plantillas
// preparadas para debug.
//
// Si el entorno tiene REPARSE_PLANTILLAS=true o
// AMBIENTE=desarrollo entonces las plantillas se volverán a
// leer cada vez que se vaya a ejecutar una, para ver
// reflejados los cambios hechos durante el desarrollo al
// recargar la página.
//
// El nombre de cada plantilla incluye subcarpetas dentro de
// tplsDir, y no incluye la extensión. Por ejemplo:
//
// "entidad/nuevo" para usar
// "/xxx/plantillas/entidad/nuevo.html"
func NuevoServicioPlantillas(carpeta string) (*TemplateResponder, error) {
	op := gecko.NewOp("plantillas.NuevoServicioPlantillas")

	tmpl, err := findAndParseTemplates(carpeta, funcMap)
	if err != nil {
		return nil, op.Err(err)
	}

	s := &TemplateResponder{
		t:       tmpl,
		carpeta: carpeta,
	}

	if os.Getenv("LISTAR_PLANTILLAS") == "true" {
		gecko.LogOkeyf("Plantillas preparadas desde %v", s.carpeta)
		s.Listar()
	}

	if os.Getenv("AMBIENTE") == "desarrollo" ||
		os.Getenv("REPARSE_PLANTILLAS") == "true" {
		s.reparse = true
	}
	if os.Getenv("REPARSE_PLANTILLAS") == "false" {
		s.reparse = false
	}

	return s, nil
}

// ================================================================ //

// Lookup returns the template with the given name that is
// associated with t, or nil if there is no such template.
func (s *TemplateResponder) Lookup(nombre string) *template.Template {
	return s.t.Lookup(nombre)
}
