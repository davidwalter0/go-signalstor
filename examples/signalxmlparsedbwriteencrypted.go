package main

import (
	"fmt"
	"github.com/davidwalter0/xml2json"
	"io/ioutil"
	"os"
)

func main() {

	var key = []byte("LKHlhb899Y09olUi")

	smsRead := xml2json.NewSmsDbIO()

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
	smsDbIO := xml2json.NewSmsDbIO()

	for _, msg := range messages.Messages {
		fmt.Println(msg)
		smsDbIO.CopySmsMessage(&msg)
		smsRead.CopyKey(smsDbIO)

		fmt.Println(smsDbIO)
		fmt.Println(smsRead)

		if err = smsDbIO.Encrypt(key); err != nil {
			fmt.Fprintf(os.Stderr, "Encrypt() failed %v\n", err)
			continue
		}

		smsDbIO.Create()

		fmt.Println(smsDbIO)
		fmt.Println(smsRead)

		if err = smsRead.Read(); err != nil {
			fmt.Printf("smsRead.Read() failed %v", err)
		} else {
			if err = smsRead.Decrypt(key); err != nil {
				fmt.Printf("smsRead.Decrypt(key) failed %v", err)
				continue
			}
			fmt.Println(">>Decrypted", smsRead)
		}
	}
}
