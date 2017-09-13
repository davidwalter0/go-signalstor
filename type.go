package signalstor

import (
	"fmt"
	"time"

	"github.com/davidwalter0/go-persist"
)

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
	Body        string `json:"body"`
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
	Debug    bool `help:"dump configuration environment or flag parse result\n\t"`
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
