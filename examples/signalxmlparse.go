package main

import (
	"fmt"
	"github.com/davidwalter0/go-signalstor"
	"io/ioutil"
	"os"
)

func main() {
	ftp := signalstor.ConfigureFtp()
	version()
	var err error
	var rawData []byte
	var messages signalstor.SmsMessages

	rawData, err = ioutil.ReadFile(ftp.Filename)
	if err != nil {
		fmt.Printf("error: Download failed %v\n", err)
		os.Exit(-1)
	}

	signalstor.XMLParse(rawData, &messages, signalstor.SmsXMLFixUp, signalstor.SmsMessageValidate)
	signalstor.DumpParsedMessages(os.Stderr, messages)
}
