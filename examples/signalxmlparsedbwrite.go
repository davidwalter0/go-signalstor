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
	smsRead := signalstor.NewSmsDbIO()

	var err error
	var rawData []byte
	var messages signalstor.SmsMessages

	
	if rawData, err = ioutil.ReadFile(ftp.Filename); err != nil {
		fmt.Printf("error: Download failed %v\n", err)
		os.Exit(-1)
	}

	signalstor.XMLParse(rawData, &messages, signalstor.SmsXMLFixUp, signalstor.NoOp)
	smsDbIO := signalstor.NewSmsDbIO()

	for _, msg := range messages.Messages {
		if ! msg.IsValid() {
			continue
		}
		smsDbIO.CopySmsMessage(&msg)
		smsRead.CopyKey(smsDbIO)

		if err = smsDbIO.Create(); err != nil {
			fmt.Println("*Error*", err)
      panic(err)

			continue
		}

		if err = smsRead.Read(); err != nil {
			fmt.Printf("smsRead.Read() failed %v", err)
		}
	}
}
