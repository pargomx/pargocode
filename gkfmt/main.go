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

func main() {
	gko.PrintLogTimestamps = false

	inputFile := flag.String("i", "input.html", "Input file path")
	outputFile := flag.String("o", "output.html", "Output file path")
	flag.Parse()

	start := time.Now()

	bytes, err := os.ReadFile(*inputFile)
	if err != nil {
		gko.FatalError(err)
	}

	startParse := time.Now()
	tokens := gkfmt.ParseTokens(string(bytes))
	timeParse := time.Since(startParse)

	var builder strings.Builder
	for _, token := range tokens {
		builder.WriteString(fmt.Sprintf("%v%v\n", strings.Repeat("\t", token.Indent), token.Txt))
	}

	file, err := os.Create(*outputFile)
	if err != nil {
		gko.FatalError(err)
	}
	defer file.Close()
	_, err = file.WriteString(builder.String())
	if err != nil {
		gko.FatalError(err)
	}

	gko.LogEventof("Saved %v tokens in %v (parsed in %v)", len(tokens), time.Since(start), timeParse)
}
