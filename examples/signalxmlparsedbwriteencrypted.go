package main

import (
	"fmt"
	"github.com/davidwalter0/go-signalstor"
	"io/ioutil"
	"os"
)

func main() {

	var key = []byte("LKHlhb899Y09olUi")

	ftp := signalstor.ConfigureFtp()
	version()
	smsRead := signalstor.NewSmsDbIO()
	fmt.Println("ftp.Password", ftp.Password)
	fmt.Println("ftp.Key", ftp.Key)
	var err error
	var rawData []byte
	var messages signalstor.SmsMessages

	rawData, err = ioutil.ReadFile(ftp.Filename)
	if err != nil {
		fmt.Printf("error: Download failed %v\n", err)
		os.Exit(-1)
	}

	signalstor.XMLParse(rawData, &messages, signalstor.SmsXMLFixUp, signalstor.SmsMessageValidate)
	// signalstor.XMLParse(rawData, &messages, signalstor.SmsXMLFixUp, signalstor.NoOp)
	smsDbIO := signalstor.NewSmsDbIO()

	for _, msg := range messages.Messages {
		if len(msg.Address) == 0 || len(msg.Timestamp) == 0 {
			continue
		}
		smsDbIO.CopySmsMessage(&msg)
		smsRead.CopyKey(smsDbIO)

		if err = smsDbIO.Encrypt(key); err != nil {
			fmt.Fprintf(os.Stderr, "Encrypt() failed %v\n", err)
			continue
		}

		if err = smsDbIO.Create(); err != nil {
			fmt.Println("*Error*", err)
			continue
		}

		if err = smsRead.Read(); err != nil {
			fmt.Printf("smsRead.Read() failed %v", err)
		} else {
			if err = smsRead.Decrypt(key); err != nil {
				fmt.Printf("smsRead.Decrypt(key) failed %v", err)
				continue
			}
			fmt.Printf("\nDecrypted\n%s", smsRead.String())
		}
	}
}
