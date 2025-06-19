package sqlitedb

import (
	"fmt"
	"reflect"
	"strings"
)

func ArgToText(arg any) string {
	if arg == nil {
		return "NULL" // nil directo
	}
	valueOf := reflect.ValueOf(arg)
	switch {

	case valueOf.Kind() == reflect.String:
		return fmt.Sprintf("'%s'", arg) // string

	case valueOf.Kind() == reflect.Ptr:
		if valueOf.IsNil() {
			return "NULL" // nil pointer
		}
		valueElem := valueOf.Elem()
		if valueElem.Kind() == reflect.String {
			return fmt.Sprintf("'%s'", valueElem.String()) // string pointer
		}
		return fmt.Sprint(valueElem.Interface()) // other pointer

	default:
		return fmt.Sprintf("%v", arg)
	}
}

const (
	tipoQuery uint8 = iota
	tipoQueryRow
	tipoExec
	tipoTX
)

func logSQL(tipo uint8, qry string, args ...any) {
	const reset = "\033[0m"
	const color = "\033[36m"
	const bold = "\033[1;34m"
	if strings.Count(qry, "?") != len(args) {
		fmt.Println(bold+"[QUERY]"+reset+" ", qry, color, "tiene numero incorrecto de argumentos"+reset)
		return
	}
	for _, arg := range args {
		qry = strings.Replace(qry, "?", color+ArgToText(arg)+reset, 1)
	}
	switch tipo {
	case tipoQuery:
		fmt.Println(bold + "[QRY]" + reset + " " + qry + ";")
	case tipoQueryRow:
		fmt.Println(bold + "[ROW]" + reset + " " + qry + ";")
	case tipoExec:
		fmt.Println(bold + "[EXE]" + reset + " " + qry + ";")
	}
}
