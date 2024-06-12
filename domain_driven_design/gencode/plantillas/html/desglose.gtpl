<!DOCTYPE html>
<html lang="es-MX">
	<head>
		##include('meta.html')
		##include('links.html')
	</head>
	<body>
		##include('preloader.html')
		##include('header.html')
		<div class="main">
			##include('navigation.html', {
				"link": "",
			})
			<section class="content">
				<div class="content__inner content__inner--sm">

					{{ range $.Tabla.CamposEspeciales -}}
					[[ $lista{{.Nombre.Camel}} := .Lista{{.Nombre.Camel}} ]]
					{{ end }}

					[[ with .{{ $.Tabla.Nombre.Camel }} ]]
					@@include('{{ $.Tabla.Nombre.ClavePlural }}/detalles.html', {
						"linkEditar": "./editar",
					})
					[[ end ]]

					##include('footer.html')
				</div>
			</section>
		</div>

		##include('scripts.html',{
			"autosize_text": false,
		})
	</body>
  </html>