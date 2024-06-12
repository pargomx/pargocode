package plantillas

import (
	"html/template"
	"io"
	"io/fs"
	"sort"
	"strings"

	"monorepo/gecko"
)

// TemplateResponder prepara plantillas html y las ejecuta con
// gracia a discreción.
type TemplateResponderFS struct {

	// t representa a una o varias plantillas identificadas por
	// un nombre único. Las plantillas pueden o no estar
	// asociadas entre sí (anidadas unas en otras).
	t *template.Template
}

// Prepara todas las plantillas encontradas en el filesystem dado.
//
// Diseñado para usarse con "_embed". El tplsDir es para quitar el
// static/plantillas si por ejemplo se hace embed:./static/ y las plantillas están
// dentro de un subdirectorio "plantillas".
//
// Las plantillas son accesibles sin el ".html"
func NuevoServicioPlantillasEmbebidas(fsys fs.FS, tplsDir string) (*TemplateResponderFS, error) {
	op := gecko.NewOp("plantillas.NuevoServicioPlantillasEmbebidas")

	// Plantilla de la cual colgarán todas.
	rootTmpl := template.New("")

	// Escanear todos los archivos y subcarpetas.
	err := fs.WalkDir(fsys, ".", func(path string, info fs.DirEntry, errWalk error) error {
		if errWalk != nil {
			return errWalk
		}

		// Solo nos interesan archivos .html
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		nombre := strings.TrimPrefix(path, tplsDir)
		nombre = strings.TrimPrefix(nombre, "/")     // ej. "/tpls/usu/hola.html" > "usu/hola.html"
		nombre = strings.TrimSuffix(nombre, ".html") // ej. "usuario/nuevo.html" > "usuario/nuevo"

		// Se ignoran plantillas que comienzan por "_"
		if strings.HasPrefix(nombre, "_") {
			return nil
		}

		bytes, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}

		// Colgar nueva plantilla.
		t := rootTmpl.New(nombre).Funcs(funcMap)
		_, err = t.Parse(string(bytes))
		// _, err := t.ParseFS(fsys, path)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, op.Err(err)
	}

	s := &TemplateResponderFS{
		t: rootTmpl,
	}

	return s, nil
}

// ================================================================ //
// ========== RENDER ============================================== //

// Render satisface la interfaz gecko.Renderer
//
// Es lo último que se debe llamar en un handler.
//
// Ejecuta una plantilla previamente instanciada al crear el servicio.
//
// Si la plantilla no existe, responde con el error definido en NuevoServicio.
func (s *TemplateResponderFS) Render(w io.Writer, nombre string, data interface{}, c *gecko.Context) error {
	if strings.HasSuffix(nombre, ".html") {
		c.LogWarnf("plantilla.Render: no es necesario poner .html a '%v'", nombre)
		nombre = strings.TrimSuffix(nombre, ".html")
	}
	return s.t.ExecuteTemplate(w, nombre, data)
}

// ================================================================ //

// Listar imprime en consola todas las plantillas
// que están preparadas en el servicio actualmente.
// Usar como herramienta de debug.
func (s *TemplateResponderFS) Listar() {
	var nombres []string
	for _, t := range s.t.Templates() {
		nombres = append(nombres, t.Name())
	}
	sort.Strings(nombres)

	tms := make(map[int]string, len(nombres))
	for i, n := range nombres {
		tms[i] = n
	}
	gecko.LogInfof("plantillas.Listar: %v", tms)
}
