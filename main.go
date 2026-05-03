package main

import (
	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/pargocode/pargocode"
)

// Información de compilación establecida con:
//
//	BUILD_INFO="$(date -I):$(git log --format="%H" -n 1)"
//	go build -ldflags "-X main.BUILD_INFO=$BUILD_INFO -X main.AMBIENTE=DEV"
var BUILD_INFO string // Información de compilación [ fecha:commit_hash ]
var AMBIENTE string   // Ambiente de ejecución [ DEV / PROD ]

func main() {
	gko.LogInfof("Versión:%s:%s", BUILD_INFO, AMBIENTE)
	pargocode.Run()
}
