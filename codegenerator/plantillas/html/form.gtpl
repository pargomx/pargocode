[[ with .{{ .Tabla.Nombre.Camel }} ]]
<header class="flex flex-wrap items-center gap-3 px-4 py-3">
	<h2 class="grow w-full text-xl text-center sm:w-auto sm:text-left">Form</h2>
	<button type="button" class="w-10 py-1 text-2xl rounded-lg text-slate-200 bg-slate-600" hx-get="" title="Recargar"><i class="fa-solid fa-rotate-right"></i></button>
	<button type="button" class="w-10 py-1 text-2xl rounded-lg text-slate-200 bg-slate-600" hx-get="" title="Cerrar"><i class="fa-solid fa-xmark"></i></button>
</header>

<form>
	{{ range .Tabla.CamposEditables }}
	{{ if .Especial -}}
	<div class="form-group">
		[[- $actual := .{{.Nombre.Camel}} ]]
		<label>{{.Nombre.Humano}}: <code>[[$actual.String]]</code></label>
		<select name="{{.NombreColumna}}" type="text" class="form-control">
			[[range $.Lista{{.Nombre.Camel}} ]]
			<option value="[[.String]]"[[if eq .String $actual.String ]] selected[[end]]>[[if .String]][[.String]]: [[.Descripcion]][[end]]</option>
			[[end]]
		</select>
	</div>
	{{ else if gt .MaxLenght 900 }}
	<div class="form-group">
		<label>{{ .Nombre.Humano }}</label>
		<textarea name="{{.NombreColumna}}"{{ if .MaxLenght }} maxlength="{{ .MaxLenght }}{{end}} class="form-control textarea-autosize">
			[[- .{{.Nombre.Camel}} -]]
		</textarea>
	</div>
	{{ else -}}
	<div class="form-group">
		<label>{{.Nombre.Humano}}</label>
		<input name="{{.NombreColumna}}" type="text" {{if .MaxLenght}}maxlength="{{.MaxLenght}}" {{ end -}}
			value="[[.{{.Nombre.Camel}}]]" placeholder="" class="form-control">
	</div>
	{{ end -}}
	{{ end }}
	
	<!-- {{ range .Tabla.PrimaryKeys }}/xxx/[[ .{{.Nombre.Camel}} ]]{{ end }} -->
	<div class="flex flex-wrap justify-center gap-4 pt-4">
		<button type="button" hx-get="" class="inline-block px-6 py-2 font-medium transition-colors bg-gray-600 rounded shadow-md text-slate-100 hover:bg-gray-500">Cancelar</button>
		[[ if .{{.Tabla.PrimerCampo.Nombre.Camel}} ]]
		<button type="button" hx-delete="" hx-confirm="¿Eliminar?" class="inline-block px-6 py-2 font-medium transition-colors bg-red-600 rounded shadow-md text-slate-100 hover:bg-red-500">Eliminar</button>
		<button type="button" hx-put="" class="inline-block px-6 py-2 font-medium transition-colors bg-blue-700 rounded shadow-md text-slate-100 hover:bg-blue-600">Guardar</button>
		[[ else ]]
		<button type="button" hx-post="" class="inline-block px-6 py-2 font-medium transition-colors bg-blue-700 rounded shadow-md text-slate-100 hover:bg-blue-600">Registrar</button>
		[[ end ]]
	</div>
	
</form>

[[ end ]]