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
				<div class="content__inner">

					{{ range $.Tabla.CamposEspeciales -}}
					[[ $lista{{.Nombre.Camel}} := .Lista{{.Nombre.Camel}} ]]
					{{ end }}

					@@include('{{$.Tabla.Nombre.ClavePlural}}/tabla.html', {

					})

					##include('footer.html')
				</div>
			</section>
		</div>

		##include('scripts.html',{
			"datatables": true,
		})
	</body>
  </html>