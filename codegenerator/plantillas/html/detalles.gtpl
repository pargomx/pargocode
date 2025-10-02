[[ with .{{ .Tabla.Nombre.Camel }} ]]
<header class="flex flex-wrap items-center gap-3 px-4 py-3">
	<h2 class="grow w-full text-xl text-center sm:w-auto sm:text-left">Detalles</h2>
	<button type="button" class="w-10 py-1 text-2xl rounded-lg text-slate-200 bg-slate-600" hx-get="" title="Recargar"><i class="fa-solid fa-rotate-right"></i></button>
	<button type="button" class="w-10 py-1 text-2xl rounded-lg text-slate-200 bg-slate-600" hx-get="" title="Cerrar"><i class="fa-solid fa-xmark"></i></button>
</header>

<table class="w-full px-1 text-center border-separate table-auto border-spacing-x-0 border-spacing-y-1">
	<tr>
		<td colspan="2">General</td>
	</tr>
	{{ range .Tabla.Campos -}}
	<tr>
		<td>{{.Nombre.Humano}}</td>
		{{ if .Especial -}}
		<td>[[ .{{.Nombre.Camel}}.String ]]: [[ .{{.Nombre.Camel}}.Descripcion ]]</td>
		{{- else if eq .Tipo "time.Time" -}}
		<td>[[ .{{.Nombre.Camel}}.String ]]</td>
		{{- else -}}
		<td>[[ .{{.Nombre.Camel}} ]]</td>
		{{- end }}
	</tr>
	{{ end }}
</table>
[[ end ]]