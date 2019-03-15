package signalstor // 	"github.com/davidwalter0/go-signalstor"

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/davidwalter0/go-signalstor"
)

func init() {
	var envCfg = map[string]string{
		// "FTP_FILENAME":          "SignalPlainXMLTextBackup.xml",
		// "FTP_HOST":              "192.168.0.12",
		// "FTP_PORT":              "2121",
		// "FTP_USER":              "ftp",
		// "FTP_PASSWORD":          "ftp",
		// "FTP_PHONE":             "+15555555555",
		"SQL_DRIVER":            "postgres",
		"SQL_HOST":              "localhost",
		"SQL_PORT":              "5432",
		"SQL_DATABASE":          "sms",
		"SQL_USER":              "sms",
		"SQL_PASSWORD":          "sms",
		"SQL_TIMEOUT":           "0",
		"SQL_SCHEMA_INITIALIZE": "false",
	}

	for k, v := range envCfg {
		if err := os.Setenv(k, v); err != nil {
			panic(fmt.Sprintf("failed configuring the environment: %v", err))
		}
	}

	fmt.Printf("%v\n", signalstor.ConfigureDb())
}

func TestSmsCopyMessage(t *testing.T) {
	var expectKey = &signalstor.SmsDbIO{
		Msg: signalstor.SmsMessage{
			Address:   "address",
			Timestamp: "date",
		},
	}

	var expectSmsDbIO = &signalstor.SmsDbIO{
		ID:      0,
		GUID:    "guid",
		Created: time.Time{},
		Changed: time.Time{},
		Msg: signalstor.SmsMessage{
			ContactName: "contact_name",
			Timestamp:   "date",
			Date:        "readable_date",
			Address:     "address",
			Subject:     "subject",
			Body:        "body",
			Type:        "1",
		},
	}

	var key = &signalstor.SmsDbIO{}
	key.CopyKey(expectKey)
	if *key != *expectKey {
		t.Fatalf("SmsDbIOKey.CopyKey failed: wanted\n%v\ngot\n%v\n", *expectKey, *key)
	}

	var smsDbIo = &signalstor.SmsDbIO{}
	smsDbIo.CopySmsDbIO(expectSmsDbIO)
	if *smsDbIo != *expectSmsDbIO {
		t.Fatalf("SmsDbIOKey.CopySmsDbIO failed: wanted\n%v\ngot\n%v\n", *expectSmsDbIO, *smsDbIo)
	}
}

/*
var XMLText = `<sms protocol="0" address="+15555555555" contact_name="Self" date="1493140602697" readable_date="Tue, 25 Apr 2017 13:16:42 EDT" type="1" subject="null" body="body XMLText" toa="null" sc_toa="null" service_center="null" read="1" status="-1" locked="0" />
<sms protocol="0" address="22000" contact_name="null" date="1493139630014" readable_date="Tue, 25 Apr 2017 13:00:30 EDT" type="1" subject="null" body="Account notification: for u@abc.123" toa="null" sc_toa="null" service_center="null" read="1" status="-1" locked="0" />`
*/

var XMLText = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<!-- File Created By Signal -->
<smses count="2">
  <sms protocol="0" address="+15555555555" contact_name="Self" date="1493140602697" readable_date="Tue, 25 Apr 2017 13:16:42 EDT" type="1" subject="null" body="body XMLText" toa="null" sc_toa="null" service_center="null" read="1" status="-1" locked="0" />
  <sms protocol="0" address="22000" contact_name="null" date="1493139630014" readable_date="Tue, 25 Apr 2017 13:00:30 EDT" type="1" subject="null" body="Account notification: for u@abc.123" toa="null" sc_toa="null" service_center="null" read="1" status="-1" locked="0" />
</smses>
`)

func JsonDump(st interface{}) {
	lhs, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	fmt.Println(string(lhs))
	return
}

func NoFixUp(xml string) (string, error) {
	return xml, nil
}

func TestSmsParseArray(t *testing.T) {

	var messages = createMessages()

	lhs, err := json.Marshal(messages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	var wantedJSON = `{"sms":[{"contact_name":"Self","date":"1493140602697","readable_date":"Tue, 25 Apr 2017 13:16:42 EDT","address":"+15555555555","subject":"null","body":"body XMLText","type":"1"},{"contact_name":"null","date":"1493139630014","readable_date":"Tue, 25 Apr 2017 13:00:30 EDT","address":"22000","subject":"null","body":"Account notification: for u@abc.123","type":"1"}]}`

	type Msg map[string][]map[string]string
	var wanted Msg
	var got Msg

	if err = json.Unmarshal([]byte(wantedJSON), &wanted); err != nil {
		t.Fatalf("json.Unmarshal fail %v", err)
	}
	if err = json.Unmarshal(lhs, &got); err != nil {
		t.Fatalf("json.Unmarshal fail %v", err)
	}

	if !reflect.DeepEqual(wanted, got) {
		t.Fatalf("parse object failed: wanted\n%s\ngot\n%s\n", wanted, got)
	}
}

func createMessages() signalstor.SmsMessages {
	var messages signalstor.SmsMessages
	signalstor.XMLParse(XMLText, &messages, signalstor.SmsXMLFixUp, signalstor.SmsMessageValidate)
	_, err := json.Marshal(messages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	return messages
}

func TestSmsIO(t *testing.T) {
	var err error
	messages := createMessages()

	smsIO := signalstor.NewSmsDbIO()
	if smsIO == nil {
		t.Fatalf("smsIO db object create is nil")
	}

	if smsIO.ConfigureDb() == nil {
		t.Fatalf("smsIO configure failed")
	}

	smsRead := signalstor.NewSmsDbIO()
	if smsRead == nil {
		t.Fatalf("smsRead db object create is nil")
	}

	for _, x := range messages.Messages {
		fmt.Printf("x %v\n", x)
		fmt.Printf("key %v\n", string(key))
		if err := smsIO.CopySmsMessage(&x).Encrypt(key); err != nil {
			t.Fatalf("Encrypt() failed%v", err)
		}
		if err = smsIO.Delete(); err != nil {
			t.Fatalf("smsIO.Delete() failed %v", err)
		}
		fmt.Printf("smsIO %v\n", *smsIO)
		if err = smsIO.Create(); err != nil {
			t.Fatalf("smsIO.Create() failed %v", err)
		}
		smsRead.CopyKey(smsIO)
		if err = smsRead.Read(); err != nil {
			t.Fatalf("smsRead.Read() failed %v", err)
		} else {
			if err = smsRead.Decrypt(key); err != nil {
				t.Fatalf("smsRead.Decrypt(key) failed %v", err)
				continue
			}
			if err = smsIO.Decrypt(key); err != nil {
				t.Fatalf("smsIO.Decrypt(key) failed %v", err)
				continue
			}

			wanted := *smsIO
			got := *smsRead
			if wanted.Msg != got.Msg {
				t.Fatalf("query returned different data\n%v\n%v\n", wanted.Msg, got.Msg)
			}
		}
	}
}
