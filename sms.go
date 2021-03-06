package signalstor // 	"github.com/davidwalter0/go-signalstor"

import (
	"database/sql"

	"fmt"
	"log"

	"github.com/davidwalter0/go-mutex"
	"github.com/davidwalter0/go-persist"
	"github.com/davidwalter0/go-persist/schema"
	"github.com/davidwalter0/go-persist/uuid"
)

var smsDB = &persist.Database{}
var standAlone = true
var dropAll = true
var smsDbInitialized = false
var monitor = mutex.NewMonitor()

// ConfigureDb alias for smsDbInitialize
func ConfigureDb() *persist.Database {
	smsDbInitialize()
	return smsDB
}

// ConfigureDb alias for smsDbInitialize
func (sms *SmsDbIO) ConfigureDb() *SmsDbIO {
	smsDbInitialize()
	return sms
}

func (sms *SmsDbIO) smsDbInitialize() {
	smsDB := sms.db
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
				// idempotent operation
				smsDB.Initialize(SmsDbIOSchema)
			}
		}
	}
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

// SmsDbIOSchema describes the table and triggers for persisting
// smsentications from totp objects from twofactor
// var SmsDbIOSchema schema.DBSchema = schema.DBSchema{
// timestamp is an sms millisecond time
// readable_date is mapped to date in the object,
var SmsDbIOSchema = schema.DBSchema{
	"sms": schema.SchemaText{ // timestamp <-> domain
		`CREATE TABLE sms (
       id  serial primary key,
       guid varchar(256) NOT NULL DEFAULT '' unique,
       address varchar(256) NOT NULL, 
       timestamp varchar(64) NOT NULL,  
       contact_name varchar(256) NOT NULL,
       readable_date varchar(256) NOT NULL, 
       subject varchar(64) DEFAULT '', 
       body text,
       type int NOT NULL default 1,
       encrypted boolean NOT NULL default false,
       created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
       changed timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
    )`,
		`CREATE UNIQUE INDEX unique_idx on sms (address, timestamp)`,
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

// NewKey create the key fields for an sms struct, notice that address
// uses account
func NewKey(address, timestamp string) *SmsDbIO {
	smsDbInitialize()
	return &SmsDbIO{
		Msg: SmsMessage{
			Address:   address,
			Timestamp: timestamp,
		},
		db: smsDB,
	}
}

// NewSmsDbIO initialize an sms struct
func NewSmsDbIO() *SmsDbIO {
	smsDbInitialize()
	return &SmsDbIO{
		db: smsDB,
	}
}

// NewSmsDbIOFromMsg initialize an sms struct
func NewSmsDbIOFromMsg(from *SmsMessage) *SmsDbIO {
	smsDbInitialize()
	return &SmsDbIO{
		Msg: *from,
		db:  smsDB,
	}
}

// Copy initialize an SmsDbIO struct from a message
func (sms *SmsMessage) Copy(from *SmsMessage) *SmsMessage {
	*sms = *from
	return sms
}

// CopySmsMessage initialize an SmsDbIO struct from a message
func (sms *SmsDbIO) CopySmsMessage(from *SmsMessage) *SmsDbIO {
	sms.Msg = *from
	return sms
}

// CopySmsDbIO initialize an sms struct
func (sms *SmsDbIO) CopySmsDbIO(from *SmsDbIO) *SmsDbIO {
	sms.ID = from.ID
	sms.GUID = from.GUID
	sms.Msg = from.Msg
	sms.Created = from.Created
	sms.Changed = from.Changed
	return sms
}

// Copy initialize an sms struct
func (sms *SmsDbIO) Copy(from *SmsDbIO) *SmsDbIO {
	return sms.CopySmsDbIO(from)
}

// CopyKey from SmsDbIO object
func (sms *SmsDbIO) CopyKey(from *SmsDbIO) *SmsDbIO {
	return sms.CopySmSDbIOKey(from)
}

// CopySmSDbIOKey initialize the sms's table key in the struct
func (sms *SmsDbIO) CopySmSDbIOKey(from *SmsDbIO) *SmsDbIO {
	sms.Msg.Address = from.Msg.Address
	sms.Msg.Timestamp = from.Msg.Timestamp
	return sms
}

// Create a row in a table
func (sms *SmsDbIO) Create() (err error) {
	if sms.db == nil {
		panic("SmsDbIO.db not initialized")
	}
	smsDB := sms.db
	if !sms.Msg.IsValid() || len(sms.Msg.Address) == 0 || len(sms.Msg.Timestamp) == 0 {
		text := fmt.Sprintf("Error message parse error broken or empty fields %v", sms.Msg)
		return fmt.Errorf(text)
	}

	// ignore DB & id
	insert := fmt.Sprintf(`
INSERT INTO sms 
(
  guid, 
  address,
  timestamp,
  subject,
  contact_name,
  body,
  readable_date,
  type,
  created,
  changed
)
VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		uuid.GUID().String(),
		sms.Msg.Address,
		sms.Msg.Timestamp,
		sms.Msg.Subject,
		sms.Msg.ContactName,
    sms.Msg.Body.Escape(),
		sms.Msg.Date,
		sms.Msg.Type,
	)
	_, err = smsDB.Exec(insert)
	return
}

// Read row from db using sms key fields for query
func (sms *SmsDbIO) Read() (err error) {
	if sms.db == nil {
		panic("SmsDbIO.db uninitialized")
	}
	smsDB := sms.db
	// ignore DB & id
	query := fmt.Sprintf(`
SELECT 
  id,
  guid, 
  address,
  timestamp,
  subject,
  contact_name,
  readable_date,
  body,
  type,
  created,
  changed
FROM
   sms 
WHERE
  address = '%s'
AND
  timestamp = '%s'
`,
		sms.Msg.Address,
		sms.Msg.Timestamp,
	)
	var rows *sql.Rows
	if rows, err = smsDB.Query(query); err != nil {
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			panic(err)
		}
	}()

	rows.Next()
	err = rows.Scan(
		&sms.ID,
		&sms.GUID,
		&sms.Msg.Address,
		&sms.Msg.Timestamp,
		&sms.Msg.Subject,
		&sms.Msg.ContactName,
		&sms.Msg.Date,
		&sms.Msg.Body,
		&sms.Msg.Type,
		&sms.Created,
		&sms.Changed)

	return
}

// Update row from db using sms key fields
func (sms *SmsDbIO) Update() (err error) {
	if sms.db == nil {
		panic("SmsDbIO.db uninitialized")
	}
	smsDB := sms.db
	// ignore DB & id
	update := fmt.Sprintf(`
UPDATE
  sms
SET
  subject       = '%s',
  contact_name  =  %s,
  readable_date = '%s',
  body          = '%s'
  type          = '%s'
WHERE
  address  = '%s'
AND
  timestamp     = '%s'
`,
		// set
		sms.Msg.Subject,
		sms.Msg.ContactName,
		sms.Msg.Date,
    sms.Msg.Body.Escape(),
		sms.Msg.Type,
		// where
		sms.Msg.Address,
		sms.Msg.Timestamp,
	)
	_, err = smsDB.Exec(update)
	return
}

// Delete row from db using sms key fields
func (sms *SmsDbIO) Delete() (err error) {
	if sms.db == nil {
		panic("SmsDbIO.db uninitialized")
	}
	smsDB := sms.db
	// ignore DB & id
	delete := fmt.Sprintf(`
DELETE FROM
  sms
WHERE
  address  = '%s'
AND
  timestamp     = '%s'
`,
		// where
		sms.Msg.Address,
		sms.Msg.Timestamp,
	)
	_, err = smsDB.Exec(delete)
	return
}

// Count rows for keys in sms
func (sms *SmsDbIO) Count() (count int) {
	if sms.db == nil {
		panic("SmsDbIO.db uninitialized")
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
  timestamp     = '%s'
`,
		// where
		sms.Msg.Address,
		sms.Msg.Timestamp,
	)

	row := smsDB.QueryRow(query)
	err := row.Scan(&count)
	if err != nil {
		log.Println("Row count query error", err)
	}
	return count
}

// ReadAll rows queried from db using sms key fields for query
func (sms *SmsDbIO) ReadAll(publish chan *SmsMessage) (err error) {
	if sms.db == nil {
		panic("SmsDbIO.db uninitialized")
	}
	smsDB := sms.db
	// ignore DB & id
	query := fmt.Sprintf(`
SELECT 
  id,
  guid, 
  address,
  timestamp,
  subject,
  contact_name,
  readable_date,
  body,
  type,
  created,
  changed
FROM
   sms 
WHERE
  address = '%s'
`,
		// where
		sms.Msg.Address,
	)
	var rows *sql.Rows
	if rows, err = smsDB.Query(query); err != nil {
		return
	}
	defer func() {
		if err = rows.Close(); err != nil {
			panic(err)
		}
	}()

	for rows.Next() {
		if err = rows.Scan(
			&sms.ID,
			&sms.GUID,
			&sms.Msg.Address,
			&sms.Msg.Timestamp,
			&sms.Msg.Subject,
			&sms.Msg.ContactName,
			&sms.Msg.Date,
			&sms.Msg.Body,
			&sms.Msg.Type,
			&sms.Created,
			&sms.Changed); err != nil {
			return err
		}
		var row = &SmsDbIO{}
		row.CopySmsDbIO(sms)
		publish <- &row.Msg
	}
	return err
}

// IsValid test for non empty records
func (sms *SmsMessage) IsValid() bool {
	return len(sms.Address) > 0 &&
		len(sms.Timestamp) > 0 &&
		// len(sms.ContactName) > 0 &&
		len(sms.Date) > 0 &&
		// len(sms.Subject) > 0 &&
		// len(sms.Body) > 0 &&
		len(sms.Type) > 0
}
