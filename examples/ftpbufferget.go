package main

import (
	"log"

	"github.com/davidwalter0/xml2json"
)

func main() {
	writingToFile := false
	cfg := xml2json.ConfigureFtp()
	if content, err := xml2json.Download(cfg, writingToFile); err != nil {
		log.Fatal(err)
	} else {
		log.Println(string(content))
	}
}
