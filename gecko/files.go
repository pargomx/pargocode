package gecko

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

// Envía un archivo como respuesta o el index.html si es un directorio.
func fsFile(c *Context, fpath string, filesystem fs.FS) error {
	// Si se pide el dir raíz fpath vendrá vacío y es inválido para fs.Stat
	if fpath == "" {
		fpath = "."
	}
	// Verificar si el archivo existe.
	fi, err := fs.Stat(filesystem, fpath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			c.LogError("fsFile.Stat('"+fpath+"'): ", err)
		}
		return ErrNotFound
	}
	// Si es un directorio se sirve el index.html
	if fi.IsDir() {
		return fsDirIndex(c, fpath, fi, filesystem)
	}
	// Abrir el archivo.
	file, err := filesystem.Open(fpath)
	if err != nil {
		c.LogError("fsFile.Open('"+fpath+"'): ", err)
		return ErrNotFound
	}
	defer file.Close()
	// Enviar el archivo.
	ff, ok := file.(io.ReadSeeker)
	if !ok {
		return errors.New("file does not implement io.ReadSeeker")
	}
	http.ServeContent(c.Response(), c.Request(), fi.Name(), fi.ModTime(), ff)
	return nil
}

// Se servirá el index.html si existe en el directorio fpath.
func fsDirIndex(c *Context, fpath string, fi fs.FileInfo, filesystem fs.FS) error {
	if !fi.IsDir() {
		return errors.New("fpath is not a directory")
	}
	fpath = filepath.ToSlash(filepath.Join(fpath, "index.html"))
	// Abrir el archivo.
	file, err := filesystem.Open(fpath)
	if err != nil {
		c.LogError("fsFile.Open('"+fpath+"'): ", err)
		return ErrNotFound
	}
	defer file.Close()
	// Enviar el archivo.
	ff, ok := file.(io.ReadSeeker)
	if !ok {
		return errors.New("file does not implement io.ReadSeeker")
	}
	http.ServeContent(c.Response(), c.Request(), fi.Name(), fi.ModTime(), ff)
	return nil
}

// staticDirectoryHandler creates handler function to serve files from provided
// file system When disablePathUnescaping is set then file name from path is not
// unescaped and is served as is.
func staticDirectoryHandler(filesystem fs.FS) HandlerFunc {
	return func(c *Context) error {
		fpath := c.Param("fpath")
		// Convertir %2F a / por ejemplo.
		fpath, err := url.PathUnescape(fpath)
		if err != nil {
			return fmt.Errorf("failed to unescape path variable: %w", err)
		}
		// Necesario en windows porque fs.FS solo usa slashes.
		fpath = filepath.ToSlash(fpath)
		// fs.Open() asume que fpath es relativa al root y rechaza el prefijo `/`.
		filepath.Clean(strings.TrimPrefix(fpath, "/"))
		// Servir el archivo solicitado.
		return fsFile(c, fpath, filesystem)
	}
}

// ================================================================ //
// ========== SERVIR ARCHIVOS ===================================== //

// Crea un handler para servir un archivo estático desde un filesystem dado.
func StaticFileHandler(file string, filesystem fs.FS) HandlerFunc {
	return func(c *Context) error {
		return fsFile(c, file, filesystem)
	}
}

// Registra una ruta para servir un archivo desde el filesystem dado.
func (g *Gecko) FileFS(path, file string, filesystem fs.FS) {
	g.registrarRuta(http.MethodGet, path, func(c *Context) error {
		return fsFile(c, file, filesystem)
	})
}

// Registra una nueva ruta para servir un archivo.
func (g *Gecko) File(path, file string) {
	g.registrarRuta(http.MethodGet, path, func(c *Context) error {
		return fsFile(c, file, c.gecko.Filesystem)
	})
}

// Envía el contenido de un archivo como respuesta desde un filesystem dado.
// Deduce el ContentType y se encarga del caché gracias a http.ServeContent.
func (c *Context) FileFS(file string, filesystem fs.FS) error {
	return fsFile(c, file, filesystem)
}

// Envía el contenido de un archivo como respuesta.
// Deduce el ContentType y se encarga del caché gracias a http.ServeContent.
func (c *Context) File(file string) error {
	return fsFile(c, file, c.gecko.Filesystem)
}

// FileAttachment es similar a File() excepto que se usa para enviar
// un archivo como adjunto, especificando un nombre para él.
//
// Hace que el navegador descargue el archivo sin visualizarlo.
//
// Content-Disposition = "attachment; filename=<FILE_NAME>"
func (c *Context) FileAttachment(file, name string) error {
	c.response.Header().Set(HeaderContentDisposition, fmt.Sprintf("attachment; filename=%q", name))
	return fsFile(c, file, c.gecko.Filesystem)
}

// FileInline es similar a File() excepto que se usa para enviar
// un archivo como inline, especificando un nombre para él.
//
// Hace que el navegador visualice el archivo sin descargarlo.
//
// Content-Disposition = "inline; filename=<FILE_NAME>"
func (c *Context) FileInline(file, name string) error {
	c.response.Header().Set(HeaderContentDisposition, fmt.Sprintf("inline; filename=%q", name))
	return fsFile(c, file, c.gecko.Filesystem)
}

// ================================================================ //
// ========== SERVIR DIRECTORIOS ================================== //

// Crea un filesystem en donde el `fsRoot` sea la nueva raíz.
func mustSubFS(currentFs fs.FS, fsRoot string) fs.FS {
	fsRoot = filepath.ToSlash(filepath.Clean(fsRoot))
	subFs, err := fs.Sub(currentFs, fsRoot)
	if err != nil {
		FatalFmt("imposible crear subFS: invalid root: %v", err)
	}
	return subFs
}

// Se registra tanto /files como /files/*. para que el mux no redireccione
// /files a /files/ en un loop infinito debido a quitarTrailingSlash(r).

// Registra una ruta para servir archivos en el directorio actual
// desde la ruta dada.
func (g *Gecko) Static(pathPrefix string) {
	handler := staticDirectoryHandler(g.Filesystem)
	g.registrarRuta(http.MethodGet, pathPrefix, handler)
	g.registrarRuta(http.MethodGet, path.Join(pathPrefix, "{fpath...}"), handler)
}

// Registra una ruta para servir archivos en el directorio dado.
func (g *Gecko) StaticSub(pathPrefix string, fsRoot string) {
	handler := staticDirectoryHandler(mustSubFS(g.Filesystem, fsRoot))
	g.registrarRuta(http.MethodGet, pathPrefix, handler)
	g.registrarRuta(http.MethodGet, path.Join(pathPrefix, "{fpath...}"), handler)
}

// Registra una ruta para servir archivos estáticos desde un filesystem dado.
// Para `//go:embed static/img` usar `fs := gecko.MustSubFS(fs, "static/img")`.
func (g *Gecko) StaticFS(pathPrefix string, filesystem fs.FS) {
	handler := staticDirectoryHandler(filesystem)
	g.registrarRuta(http.MethodGet, pathPrefix, handler)
	g.registrarRuta(http.MethodGet, path.Join(pathPrefix, "{fpath...}"), handler)
}

// Registra una ruta para servir archivos estáticos desde un filesystem dado
// quitando el prefijo de la ruta.
func (g *Gecko) StaticSubFS(pathPrefix string, fsRoot string, filesystem fs.FS) {
	handler := staticDirectoryHandler(mustSubFS(filesystem, fsRoot))
	g.registrarRuta(http.MethodGet, pathPrefix, handler)
	g.registrarRuta(http.MethodGet, path.Join(pathPrefix, "{fpath...}"), handler)
}

// ================================================================ //
