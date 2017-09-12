package main

import (
	"fmt"
	"github.com/davidwalter0/xml2json"
	"io/ioutil"
	"os"
)

func loadData() []byte {
	ftp := xml2json.ConfigureFtp()
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
	ftp := xml2json.ConfigureFtp()
	return ftp.Filename
}

func main() {
	var done = make(chan bool)

	var messages = make(chan *xml2json.SmsMessage)

	go xml2json.XMLParsePublish(configure(), messages)
	go xml2json.DumpParsedMessagesSubscribe(os.Stderr, messages, done)
	<-done
}
