package signalstor // 	"github.com/davidwalter0/go-signalstor"

import (
	"strings"
	"fmt"
	"time"

	"github.com/davidwalter0/go-persist"
)

type Body string

func (body *Body) Escape() (text string) {
  if body == nil {
    panic("body is nil")
  }
  text = string(*body)
	text = strings.Replace(text, "'", "''", -1)
  text = strings.Replace(text, "%", "%%", -1)
  *body = Body(text)
  return 
}

// SmsDbIO object db I/O for sms table
type SmsDbIO struct {
	ID      int               `json:"id"`
	GUID    string            `json:"guid"`
	Created time.Time         `json:"created"`
	Changed time.Time         `json:"changed"`
	db      *persist.Database `ignore:"true"`
	Msg     SmsMessage
}

// SmsMessage sms message content
type SmsMessage struct {
	Address     string `json:"address"`
	Timestamp   string `json:"date"` // millisecond resolution sms timestamp
	ContactName string `json:"contact_name"`
	Date        string `json:"readable_date"`
	Subject     string `json:"subject"`
	Body        Body `json:"body"`
	Type        string `json:"type" doc:"1 received, 2 sent"`
}

// Ftp options to load from flags or env variables
type Ftp struct {
	Host     string
	Port     string
	User     string
	Password string
	Filename string
	Phone    string
	Key      string `json:"key"    doc:"encryption key for db write"`
	Debug    bool   `json:"debug"  doc:"dump configuration environment or flag parse result"`
}

var timeout uint = 5 // seconds
var writingToFile = false

// SmsMessageUnmarshal messages from xml parse are in a
// map[string]map[string]string. map["sms"] is the only entry in the
// parent map
type SmsMessageUnmarshal struct {
	SmsMessage `json:"sms"`
}
type SmsMessageArray []*SmsMessage

// SmsMessages array of mesages returned from full document parse
type SmsMessages struct {
	Messages []SmsMessage `json:"sms"`
}

type ByDateSmsMessageMap map[int]*SmsMessage
type UserTimestampMap map[string]ByDateSmsMessageMap

// BySMSTimestamp sortable type
type BySMSTimestamp SmsMessageArray

func (a BySMSTimestamp) Len() int           { return len(a) }
func (a BySMSTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySMSTimestamp) Less(i, j int) bool { return a[i].Timestamp < a[j].Timestamp }

// BySMSAddress sortable type
type BySMSAddress SmsMessageArray

func (a BySMSAddress) Len() int           { return len(a) }
func (a BySMSAddress) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySMSAddress) Less(i, j int) bool { return a[i].Address < a[j].Address }

func (message SmsMessage) String() string {
	contact := message.ContactName
	return fmt.Sprintf(
		`
Address   : %s
Timestamp : %s

Date      : %s
Contact   : %s
Message   : %s
`,
		message.Address,
		message.Timestamp,

		message.Date,
		contact,
		message.Body,
	)
}

// String formatted SmsDbIO text
func (message SmsDbIO) String() string {
	contact := message.Msg.ContactName
	return fmt.Sprintf(
		`
ID        : %12d
GUID      : %s
Address   : %s
Timestamp : %s

Date      : %s
Contact   : %s
Message   : %s
`,
		message.ID,
		message.GUID,
		message.Msg.Address,
		message.Msg.Timestamp,

		message.Msg.Date,
		contact,
		message.Msg.Body,
	)
}
