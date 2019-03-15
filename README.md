---
This is test software experimenting with:

- xml - json conversion using an xml library
- ftp transfer tests
- download signal formatted xml backup
- load to database
- client side encrypt before write to database

*New*
- Update to work with newer xml output format of golang
  `github.com/xeals/signal-back` application
- Update example db.environment script to test the environment setup,
  assume that user has create db and user permissions via psql and the
  user, password and db name are sms, update for your environment.

The goal is to demonstrate transformation and secure backup of some
source data using whisper systems signal xml backup format as an
exmaple.

---

```
go get github.com/davidwalter0/go-signalstor

cd ${GOPATH}/src/github.com/davidwalter0/go-signalstor/examples

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
Commands

- ftpbufferget.go
  read a file into a buffer from ftp
- ftpfileget.go
  download a file
- signalxmlparsechan.go
  parse an input file using channels for unbuffered i/o
- signaldbchanreadall.go
  read all data previously loaded into a database table using channels
- signalxmlparsechanwriteall.go
  parse xml from file using channels between functions and write to
  database
- signalxmlparsedbwriteencrypted.go
  encrypt non-key data and write to database, ignore single byte data
  with little variance to limit simplifying known inputs for brute
  force attacks
- signalxmlparsedbwrite.go
  write data to database
- signalxmlparse.go
  parse xml and dump to stderr

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
