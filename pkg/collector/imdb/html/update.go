package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/florinutz/filme/pkg/collector"
	"github.com/florinutz/filme/pkg/collector/imdb/detail"
)

const dataFile = "test-data"

var urls = []string{
	detail.TestPageFilm,
	detail.TestPageSeriesFinished,
	detail.TestPageSeriesUnfinished,
}

func main() {
	outputPath := dataFile
	if len(os.Args) > 1 {
		outputPath = strings.Join(os.Args[1:], " ")
	}

	reqs, err := collector.GenerateSimpleRequests(urls, nil)
	if err != nil {
		log.Fatal(err)
	}

	errs := collector.UpdateTestData(reqs, outputPath)
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, err)
	}
}
