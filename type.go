package xml2json

import (
	"fmt"
)

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

// type SMSItem struct {
// 	FmtDate       string `json:"readable_date"`
// 	ScToa         string `json:"sc_toa"`
// 	Read          string `json:"read"`
// 	Status        string `json:"status"`
// 	Protocol      string `json:"protocol"`
// 	ContactName   string `json:"contact_name"`
// 	Type          string `json:"type"`
// 	Subject       string `json:"subject"`
// 	Toa           string `json:"toa"`
// 	Date          string `json:"date"`
// 	Address       string `json:"address"`
// 	Body          string `json:"body"`
// 	ServiceCenter string `json:"service_center"`
// 	Locked        string `json:"locked"`
// }

// SmsMessage sms message content
type SmsMessage struct {
	ContactName  string `json:"contact_name"`
	Date         string `json:"date"`
	ReadableDate string `json:"readable_date"`
	Address      string `json:"address"`
	Subject      string `json:"subject"`
	Body         string `json:"body"`
	Type         string `json:"type" doc:"1 received, 2 sent"`
}

// type SMS map[string]SMSItem
type SmsMessageUnmarshal struct {
	SmsMessage `json:"sms"`
}
type SmsMessageArray []*SmsMessage

// SmsMessages array of mesages returned from full document parse
type SmsMessages struct {
	Messages []SmsMessage `json:"sms"`
}

type ByDateSmsMessageMap map[int]*SmsMessage
type UserDateMap map[string]ByDateSmsMessageMap

// var userDateMap UserDateMap = make(UserDateMap)

// BySMSDate sortable type
type BySMSDate SmsMessageArray

func (a BySMSDate) Len() int           { return len(a) }
func (a BySMSDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySMSDate) Less(i, j int) bool { return a[i].Date < a[j].Date }

// BySMSAddress sortable type
type BySMSAddress SmsMessageArray

func (a BySMSAddress) Len() int           { return len(a) }
func (a BySMSAddress) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySMSAddress) Less(i, j int) bool { return a[i].Address < a[j].Address }

// func (message *SmsMessageUnmarshal) String() string {
// 	contact := message.SmsMessage.ContactName
// 	return fmt.Sprintf(
// 		`Contact :%s
// Date    :%s
// Address :%s
// Message :%s
// `,
// 		// message.From,
// 		contact,
// 		message.SmsMessage.ReadableDate,
// 		message.SmsMessage.Address,
// 		// message.Subject,
// 		message.SmsMessage.Body,
// 	)
// }

func (message SmsMessage) String() string {
	contact := message.ContactName
	return fmt.Sprintf(
		`Contact : %s
Date    : %s
Address : %s
Message : %s
`,
		// message.From,
		contact,
		message.ReadableDate,
		message.Address,
		// message.Subject,
		message.Body,
	)
}
