package gecko

import (
	"encoding/json"
	"net/http"
)

// Recibe un JSON del request y pone los datos en v con json.Unmarshal
func (c *Context) JSONUnmarshal(v any) error {
	err := json.NewDecoder(c.Request().Body).Decode(v)
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		return NewErr(http.StatusBadRequest).Err(err).Msgf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)
	} else if se, ok := err.(*json.SyntaxError); ok {
		return NewErr(http.StatusBadRequest).Err(err).Msgf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())
	}
	return err
}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func jsonSerialize(c *Context, i interface{}, indent string) error {
	enc := json.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func jsonDeserialize(c *Context, i interface{}) error {
	err := json.NewDecoder(c.Request().Body).Decode(i)
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		return NewErr(http.StatusBadRequest).Msgf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset).Err(err)
	} else if se, ok := err.(*json.SyntaxError); ok {
		return NewErr(http.StatusBadRequest).Msgf("Syntax error: offset=%v, error=%v", se.Offset, se.Error()).Err(err)
	}
	return err
}

// ================================================================ //
// ========== ENVIAR JSON ========================================= //

// Responder con MIME "application/json" UTF8.
//
// Si hay un QueryParam "pretty" se indenta la respuesta.
//
// Utiliza json.Marshal y puede no ser eficiente para grandes objetos.
func (c *Context) JSON(code int, i interface{}) (err error) {
	indent := ""
	if _, pretty := c.QueryParams()["pretty"]; pretty {
		indent = "  "
	}
	return c.json(code, i, indent)
}

// Responder con MIME "application/json" UTF8.
//
// Se imprmie con la indentación dada (espacios o tabs).
//
// Utiliza json.Marshal y puede no ser eficiente para grandes objetos.
func (c *Context) JSONPretty(code int, i interface{}, indent string) (err error) {
	return c.json(code, i, indent)
}

// Responder con MIME "application/json" UTF8.
//
// Útil para JSON que ya fue codificado en otro lugar.
func (c *Context) JSONBlob(code int, b []byte) (err error) {
	return c.Blob(code, MIMEApplicationJSONCharsetUTF8, b)
}

// JSON Pading encapsula el JSON en la función callback dada.
//
// Responder con MIME "application/json" UTF8.
//
// Si hay un QueryParam "pretty" se indenta la respuesta.
//
// Utiliza json.Marshal y puede no ser eficiente para grandes objetos.
//
// Ejemplo: miFuncion({x=1,y=2});
func (c *Context) JSONP(code int, callback string, i interface{}) (err error) {
	return c.jsonPBlob(code, callback, i)
}

// JSON Pading encapsula el JSON en la función callback dada.
//
// Responder con MIME "application/json" UTF8.
//
// Útil para JSON que ya fue codificado en otro lugar.
//
// Ejemplo: miFuncion({x=1,y=2});
func (c *Context) JSONPBlob(code int, callback string, b []byte) (err error) {
	c.writeContentType(MIMEApplicationJavaScriptCharsetUTF8)
	c.response.WriteHeader(code)
	if _, err = c.response.Write([]byte(callback + "(")); err != nil {
		return
	}
	if _, err = c.response.Write(b); err != nil {
		return
	}
	_, err = c.response.Write([]byte(");"))
	return
}

// ================================================================ //

func (c *Context) jsonPBlob(code int, callback string, i interface{}) (err error) {
	indent := ""
	if _, pretty := c.QueryParams()["pretty"]; pretty {
		indent = "  "
	}
	c.writeContentType(MIMEApplicationJavaScriptCharsetUTF8)
	c.response.WriteHeader(code)
	if _, err = c.response.Write([]byte(callback + "(")); err != nil {
		return
	}
	if err = jsonSerialize(c, i, indent); err != nil {
		return
	}
	if _, err = c.response.Write([]byte(");")); err != nil {
		return
	}
	return
}

func (c *Context) json(code int, i interface{}, indent string) error {
	c.writeContentType(MIMEApplicationJSONCharsetUTF8)
	c.response.Status = code
	return jsonSerialize(c, i, indent)
}
