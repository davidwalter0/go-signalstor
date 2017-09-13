package signalstor

import (
	"fmt"
	"log"
	"os"

	"github.com/davidwalter0/db/uuid"
	"github.com/davidwalter0/go-ftp"
)

func ftpDeferableClose(connection *ftp.Connection) {
	err := connection.Logout()
	if err != nil {
		panic(err)
	}
}

// Download to file if writingToFile and use configured filename
// returning and empty buffer. If writingToFile is true return a
// buffer
func Download(cfg *Ftp, writingToFile bool) ([]byte, error) {
	var result []byte
	var err error

	var connection *ftp.Connection
	var connectString = fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	connection, err = ftp.Dial(connectString)
	defer ftpDeferableClose(connection)
	if err != nil {
		fmt.Printf("error: connection error %v\n", err)
		os.Exit(-1)
	}

	// login
	err = connection.Login(cfg.User, cfg.Password)

	if err != nil {
		log.Fatalf("error: login failure %v\n", err)
		os.Exit(-1)
	}

	if writingToFile {
		err = connection.Get(cfg.Filename, fmt.Sprintf("%s-%s", cfg.Filename, uuid.GUID()), ftp.BINARY, timeout)
		if err != nil {
			fmt.Printf("error: download failed %v\n", err)
			os.Exit(-1)
		}
	} else {
		result, err = connection.GetBuffer(cfg.Filename, ftp.BINARY, timeout)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	return result, err
}
