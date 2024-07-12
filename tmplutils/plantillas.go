package tmplutils

// import (
// 	"bytes"
// 	"go/format"
// 	"io"
// 	"os"
// 	"os/exec"
// 	"strings"
// 	"text/template"

// 	"github.com/pargomx/gecko"
// )

// // ================================================================ //
// // ========== EJECUTAR ============================================ //

// func HaciaBuffer(tmplString string, data any, buf io.Writer) error {
// 	tmpl, err := template.New("").Funcs(funcMap).Parse(tmplString)
// 	if err != nil {
// 		return gecko.NewErr(800).Err(err).Op("parseTemplate")
// 	}

// 	if err = tmpl.Execute(buf, data); err != nil {
// 		return gecko.NewErr(800).Err(err).Op("execTemplate")
// 	}

// 	return nil
// }

// func HaciaString(tmplString string, data any) (string, error) {
// 	tmpl, err := template.New("").Funcs(funcMap).Parse(tmplString)
// 	if err != nil {
// 		return "", gecko.NewErr(800).Err(err).Op("parseTemplate")
// 	}

// 	buf := new(bytes.Buffer)

// 	if err = tmpl.Execute(buf, data); err != nil {
// 		return "", gecko.NewErr(800).Err(err).Op("execTemplate")
// 	}

// 	return buf.String(), nil
// }

// func HaciaArchivo(tmplString string, data any, filename string) error {
// 	tmpl, err := template.New("").Funcs(funcMap).Parse(tmplString)
// 	if err != nil {
// 		return gecko.NewErr(800).Err(err).Op("parseTemplate")
// 	}

// 	fileOut, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}

// 	if err = tmpl.Execute(fileOut, data); err != nil {
// 		return gecko.NewErr(800).Err(err).Op("execTemplate")
// 	}

// 	return fileOut.Close()
// }

// // ================================================================ //
// // ========== FORMAT GO =========================================== //

// // Da formato al código fuente de Go y lo escribe en buf.
// // Si encuentra error en el formato escribe el resultado tal cual y agrega un comentario.
// func HaciaBufferGo(tmplString string, data any, buf io.Writer) error {
// 	tmpl, err := template.New("").Funcs(funcMap).Parse(tmplString)
// 	if err != nil {
// 		return gecko.NewErr(800).Err(err).Op("parseTemplate")
// 	}

// 	var buf2 bytes.Buffer

// 	if err = tmpl.Execute(&buf2, data); err != nil {
// 		return gecko.NewErr(800).Err(err).Op("execTemplate")
// 	}

// 	codigo, err := format.Source(buf2.Bytes())
// 	if err != nil { // no retornar error para facilitar debug del código generado
// 		codigo = []byte("// ERROR: " + gecko.NewErr(800).Err(err).Op("fmt_gocode").Error() + "\n\n" + buf2.String())
// 	}
// 	_, err = buf.Write(codigo)
// 	if err != nil {
// 		return gecko.NewErr(800).Err(err).Op("write_buffer")
// 	}
// 	return nil
// }

// // Escribe en el archivo dado y ejecuta goimports en él.
// func HaciaArchivoGo(tmplString string, data any, filename string) error {
// 	tmpl, err := template.New("").Funcs(funcMap).Parse(tmplString)
// 	if err != nil {
// 		return gecko.NewErr(800).Err(err).Op("parseTemplate")
// 	}

// 	fileOut, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}

// 	if err = tmpl.Execute(fileOut, data); err != nil {
// 		return gecko.NewErr(800).Err(err).Op("execTemplate")
// 	}
// 	fileOut.Close()

// 	cmd := exec.Command("goimports", "-w", filename)
// 	if errOut, err := cmd.CombinedOutput(); err != nil {
// 		return gecko.Err(err).Op("goimports").Msg("no se pudo dar formato").Msg(string(errOut))
// 	}
// 	return nil
// }

// // ================================================================ //
// // ========== FORMAT HTML ========================================= //

// // Transforma los '[[ ]]' en '{{ }}'.
// func HaciaBufferHTML(tmplString string, data any, buf io.Writer) (err error) {
// 	tmpl, err := template.New("").Funcs(funcMap).Parse(tmplString)
// 	if err != nil {
// 		return gecko.NewErr(800).Err(err).Op("parseTemplate")
// 	}

// 	buf2 := new(bytes.Buffer)

// 	if err = tmpl.Execute(buf2, data); err != nil {
// 		return gecko.NewErr(800).Err(err).Op("execTemplate")
// 	}

// 	strings.NewReplacer("[[", "{{", "]]", "}}").WriteString(buf, buf2.String())
// 	return nil
// }

// // Transforma los '[[ ]]' en '{{ }}'.
// func HaciaStringHTML(tmplString string, data any) (html string, err error) {
// 	tmpl, err := template.New("").Funcs(funcMap).Parse(tmplString)
// 	if err != nil {
// 		return "", gecko.NewErr(800).Err(err).Op("parseTemplate")
// 	}

// 	buf := new(bytes.Buffer)

// 	if err = tmpl.Execute(buf, data); err != nil {
// 		return "", gecko.NewErr(800).Err(err).Op("execTemplate")
// 	}

// 	return strings.NewReplacer("[[", "{{", "]]", "}}").Replace(buf.String()), nil
// }

// // Transforma los '[[ ]]' en '{{ }}' y escribe en el archivo dado.
// func HaciaArchivoHTML(tmplString string, data any, filename string) error {
// 	tmpl, err := template.New("").Funcs(funcMap).Parse(tmplString)
// 	if err != nil {
// 		return gecko.NewErr(800).Err(err).Op("parseTemplate")
// 	}

// 	buf := new(bytes.Buffer)

// 	if err = tmpl.Execute(buf, data); err != nil {
// 		return gecko.NewErr(800).Err(err).Op("execTemplate")
// 	}

// 	fileOut, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	strings.NewReplacer("[[", "{{", "]]", "}}").WriteString(fileOut, buf.String())

// 	return fileOut.Close()
// }
