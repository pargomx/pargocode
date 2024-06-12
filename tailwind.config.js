/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
	"./pargo/plantillas/*.html",
	"./pargo/plantillas/**/*.html",
	
],
  theme: {
    extend: {
		cursor: {
			'edit': 'url("/assets/img/pen.svg"), pointer',
		},
	},
  },
  plugins: [],
}

