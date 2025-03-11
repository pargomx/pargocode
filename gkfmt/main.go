package main

import (
	"flag"
	"fmt"
	"monorepo/gkfmt/gkfmt"
	"os"
	"strings"
	"time"

	"github.com/pargomx/gecko/gko"
)

var BUILD_INFO string

func main() {
	gko.PrintLogTimestamps = false

	if info := strings.Split(BUILD_INFO, ":"); len(info) == 3 {
		fmt.Printf("\033[1;36m%v\033[0;36m v%v (%v)\033[0m\n", info[0], info[1], info[2])
	}

	var inputFile, outputFile string

	flag.StringVar(&inputFile, "i", "input.html", "Input file path")
	flag.StringVar(&outputFile, "o", "", "Output file path. Default: rewrite input")
	useTokens := flag.Bool("t", false, "Format using only tokens (dumb)")
	debug := flag.Bool("d", false, "Debug mode: print each token")
	flag.Parse()

	start := time.Now()

	if !strings.HasSuffix(inputFile, ".html") {
		gko.FatalExit("Solo se admiten archivos HTML")
	}

	bytes, err := os.ReadFile(inputFile)
	if err != nil {
		gko.FatalError(err)
	}

	var builder strings.Builder
	gkfmt.FormatGeckoTemplate(string(bytes), &builder, *useTokens, *debug)

	if outputFile == "" {
		outputFile = inputFile
	}
	file, err := os.Create(outputFile)
	if err != nil {
		gko.FatalError(err)
	}
	defer file.Close()
	_, err = file.WriteString(builder.String())
	if err != nil {
		gko.FatalError(err)
	}

	fmt.Printf("Formated in %v\n", time.Since(start))
}
