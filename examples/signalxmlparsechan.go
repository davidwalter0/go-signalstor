package main

import (
	"fmt"
	"github.com/davidwalter0/go-signalstor"
	"io/ioutil"
	"os"
)

func loadData() []byte {
	ftp := signalstor.ConfigureFtp()
	var err error
	var rawData []byte
	rawData, err = ioutil.ReadFile(ftp.Filename)
	if err != nil {
		fmt.Printf("error: Download failed %v\n", err)
		os.Exit(-1)
	}
	return rawData
}

func configure() (filename string) {
	ftp := signalstor.ConfigureFtp()
	return ftp.Filename
}

func main() {
	var done = make(chan bool)

	var messages = make(chan *signalstor.SmsMessage)

	go signalstor.XMLParsePublish(configure(), messages)
	go signalstor.DumpParsedMessagesSubscribe(os.Stderr, messages, done)
	<-done
}
