package plantillas

import (
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"monorepo/gecko"
)

// Camina por el directorio especificado y prepara todas las
// plantillas html que encuentra con el funcMap dado.
//
// Cualquier plantilla con errores sintácticos provoca que
// ninguna esté disponible. Es útil saberlo desde el inicio y
// no hasta la hora de ejecutarla.
//
// El nombre de cada plantilla incluye subcarpetas dentro del
// directorio raíz, y no incluye la extensión html. Ejemplo:
//
// Ejecutar "entidad/nuevo" utilizará el archivo
// "/xxx/plantillas/entidad/nuevo.html"
func findAndParseTemplates(tplsDir string, funcMap template.FuncMap) (*template.Template, error) {
	op := gecko.NewOp("plantillas.findAndParseTemplates")

	// Validar ruta a carpeta de plantillas.
	tplsDir = filepath.Clean(tplsDir)
	f, err := os.Stat(tplsDir)
	if err != nil {
		return nil, op.Msg("directorio inaccesible: " + tplsDir)
	}
	if !f.IsDir() {
		return nil, op.Msg("no es un directorio: " + tplsDir)
	}

	// Plantilla de la cual colgarán todas.
	rootTmpl := template.New("")

	// TODO: Cambiar a WalkDir
	// Escanear todos los archivos y subcarpetas.
	err = filepath.Walk(tplsDir, func(path string, info os.FileInfo, errWalk error) error {
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

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Colgar nueva plantilla.
		t := rootTmpl.New(nombre).Funcs(funcMap)
		_, err = t.Parse(string(b))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, op.Err(err)
	}

	return rootTmpl, nil
}

// ================================================================ //

// ReParse vuelve a leer todas las plantillas en carpeta y si
// no hay error, las reemplaza en memoria para que el servicio
// utilice las actualizadas.
func (s *TemplateResponder) ReParse() {
	newTmpl, err := findAndParseTemplates(s.carpeta, funcMap)
	if err != nil {
		gecko.LogWarnf("plantillas.ReParse: usando plantillas anteriores: %v", err)
		return
	}
	s.t = newTmpl
}

// ================================================================ //

// Listar imprime en consola todas las plantillas
// que están preparadas en el servicio actualmente.
// Usar como herramienta de debug.
func (s *TemplateResponder) Listar() {
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
