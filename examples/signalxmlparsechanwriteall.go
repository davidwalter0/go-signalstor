package main

import (
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
	go func() {
		var smsDbIO = signalstor.NewSmsDbIO()
		for message := range messages {
			smsDbIO.Msg = *message
			smsDbIO.Create()
		}
		done <- true
	}()
	// go signalstor.DumpParsedMessagesSubscribe(os.Stderr, messages, done)
	<-done
}
