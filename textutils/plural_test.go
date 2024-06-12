package textutils

import "testing"

func TestDeducirNombrePlural(t *testing.T) {

	type caso struct {
		input    string
		esperado string
	}

	casos := []caso{
		{input: "otro", esperado: "otros"},
		{input: "cosa", esperado: "cosas"},
		{input: "nuevo carro", esperado: "nuevos carros"},
		{input: "juego de azar", esperado: "juegos de azar"},
		{input: "flor de bach", esperado: "flores de bach"},
		{input: "datos personales", esperado: "datos personales"},
		{input: "Tipo de calificación", esperado: "Tipos de calificación"},

		{input: "aprendiz", esperado: "aprendices"},
		{input: "asesoria", esperado: "asesorias"},
		{input: "asignatura", esperado: "asignaturas"},
		{input: "calificacion", esperado: "calificaciones"},
		{input: "calificación", esperado: "calificaciones"},
		{input: "contacto", esperado: "contactos"},
		{input: "contenido", esperado: "contenidos"},
		{input: "convenio", esperado: "convenios"},
		{input: "dbancarios", esperado: "dbancarios"},
		{input: "dfacturacion", esperado: "dfacturaciones"},
		{input: "dpersonales", esperado: "dpersonales"},
		{input: "dprofesionales", esperado: "dprofesionales"},
		{input: "encuentro", esperado: "encuentros"},
		{input: "egreso", esperado: "egresos"},
		{input: "examen", esperado: "examenes"},
		{input: "expediente", esperado: "expedientes"},
		{input: "experiencia", esperado: "experiencias"},
		{input: "falta", esperado: "faltas"},
		{input: "inscrito", esperado: "inscritos"},
		{input: "interesado", esperado: "interesados"},
		{input: "materia", esperado: "materias"},
		{input: "mentor", esperado: "mentores"},
		{input: "oferta", esperado: "ofertas"},
		{input: "pago", esperado: "pagos"},
		{input: "programa", esperado: "programas"},
		{input: "sesion", esperado: "sesiones"},
		{input: "tabulador", esperado: "tabuladores"},
		{input: "titulacion", esperado: "titulaciones"},
		{input: "usuario", esperado: "usuarios"},
		{input: "Titulación", esperado: "Titulaciones"},

		// mas casos input

		{input: "Enlace", esperado: "Enlaces"},
		{input: "Evaluación", esperado: "Evaluaciones"},
		{input: "Eficaz", esperado: "Eficaces"},
		{input: "Faz", esperado: "Faces"},
		{input: "Lapso", esperado: "Lapsos"},
		{input: "Luz", esperado: "Luces"},
		{input: "Voz", esperado: "Voces"},
		{input: "Raíz", esperado: "Raíces"},
	}

	for _, c := range casos {
		got := DeducirNombrePlural(c.input)
		if got != c.esperado {
			t.Errorf("from %q got %q, wanted %q", c.input, got, c.esperado)
		}
	}
}
