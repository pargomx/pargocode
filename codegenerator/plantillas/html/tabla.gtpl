<header class="flex flex-wrap items-center gap-3 px-4 py-3">
	<h2 class="flex-grow w-full text-xl text-center sm:w-auto sm:text-left">{{ $.Tabla.Nombre.HumanoPlural }}</h2>
	<input class="flex-grow px-3 py-1 rounded sm:flex-grow-0 text-slate-900" type="search" oninput="filtrarTabla(this.value)" placeholder="Buscar en tabla...">
	<button type="button" class="w-10 py-1 text-2xl rounded-lg text-slate-200 bg-slate-600" hx-get="/" title="Recargar"><i class="fa-solid fa-rotate-right"></i></button>
	<button type="button" class="w-10 py-1 text-2xl rounded-lg text-slate-200 bg-slate-600" hx-get="/" title="Cerrar"><i class="fa-solid fa-xmark"></i></button>
	<button type="button" class="w-10 py-1 text-2xl bg-teal-700 rounded-lg text-slate-200" hx-get="/" title="Agregar"><i class="fa-solid fa-plus"></i></button>
	{{ if $.Tabla.CamposFiltro }}
	<form class="inline-flex" id="filtros">
		{{- range $.Tabla.CamposFiltro }}{{ if .Especial }}
		@@include('_select/filtro.html',{
			"name":"{{.NombreColumna}}",
			"nombre":"{{.Nombre.Humano}}",
			"param":".Filtro{{.Nombre.Camel}}",
			"lista":".Lista{{.Nombre.Camel}}",
			"size":"small",
		})
		{{- end }}{{ end }}
	</form>
	<button type="button" class="w-10 py-1 text-2xl rounded-lg text-slate-200 bg-slate-600" hx-get="/" hx-include="#filtros" title="Recargar"><i class="fa-solid fa-rotate-right"></i></button>
	{{ end }}
</header>

[[ if .{{$.Tabla.Nombre.CamelPlural}} ]]
<table class="w-full px-1 text-center border-separate table-auto border-spacing-x-0 border-spacing-y-1">
	<thead>
		<tr class="bg-slate-400 dark:bg-slate-900">
			{{ range $.Tabla.Campos }}
			<th>{{.Nombre.Humano}}</th>{{ end }}
		</tr>
	</thead>
	<tbody>
		[[ range .{{$.Tabla.Nombre.CamelPlural}} ]]
		<tr class="bg-slate-400/25 hover:bg-slate-400/50 dark:bg-slate-800/50 dark:hover:bg-slate-800">
			{{- range $.Tabla.Campos }}
		{{ if .Especial }}
			<td>[[ .{{.Nombre.Camel}}.String ]]</td>
		{{ else if eq .Tipo "time.Time" }}
			<td>[[ .{{.Nombre.Camel}}.Format "2006-01-02" ]]</td>
		{{ else }}
			<td>[[ .{{.Nombre.Camel}} ]]</td>
		{{ end }}
			{{- end }}
		</tr>
		[[ end ]]
	</tbody>
</table>

[[ else ]]
<h4 class="p-2 mx-1 text-lg text-center bg-slate-400 dark:bg-slate-900">Nada que mostrar</h4>
<p class="p-4 text-center">Intente cambiar los filtros</p>
[[ end ]]