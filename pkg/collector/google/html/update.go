package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/florinutz/filme/pkg/collector/google/search"
	"github.com/florinutz/filme/pkg/collector/gz_http"
)

const dataFile = "test-data"

var urls = []string{search.TestSearch}

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
