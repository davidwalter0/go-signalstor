---
This is test software experimenting with:

- xml - json conversion using an xml library
- ftp transfer tests

Todo 

- database i/o using persist library
- encryption

The goal is to demonstrate transformation and secure backup of some
source data

---

```
go get github.com/davidwalter0/xml2json

cd ${GOPATH}/src/github.com/davidwalter0/xml2json/examples
```

- fetch a version of a signal backup xml file 

- edit an environment file, export the corresponding environment
  variables or set the corresponding command line flags

```
. ftp.environment
export FTP_FILENAME=/path/to/SignalPlaintextBackup.xml
go run signalxmlparse.go
```



---

Ignoring some fields and renaming date to timestamp and readable_date
to Date

```
    type SMSItem struct {
            FmtDate       string `json:"readable_date"`
            ScToa         string `json:"sc_toa"`
            Read          string `json:"read"`
            Status        string `json:"status"`
            Protocol      string `json:"protocol"`
            ContactName   string `json:"contact_name"`
            Type          string `json:"type"`
            Subject       string `json:"subject"`
            Toa           string `json:"toa"`
            Date          string `json:"date"`
            Address       string `json:"address"`
            Body          string `json:"body"`
            ServiceCenter string `json:"service_center"`
            Locked        string `json:"locked"`
    }

```

- SmsMessage sms message content
- timestamp is an sms message time to millisecond
- Date is the text version to seconds

```
type SmsMessage struct {
	Address     string `json:"address"`
	Timestamp   string `json:"date"` // millisecond resolution sms timestamp
	ContactName string `json:"contact_name"`
	Date        string `json:"readable_date"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	Type        string `json:"type" doc:"1 received, 2 sent"`
}
```