package main

import (
	"flag"
	"monorepo/gkfmt/gkfmt"
	"os"
	"strings"
	"time"

	"github.com/pargomx/gecko/gko"
)

func main() {
	gko.PrintLogTimestamps = false

	inputFile := flag.String("i", "input.html", "Input file path")
	outputFile := flag.String("o", "output.html", "Output file path")
	useTokens := flag.Bool("t", false, "Format using only tokens (dumb)")
	flag.Parse()

	start := time.Now()

	bytes, err := os.ReadFile(*inputFile)
	if err != nil {
		gko.FatalError(err)
	}

	var builder strings.Builder
	gkfmt.FormatGeckoTemplate(string(bytes), &builder, *useTokens)

	file, err := os.Create(*outputFile)
	if err != nil {
		gko.FatalError(err)
	}
	defer file.Close()
	_, err = file.WriteString(builder.String())
	if err != nil {
		gko.FatalError(err)
	}

	gko.LogEventof("Formated in %v", time.Since(start))
}
