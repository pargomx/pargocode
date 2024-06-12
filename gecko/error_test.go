package gecko

import (
	"errors"
	"testing"
)

func TestEsErrNotFound(t *testing.T) {

	if EsErrNotFound(nil) {
		t.Fatal("EsErrNotFound(nil) was true")
	}

	var err1 error
	if EsErrNotFound(err1) {
		t.Fatal("EsErrNotFound(err1) was true")
	}

	err2 := errors.New("error genérico")
	if EsErrNotFound(err2) {
		t.Fatal("EsErrNotFound(err2) was true")
	}

	err3 := NewErr(400)
	if EsErrNotFound(err3) {
		t.Fatal("EsErrNotFound(err3) was true")
	}

	err4 := NewErr(404)
	if !EsErrNotFound(err4) {
		t.Fatal("EsErrNotFound(err4) was false")
	}

	err5 := &Gkerror{}
	if EsErrNotFound(err5) {
		t.Fatal("EsErrNotFound(err5) was true")
	}

	err6 := &Gkerror{}
	err6 = nil
	if EsErrNotFound(err6) {
		t.Fatal("EsErrNotFound(err6) was true")
	}

	err7 := NewErr(404).Code(500)
	if EsErrNotFound(err7) {
		t.Fatal("EsErrNotFound(err7) was true")
	}

	err8 := NewErr(500).Code(404)
	if !EsErrNotFound(err8) {
		t.Fatal("EsErrNotFound(err8) was false")
	}

	err9 := NewErr(500).Code(404).Code(501).Code(404)
	if !EsErrNotFound(err9) {
		t.Fatal("EsErrNotFound(err9) was false")
	}
}
func BenchmarkEsErrNotFound1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := NewErr(404)
		if !err.EsNotFound() {
			b.Fatal("debería ser true")
		}
	}
}

func BenchmarkEsErrNotFound2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var err error = NewErr(404)
		if !Err(err).EsNotFound() {
			b.Fatal("debería ser true")
		}
	}
}
func BenchmarkEsErrNotFound3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var err error = NewErr(404)
		if !EsErrNotFound(err) {
			b.Fatal("debería ser true")
		}
	}
}
