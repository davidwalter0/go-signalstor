package main

import (
	"github.com/davidwalter0/xml2json"
)

func configure() (filename string) {
	ftp := xml2json.ConfigureFtp()
	return ftp.Filename
}

func main() {
	var done = make(chan bool)

	var messages = make(chan *xml2json.SmsMessage)

	go xml2json.XMLParsePublish(configure(), messages)
	go func() {
		var smsDbIO = xml2json.NewSmsDbIO()
		for message := range messages {
			smsDbIO.Msg = *message
			smsDbIO.Create()
		}
		done <- true
	}()
	// go xml2json.DumpParsedMessagesSubscribe(os.Stderr, messages, done)
	<-done
}
