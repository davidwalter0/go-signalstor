package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/davidwalter0/db/uuid"
	"github.com/davidwalter0/go-cfg"
	"github.com/davidwalter0/go-ftp"
)

type Ftp struct {
	Host     string
	Port     string
	User     string
	Password string
	Filename string
	Debug    bool
}

var timeout uint = 5 // seconds

// Connect to the server using go-cfg setup
func Connect() {
	var err error
	var ftpcfg = &Ftp{}

	cfg.Parse(ftpcfg)
	if ftpcfg.Debug {
		fmt.Printf("%v %T\n", *ftpcfg, *ftpcfg)
		jsonText, _ := json.MarshalIndent(ftpcfg, "", "  ")
		fmt.Printf("\n%v\n", string(jsonText))
	}

	var connection *ftp.Connection
	var connectString = fmt.Sprintf("%s:%s", ftpcfg.Host, ftpcfg.Port)
	connection, err = ftp.Dial(connectString)
	if err != nil {
		fmt.Printf("error: connection error %v\n", err)
		os.Exit(-1)
	}

	// login
	err = connection.Login(ftpcfg.User, ftpcfg.Password)
	if err != nil {
		fmt.Printf("error: login failure %v\n", err)
		os.Exit(-1)
	}

	var code uint
	var response string
	code, response, err = connection.Cmd("list", ftp.BINARY)
	if err != nil {
		fmt.Printf("error: list command failed %v\n", err)
		fmt.Printf("code %d response %s\n", code, response)
		os.Exit(-1)
	}

	err = connection.Get(ftpcfg.Filename, fmt.Sprintf("%s-%s", ftpcfg.Filename, uuid.GUID()), ftp.BINARY, timeout)
	if err != nil {
		fmt.Printf("error: download failed %v\n", err)
		os.Exit(-1)
	}
	connection.Logout()
	return
}

func main() {
	Connect()
	os.Exit(0)
}
