package tmplutils

import (
	"io/fs"
	"sort"
	"strings"
	"text/template"
)

type Renderer struct {
	tmpls *template.Template
}

func NuevoRenderer(fs fs.FS, tmplsDir string) (*Renderer, error) {
	tmpls, err := PrepararPlantillasFromDir(fs, tmplsDir)
	if err != nil {
		return nil, err
	}
	return &Renderer{
		tmpls: tmpls,
	}, nil
}

func (s *Renderer) ListarPlantillas() []string {
	var nombres []string
	for _, t := range s.tmpls.Templates() {
		nombres = append(nombres, t.Name())
	}
	sort.Strings(nombres)
	return nombres
}

// Prepara todas las plantillas encontradas en el filesystem dado.
//
// Diseñado para usarse con "_embed". El tplsDir es para quitar el
// static/plantillas si por ejemplo se hace embed:./static/ y las plantillas están
// dentro de un subdirectorio "plantillas".
//
// Las plantillas son accesibles sin el ".html", ".gtmpl", ".gtpl"
func PrepararPlantillasFromDir(fsys fs.FS, tplsDir string) (*template.Template, error) {

	// Plantilla de la cual colgarán todas.
	rootTmpl := template.New("root")

	// Escanear todos los archivos y subcarpetas.
	err := fs.WalkDir(fsys, ".", func(path string, info fs.DirEntry, errWalk error) error {
		if errWalk != nil {
			return errWalk
		}

		// Solo nos interesan archivos con estas extensiones.
		if info.IsDir() || !(strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".gtmpl") || strings.HasSuffix(path, ".gtpl")) {
			return nil
		}

		nombre := strings.TrimPrefix(path, tplsDir)
		nombre = strings.TrimPrefix(nombre, "/")     // ej. "/tpls/usu/hola.html" > "usu/hola.html"
		nombre = strings.TrimSuffix(nombre, ".html") // ej. "usuario/nuevo.html" > "usuario/nuevo"
		nombre = strings.TrimSuffix(nombre, ".gtmpl")
		nombre = strings.TrimSuffix(nombre, ".gtpl")

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
		return nil, err
	}

	return rootTmpl, nil
}
