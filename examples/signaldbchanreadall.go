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
	configure()
	var messages = make(chan *signalstor.SmsMessage)
	var smsDbIO = signalstor.NewSmsDbIO()
	smsDbIO.Msg.Address = "+15555555555"

	go func() {
		defer close(messages)
		smsDbIO.ReadAll(messages)
	}()

	go signalstor.DumpParsedMessagesSubscribe(os.Stderr, messages, done)
	<-done
}
