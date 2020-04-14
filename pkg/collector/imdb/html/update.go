package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/florinutz/filme/pkg/collector/gz_http"
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

	reqs, err := gz_http.GenerateSimpleRequests(urls, nil)
	if err != nil {
		log.Fatal(err)
	}

	errs := gz_http.UpdateTestData(reqs, outputPath)
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, err)
	}
}
