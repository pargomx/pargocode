package sqlitedb_test

import (
	"monorepo/sqlitedb"
	"testing"
)

func TestLogSQLArg(t *testing.T) {
	type testcase struct {
		name string
		arg  any
		want string
	}
	str := "hola"
	flg := true
	num := 123
	dec := 123.456
	itm := struct {
		nullField  *int
		zeroString string
		zeroBool   bool
		zeroInt    int
		zeroFloat  float64
	}{}
	tests := []testcase{

		{name: "nil", arg: nil, want: "NULL"},
		{name: "nil field", arg: itm.nullField, want: "NULL"},

		{name: "string", arg: str, want: "'hola'"},
		{name: "zero string", arg: itm.zeroString, want: "''"},
		{name: "*string", arg: &str, want: "'hola'"},

		{name: "bool", arg: flg, want: "true"},
		{name: "zero bool", arg: itm.zeroBool, want: "false"},
		{name: "*bool", arg: &flg, want: "true"},

		{name: "int", arg: num, want: "123"},
		{name: "zero int", arg: itm.zeroInt, want: "0"},
		{name: "*int", arg: &num, want: "123"},

		{name: "float", arg: dec, want: "123.456"},
		{name: "zero float", arg: itm.zeroFloat, want: "0"},
		{name: "*float", arg: &dec, want: "123.456"},
	}
	for _, tt := range tests {
		got := sqlitedb.ArgToText(tt.arg)
		if got != tt.want {
			t.Errorf("LogSQLArg(%v) got %v; want %v", tt.name, got, tt.want)
		}
	}
}
