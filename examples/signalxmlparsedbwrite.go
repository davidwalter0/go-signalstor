package main

import (
	"fmt"
	"github.com/davidwalter0/go-signalstor"
	"io/ioutil"
	"os"
)

func main() {

	ftp := signalstor.ConfigureFtp()
	var err error
	var rawData []byte
	var messages signalstor.SmsMessages

	rawData, err = ioutil.ReadFile(ftp.Filename)
	if err != nil {
		fmt.Printf("error: Download failed %v\n", err)
		os.Exit(-1)
	}
	signalstor.XMLParse(rawData, &messages, signalstor.SmsXMLFixUp, signalstor.NoOp)
	smsDbIO := signalstor.NewSmsDbIO()
	fmt.Println("here", len(messages.Messages))
	for _, msg := range messages.Messages {
		fmt.Println(msg)
		smsDbIO.CopySmsMessage(&msg)
		smsDbIO.Create()
	}
}
