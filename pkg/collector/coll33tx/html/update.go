package main

import (
	"log"

	"github.com/florinutz/filme/pkg/collector/coll33tx/html/mockloader"
)

const dataFile = "data.json"

func main() {
	err := mockloader.Write(dataFile)
	if err != nil {
		log.Fatal(err)
	}
}
