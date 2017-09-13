package main

import (
	"os"

	"github.com/davidwalter0/xml2json"
)

func configure() (filename string) {
	ftp := xml2json.ConfigureFtp()
	return ftp.Filename
}

func main() {
	var done = make(chan bool)

	var messages = make(chan *xml2json.SmsMessage)
	var smsDbIO = xml2json.NewSmsDbIO()
	smsDbIO.Msg.Address = "+15555555555"

	go func() {
		defer close(messages)
		smsDbIO.ReadAll(messages)
	}()

	go xml2json.DumpParsedMessagesSubscribe(os.Stderr, messages, done)
	<-done
}
