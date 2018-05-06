package signalstor // 	"github.com/davidwalter0/go-signalstor"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/davidwalter0/go-signalstor"
)

func init() {
	var envCfg = map[string]string{
		// "FTP_FILENAME":          "SignalPlainJSONTextBackup.xml",
		// "FTP_HOST":              "192.168.0.12",
		// "FTP_PORT":              "2121",
		// "FTP_USER":              "ftp",
		// "FTP_PASSWORD":          "ftp",
		// "FTP_PHONE":             "+15555555555",
		"SQL_DRIVER":            "postgres",
		"SQL_HOST":              "localhost",
		"SQL_PORT":              "5432",
		"SQL_DATABASE":          "sms",
		"SQL_USER":              "USER_ID",
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

var originSmsDbIO = &signalstor.SmsDbIO{
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

func MakeMsg() *signalstor.SmsMessage {
	var msgJSON = []byte(`{"contact_name":"Self","date":"1493140602697","readable_date":"Tue, 25 Apr 2017 13:16:42 EDT","address":"+15555555555","subject":"null","body":"body JSONText","type":"1"}`)

	var msg = signalstor.SmsMessage{}
	var err error
	if err = json.Unmarshal(msgJSON, &msg); err != nil {
		fmt.Println(err)
	}
	return &msg
}

var JSONText = `{"sms":[{"contact_name":"Self","date":"1493140602697","readable_date":"Tue, 25 Apr 2017 13:16:42 EDT","address":"+15555555555","subject":"null","body":"body JSONText","type":"1"},{"contact_name":"null","date":"1493139630014","readable_date":"Tue, 25 Apr 2017 13:00:30 EDT","address":"22000","subject":"null","body":"Account notification: for u@abc.123","type":"1"}]}`

func TestEncrypt(t *testing.T) {
	var err error
	var nti *signalstor.NaclTool
	if nti = signalstor.NewEncryptionNaclTool([]byte(JSONText)); nti == nil {
		t.Fatalf(fmt.Sprintf("NewEncryptionNaclTool returned nil"))
	}

	if err = nti.NaclEncrypt(); err != nil {
		t.Fatalf(fmt.Sprintf("NaclEncrypt %v", err))
	}
	if debug {
		fmt.Println(string((*nti.Key)[:]))
		fmt.Println(string((*nti.Nonce)[:]))
		fmt.Println(string(*nti.CypherMessage))
	}
	var key = nti.Key
	var msg = nti.CypherMessage
	var nto *signalstor.NaclTool

	if nto = signalstor.NewDecryptionNaclTool(key, *msg); nto == nil {
		t.Fatalf(fmt.Sprintf("NewDecryptionNaclTool returned nil"))
	}

	if err = nto.NaclDecrypt(); err != nil {
		t.Fatalf(fmt.Sprintf("NaclEncrypt %v", err))
	}

	if !bytes.Equal(*nto.PlainMessage, *nti.PlainMessage) {
		t.Fatalf("Encryption/Decryption failure: input and output aren't equal")
	}
	if debug {
		fmt.Println("in  ", string(*nti.PlainMessage))
		fmt.Println("out ", string(*nto.PlainMessage))
		fmt.Println("in  ", string((*nti.Key)[:]))
		fmt.Println("out ", string((*nto.Key)[:]))
	}
}

func TestBase58Encode(t *testing.T) {
	var nti *signalstor.NaclTool
	if nti = signalstor.NewEncryptionNaclTool([]byte(JSONText)); nti == nil {
		t.Fatalf(fmt.Sprintf("NewEncryptionNaclTool returned nil"))
	}
	{
		var encoded = nti.Key.Encode()
		var got = signalstor.DecodeKey(string(encoded))

		if debug {
			fmt.Println("nti ", string((*nti.Key)[:]))
			fmt.Println("got ", string(got[:]))
		}
		if !bytes.Equal(got[:], (*nti.Key)[:]) {
			t.Fatalf("Symmetric operations failed encode/decode")
		}
		if debug {
			fmt.Println("in  ", got.Encode())
			fmt.Println("out ", encoded)
		}
	}
	{
		var encoded = nti.Nonce.Encode()
		var got = signalstor.DecodeNonce(string(encoded))
		if debug {
			fmt.Println("nti  ", string((*nti.Nonce)[:]))
			fmt.Println("got ", string(got[:]))
		}
		if !bytes.Equal(got[:], (*nti.Nonce)[:]) {
			t.Fatalf("Symmetric operations failed encode/decode")
		}
		if debug {
			fmt.Println("in  ", got.Encode())
			fmt.Println("out ", encoded)
		}
	}

	{
		var ekey = nti.Key.Encode()
		var got = ekey.Decode()

		if debug {
			fmt.Println("nti ", string((*nti.Key)[:]))
			fmt.Println("got ", string(got[:]))
		}
		if !bytes.Equal(got[:], (*nti.Key)[:]) {
			t.Fatalf("Symmetric operations failed encode/decode")
		}
		if debug {
			fmt.Println("in  ", got.Encode())
			fmt.Println("out ", ekey)
		}
	}
	{
		var enonce = nti.Nonce.Encode()
		var got = enonce.Decode()
		if debug {
			fmt.Println("nti  ", string((*nti.Nonce)[:]))
			fmt.Println("got ", string(got[:]))
		}
		if !bytes.Equal(got[:], (*nti.Nonce)[:]) {
			t.Fatalf("Symmetric operations failed encode/decode")
		}
		if debug {
			fmt.Println("in  ", got.Encode())
			fmt.Println("out ", enonce)
		}
	}
}
