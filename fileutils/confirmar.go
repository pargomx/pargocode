package fileutils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func confirmadoPorUsuario(path string) bool {
	fmt.Printf("escribir \033[34m%s\033[0m [Y/n] ", path)
	res := ""
	scanner := bufio.NewScanner(os.Stdin)
	ok := scanner.Scan()
	if ok {
		res = strings.TrimRight(scanner.Text(), "\r\n")
	}
	if !(res == "" || res == "y" || res == "Y") {
		return true
	}
	return false
}
