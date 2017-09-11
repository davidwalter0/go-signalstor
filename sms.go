package xml2json

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/davidwalter0/go-mutex"
	"github.com/davidwalter0/go-persist"
	"github.com/davidwalter0/go-persist/schema"
	"github.com/davidwalter0/go-persist/uuid"
)

// SmsDb db connection and i/o interface object
type SmsDb *persist.Database

var smsDB = &persist.Database{}
var standAlone = true
var dropAll = true
var smsDbInitialized = false
var monitor = mutex.NewMonitor()

// ConfigureDb alias for smsDbInitialize
func ConfigureDb() SmsDb {
	smsDbInitialize()
	return smsDB
}

// Initialize a database connection
func smsDbInitialize() {
	if !smsDbInitialized {
		defer monitor()()
		if !smsDbInitialized {
			smsDbInitialized = true
			if standAlone {
				smsDB.ConfigEnvWPrefix("SQL", false)
				smsDB.Connect()
				if dropAll {
					smsDB.DropAll(SmsDbIOSchema)
				}
				smsDB.Initialize(SmsDbIOSchema)
			}
		}
	}
}

// func Initialize() {
// 	init()
// }

// SmsDbIOSchema describes the table and triggers for persisting
// smsentications from totp objects from twofactor
// var SmsDbIOSchema schema.DBSchema = schema.DBSchema{
var SmsDbIOSchema = schema.DBSchema{
	"sms": schema.SchemaText{ // date <-> domain
		`CREATE TABLE sms (
       id  serial primary key,
       guid varchar(256) NOT NULL DEFAULT '' unique,
       address varchar(256) NOT NULL, 
       date varchar(32) NOT NULL, 
       contact_name varchar(256) NOT NULL,
       readable_date varchar(256) NOT NULL, 
       subject varchar(32) DEFAULT '', 
       body text,
       type int NOT NULL default 1,
       created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
       changed timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
    )`,
		`CREATE UNIQUE INDEX unique_idx on sms (address, date)`,
		`CREATE OR REPLACE FUNCTION update_created_column()
       RETURNS TRIGGER AS $$
       BEGIN
          NEW.changed = now(); 
          RETURN NEW;
       END;
       $$ language 'plpgsql'`,
		`CREATE TRIGGER update_ab_changetimestamp 
       BEFORE UPDATE ON sms 
       FOR EACH ROW EXECUTE PROCEDURE update_created_column()`,
	},
}

// SmsDbIOKey accessible object for database smsentication table I/O
// type SmsDbIOKey struct {
// 	Address string `json:"address"`
// 	Date    string `json:"date"`
// }

// SmsDbIO object db I/O for sms table
type SmsDbIO struct {
	ID      int               `json:"id"`
	GUID    string            `json:"guid"`
	Created time.Time         `json:"created"`
	Changed time.Time         `json:"changed"`
	db      *persist.Database `ignore:"true"`
	Msg     SmsMessage
}

// NewKey create the key fields for an sms struct, notice that address
// uses account
func NewKey(address, date string) *SmsDbIO {
	smsDbInitialize()
	return &SmsDbIO{
		Msg: SmsMessage{
			Address: address,
			Date:    date,
		},
		db: smsDB,
	}
}

// NewSmsDbIO smsDbInitialize an sms struct
func NewSmsDbIO() *SmsDbIO {
	smsDbInitialize()
	return &SmsDbIO{
		db: smsDB,
	}
}

// NewSmsDbIOFromMsg smsDbInitialize an sms struct
func NewSmsDbIOFromMsg(from *SmsMessage) *SmsDbIO {
	smsDbInitialize()
	return &SmsDbIO{
		Msg: *from,
		db:  smsDB,
	}
}

// // NewSmsDbIO smsDbInitialize an sms struct
// func NewSmsDbIO(address, date, subject, key, body, contactName, readableDate string) *SmsDbIO {
// 	return &SmsDbIO{
// 		Msg: SmsMessage{
// 			Address:      address,
// 			Date:         date,
// 			Subject:      subject,
// 			ContactName:  contactName,
// 			Body:         body,
// 			ReadableDate: readableDate,
// 		},
// 		db: smsDB,
// 	}
// }

// CopySmsMessage smsDbInitialize an SmsDbIO struct from a message
func (sms *SmsDbIO) CopySmsMessage(from *SmsMessage) {
	sms.Msg = *from
}

// CopySmsDbIO smsDbInitialize an sms struct
func (sms *SmsDbIO) CopySmsDbIO(from *SmsDbIO) {
	sms.ID = from.ID
	sms.GUID = from.GUID
	sms.Msg = from.Msg
	sms.Created = from.Created
	sms.Changed = from.Changed
}

// CopyKey smsDbInitialize the sms's table key in the struct
func (sms *SmsDbIO) CopySmSDbIOKey(from *SmsDbIO) {
	sms.Msg.Address = from.Msg.Address
	sms.Msg.Date = from.Msg.Date
}

// Create a row in a table
func (sms *SmsDbIO) Create() {
	if sms.db == nil {
		panic("SmsDbIO.db unsmsDbInitialized")
	}
	smsDB := sms.db
	// ignore DB & id
	insert := fmt.Sprintf(`
INSERT INTO sms 
(
  guid, 
  address,
  date,
  subject,
  contact_name,
  body,
  readable_date,
  created,
  changed
)
VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		uuid.GUID().String(),
		sms.Msg.Address,
		sms.Msg.Date,
		sms.Msg.Subject,
		sms.Msg.ContactName,
		sms.Msg.Body,
		sms.Msg.ReadableDate,
	)
	// fmt.Println(insert)
	// fmt.Println(smsDB.Exec(insert))
	_, err := smsDB.Exec(insert)
	if err != nil {
		log.Println("Row count query error", err)
	}
	// fmt.Println("Count", sms.Count())
}

// Read row from db using sms key fields for query
func (sms *SmsDbIO) Read() bool {
	if sms.db == nil {
		panic("SmsDbIO.db unsmsDbInitialized")
	}
	smsDB := sms.db
	// ignore DB & id
	query := fmt.Sprintf(`
SELECT 
  id,
  guid, 
  address,
  date,
  subject,
  contact_name,
  readable_date,
  body,
  created,
  changed
FROM
   sms 
WHERE
  address = '%s'
AND
  date = '%s'
`,
		sms.Msg.Address,
		sms.Msg.Date,
	)
	fmt.Println(query)
	rows := smsDB.Query(query)
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&sms.ID,
			&sms.GUID,
			&sms.Msg.Address,
			&sms.Msg.Date,
			&sms.Msg.Subject,
			&sms.Msg.ContactName,
			&sms.Msg.ReadableDate,
			&sms.Msg.Body,
			&sms.Created,
			&sms.Changed); err != nil {
			panic(fmt.Sprintf("%v", err))
		}
		fmt.Println(
			sms.ID,
			sms.GUID,
			sms.Msg.Address,
			sms.Msg.Date,
			sms.Msg.Subject,
			sms.Msg.ContactName,
			sms.Msg.ReadableDate,
			sms.Msg.Body,
			sms.Created,
			sms.Changed)
	}
	count := sms.Count()
	// fmt.Println("Count", count)
	return count != 0
}

// Update row from db using sms key fields
func (sms *SmsDbIO) Update() {
	if sms.db == nil {
		panic("SmsDbIO.db unsmsDbInitialized")
	}
	smsDB := sms.db
	// ignore DB & id
	update := fmt.Sprintf(`
UPDATE
  sms
SET
  subject       = '%s',
  contact_name  =  %d,
  readable_date = '%s',
  body          = '%s'
WHERE
  address  = '%s'
AND
  date     = '%s'
`,
		// set
		sms.Msg.Subject,
		sms.Msg.ContactName,
		sms.Msg.ReadableDate,
		sms.Msg.Body,
		// where
		sms.Msg.Address,
		sms.Msg.Date,
	)
	var err error
	var rows *sql.Rows
	var result sql.Result
	fmt.Println(update)
	result, err = smsDB.Exec(update)
	fmt.Println("update result", result, "error", err)
	rows = smsDB.Query("SELECT * FROM sms")
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&sms.ID,
			&sms.GUID,
			&sms.Msg.Address,
			&sms.Msg.Date,
			&sms.Msg.Subject,
			&sms.Msg.ContactName,
			&sms.Msg.ReadableDate,
			&sms.Msg.Body,
			&sms.Created,
			&sms.Changed); err != nil {
			panic(fmt.Sprintf("%v", err))
		}
		fmt.Println(
			sms.ID,
			sms.GUID,
			sms.Msg.Address,
			sms.Msg.Date,
			sms.Msg.Subject,
			sms.Msg.ContactName,
			sms.Msg.ReadableDate,
			sms.Msg.Body,
			sms.Created,
			sms.Changed)
	}
	// fmt.Println("Count", sms.Count())
}

// Delete row from db using sms key fields
func (sms *SmsDbIO) Delete() {
	if sms.db == nil {
		panic("SmsDbIO.db unsmsDbInitialized")
	}
	smsDB := sms.db
	// ignore DB & id
	delete := fmt.Sprintf(`
DELETE FROM
  sms
WHERE
  address  = '%s'
AND
  date     = '%s'
`,
		// where
		sms.Msg.Address,
		sms.Msg.Date,
	)
	var err error
	var rows *sql.Rows
	var result sql.Result
	fmt.Println(delete)
	result, err = smsDB.Exec(delete)
	fmt.Println("delete result", result, "error", err)
	rows = smsDB.Query("SELECT * FROM sms")
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&sms.ID,
			&sms.GUID,
			&sms.Msg.Address,
			&sms.Msg.Date,
			&sms.Msg.Subject,
			&sms.Msg.ContactName,
			&sms.Msg.ReadableDate,
			&sms.Msg.Body,
			&sms.Created,
			&sms.Changed); err != nil {
			panic(fmt.Sprintf("%v", err))
		}
		fmt.Println(
			sms.ID,
			sms.GUID,
			sms.Msg.Address,
			sms.Msg.Date,
			sms.Msg.Subject,
			sms.Msg.ContactName,
			sms.Msg.ReadableDate,
			sms.Msg.Body,
			sms.Created,
			sms.Changed)
	}
	// fmt.Println("Count", sms.Count())
}

// Count rows for keys in sms
func (sms *SmsDbIO) Count() (count int) {
	if sms.db == nil {
		panic("SmsDbIO.db unsmsDbInitialized")
	}
	smsDB := sms.db
	query := fmt.Sprintf(`
SELECT
  COUNT(*) 
FROM
  sms
WHERE 
  address  = '%s'
AND
  date     = '%s'
`,
		// where
		sms.Msg.Address,
		sms.Msg.Date,
	)

	row := smsDB.QueryRow(query)
	err := row.Scan(&count)
	if err != nil {
		log.Println("Row count query error", err)
	}
	return count
}
