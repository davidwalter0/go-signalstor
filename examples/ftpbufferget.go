package main

import (
	"log"

	"github.com/davidwalter0/go-signalstor"
)

func main() {
	writingToFile := false
	cfg := signalstor.ConfigureFtp()
	version()

	if content, err := signalstor.Download(cfg, writingToFile); err != nil {
		log.Fatal(err)
	} else {
		log.Println(string(content))
	}
}
