{{ if .PackageDef }}package {{ .Tabla.Paquete.Nombre }}{{br}}{{ end }}

var (
	Err{{ $.Tabla.NombreItem }}NotFound      error = errors.New("{{ .Tabla.NombreNominativo }} no se encuentra")
	Err{{ $.Tabla.NombreItem }}AlreadyExists error = errors.New("{{ .Tabla.NombreNominativo }} ya existe")
)

func ({{ .Tabla.NombreAbrev }} *{{ .Tabla.NombreItem }}) Validar() error {

	{{ range .Tabla.CamposEspeciales }}
	if {{ $.Tabla.NombreAbrev }}.{{ .NombreCampo }}.EsTodos() {
		return errors.New("{{ $.Tabla.Paquete.Nombre }}.{{ $.Tabla.NombreItem }} no admite propiedad {{ .NombreCampo }}Todos")
	}
	{{end}}

	return nil
}