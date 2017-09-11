package main

import (
	"fmt"
	"github.com/davidwalter0/xml2json"
	"io/ioutil"
	"os"
)

func main() {
	ftp := xml2json.ConfigureFtp()
	var err error
	var rawData []byte
	var messages xml2json.SmsMessages

	rawData, err = ioutil.ReadFile(ftp.Filename)
	if err != nil {
		fmt.Printf("error: Download failed %v\n", err)
		os.Exit(-1)
	}
	xml2json.XMLParse(rawData, &messages, xml2json.SmsXMLFixUp, xml2json.NoOp)
	xml2json.DumpParsedMessages(os.Stderr, messages)
}
