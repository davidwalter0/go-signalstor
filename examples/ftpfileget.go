package main

import (
	"log"

	"github.com/davidwalter0/xml2json"
)

func main() {
	writingToFile := true
	cfg := xml2json.ConfigureFtp()
	if _, err := xml2json.Download(cfg, writingToFile); err != nil {
		log.Fatal(err)
	}
}
