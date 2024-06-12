package gecko

import (
	"net/http"
	"strings"
)

// Tratar rutas con trailing slash como si no lo tuvieran.
// Utilizado como middleware global antes del router.
func quitarTrailingSlash(r *http.Request) {
	url := r.URL
	path := url.Path
	queryString := r.URL.RawQuery
	l := len(path) - 1
	if l > 0 && strings.HasSuffix(path, "/") {
		path = path[:l]
		uri := path
		if queryString != "" {
			uri += "?" + queryString
		}
		r.RequestURI = uri
		url.Path = path
	}
}

// ================================================================ //
// ========== RUTAS Y HANDLERS ==================================== //

// HandlerFunc defines a function to serve HTTP requests.
type HandlerFunc func(c *Context) error

// Registrar una nueva ruta con un http.HandlerFunc
// que prepare el gecko.Context y ejecute el gecko.HandlerFunc.
func (g *Gecko) registrarRuta(método string, ruta string, handler HandlerFunc) {
	patrón := toMuxPattern(método, ruta)
	g.mux.HandleFunc(patrón, func(w http.ResponseWriter, r *http.Request) {
		// Preparar contexto.
		c := &Context{
			request:  r,
			response: NewResponse(w, g),
			path:     patrón,
			gecko:    g,
		}
		// Ejecutar handler.
		err := handler(c)
		if err != nil {
			// c.LogError(err) // TODO: handle errors.
			// g.HTTPErrorHandler(err, c)
			g.GeckoHTTPErrorHandler(err, c)
		}
	})
	// fmt.Println("RUTA:", patrón)
}

// Necesario para validar patrón de ruta con método y las reglas de gecko.
func toMuxPattern(método string, ruta string) string {
	// Validar la ruta.
	if strings.Contains(ruta, " ") {
		FatalFmt("La ruta no puede contener espacios en blanco: '%s'", ruta)
	}
	ruta = strings.TrimSuffix(ruta, "/")
	if ruta == "" {
		ruta = "/{$}"
	}
	if ruta[0] != '/' {
		FatalFmt("La ruta debe comenzar con slash: '%s'", ruta)
	}
	// Validar método.
	if método == "" {
		FatalFmt("El método no puede estar indefinido: '%s'", ruta)
	}
	return método + " " + ruta
}

// ================================================================ //
// ========== Registrar handlers con métodos ====================== //

func (g *Gecko) GET(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodGet, path, handler)
}
func (g *Gecko) POST(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodPost, path, handler)
}
func (g *Gecko) PUT(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodPut, path, handler)
}
func (g *Gecko) PATCH(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodPatch, path, handler)
}
func (g *Gecko) DELETE(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodDelete, path, handler)
}

func (g *Gecko) POS(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodPost, path, handler)
}
func (g *Gecko) PCH(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodPatch, path, handler)
}
func (g *Gecko) DEL(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodDelete, path, handler)
}

/*
func (g *Gecko) OPTIONS(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodOptions, path, handler)
}
func (g *Gecko) HEAD(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodHead, path, handler)
}
func (g *Gecko) CONNECT(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodConnect, path, handler)
}
func (g *Gecko) TRACE(path string, handler HandlerFunc) {
	g.registrarRuta(http.MethodTrace, path, handler)
}
*/

// ================================================================ //
