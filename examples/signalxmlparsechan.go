package main

import (
	"os"

	"github.com/davidwalter0/go-signalstor"
)

func configure() (filename string) {
	ftp := signalstor.ConfigureFtp()
	version()
	return ftp.Filename
}

func main() {
	var done = make(chan bool)

	var messages = make(chan *signalstor.SmsMessage)

	go signalstor.XMLParsePublish(configure(), messages)
	go signalstor.DumpParsedMessagesSubscribe(os.Stderr, messages, done)
	<-done
}
