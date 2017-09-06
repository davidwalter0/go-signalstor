package main

import (
	"fmt"
	"github.com/davidwalter0/xml2json"
	"io/ioutil"
	"os"
)

func main() {

	ftp := xml2json.Configure()
	var err error
	var rawData []byte
	var messages xml2json.SmsMessages

	rawData, err = ioutil.ReadFile(ftp.Filename)
	if err != nil {
		fmt.Printf("error: Download failed %v\n", err)
		os.Exit(-1)
	}
	// xml2json.XMLParseArray(rawData, &messages, xml2json.SmsXmlFixUp, xml2json.NoOp)
	xml2json.XMLParse(rawData, &messages, xml2json.SmsXmlFixUp, xml2json.NoOp)
	xml2json.DumpParsedMessages(os.Stderr, messages)
	smsDbIO := xml2json.NewSmsDbIO()
	for _, msg := range messages.Messages {
		fmt.Println(msg)
		smsDbIO.CopySmsMessage(&msg)
		smsDbIO.Create()
	}
}
