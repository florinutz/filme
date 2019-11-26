package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/florinutz/filme/pkg/collector"
	"github.com/florinutz/filme/pkg/collector/coll33tx/detail"
	"github.com/florinutz/filme/pkg/filme/l33tx/list"
)

const dataFile = "test-data"

var urls = []string{
	detail.TestPageDetail,
	list.TestPageListWithPagination,
	list.TestPageList,
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
