package ddd

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func NewPaqueteID() int {
	num, err := rand.Int(rand.Reader, big.NewInt(9999))
	if err != nil {
		fmt.Println("Error al generar un número aleatorio", err)
		return 0
	}
	return int(num.Int64())
}

func NewTablaID() int {
	num, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		fmt.Println("Error al generar un número aleatorio", err)
		return 0
	}
	return int(num.Int64())
}

func NewCampoID() int {
	num, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		fmt.Println("Error al generar un número aleatorio", err)
		return 0
	}
	return int(num.Int64())
}

func NewConsultaID() int {
	num, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		fmt.Println("Error al generar un número aleatorio", err)
		return 0
	}
	return int(num.Int64())
}
