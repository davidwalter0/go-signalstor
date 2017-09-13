package signalstor

import (
	"testing"
	"time"

	"github.com/davidwalter0/go-signalstor"
)

var key = []byte("LKHlhb899Y09olUi")

// pwgen --secure --no-capitalize --numerals 16 1
func Test_Encrypt(t *testing.T) {

	var encryptSmsDbIO = &signalstor.SmsDbIO{
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

	if err := encryptSmsDbIO.Encrypt([]byte(key)); err != nil {
		t.Fatalf("Decrypt failed %v", err)
	}

	var wanted = &signalstor.SmsDbIO{
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

	if err := encryptSmsDbIO.Decrypt([]byte(key)); err != nil {
		t.Fatalf("Decrypt failed %v", err)
	}

	got := encryptSmsDbIO

	if *wanted != *got {
		t.Fatalf("object failed to decrypt wanted\n%v\ngot\n%v\n", wanted, got)
	}
}
