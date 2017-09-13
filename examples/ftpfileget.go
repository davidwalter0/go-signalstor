package main

import (
	"log"

	"github.com/davidwalter0/go-signalstor"
)

func main() {
	writingToFile := true
	cfg := signalstor.ConfigureFtp()
	if _, err := signalstor.Download(cfg, writingToFile); err != nil {
		log.Fatal(err)
	}
}
