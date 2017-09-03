package main

import (
	"encoding/json"
	"os"
	// "encoding/xml"
	"bufio"
	"fmt"
	xj "github.com/basgys/goxml2json"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
)

type Ftp struct {
	Host     string
	Port     string
	User     string
	Password string
	Filename string
}

var timeout uint = 5 // seconds

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

type SMS map[string]SMSItem

type Writer struct {
	From         string
	ContactName  string
	Date         string
	ReadableDate string
	Address      string
	Subject      string
	Body         string
}

func (w *Writer) String() string {
	contact := w.ContactName
	if w.ContactName == "null" {
		contact = "*Unknown*"
	}

	// From    :%s
	return fmt.Sprintf(
		`Contact :%s
Date    :%s
Address :%s
Message :%s
`,
		// w.From,
		contact,
		w.ReadableDate,
		w.Address,
		// w.Subject,
		w.Body,
	)
}

func XmlJsonize(text string) {

	xml := strings.NewReader(text)
	var err error
	if err != nil {
		fmt.Println(fmt.Errorf("%s", err))
		os.Exit(1)
	}

	jsonx, err := xj.Convert(xml)
	if err != nil {
		panic("That's embarrassing...")
	}

	fmt.Println("typeof jsonx", reflect.TypeOf(jsonx))

	var temp map[string][]interface{}
	var x interface{}
	err = json.Unmarshal(jsonx.Bytes(), &temp)
	fmt.Println("typeof temp", reflect.TypeOf(temp))

	err = json.Unmarshal(jsonx.Bytes(), &x)
	fmt.Println("typeof x", reflect.TypeOf(x))

	if err != nil {
		fmt.Println(fmt.Errorf("%s", err))
		os.Exit(1)
	}

	fmt.Println(x)
	for k, v := range temp {
		fmt.Printf("Key : %s\n", k)
		var w Writer
		parsed := false
		for _, a := range v {
			for mk, mv := range a.(map[string]interface{}) {
				key := mk[1:]
				mvx := mv.(string)
				switch key {
				case "from":
					w.From = mvx
					parsed = true
				case "date":
					w.Date = mvx
					parsed = true
				case "readable_date":
					w.ReadableDate = mvx
					parsed = true
				case "address":
					w.Address = mvx
					parsed = true
				case "subject":
					w.Subject = mvx
					parsed = true
				case "contact_name":
					w.ContactName = mvx
					parsed = true
				case "body":
					w.Body = mvx
					parsed = true
				}
			}
			if parsed {
				fmt.Printf("%s", w.String())
				parsed = false
			}
		}
	}
}

func main() {
	filename := "SignalPlaintextBackup.xml"
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var xmlText string
	for scanner.Scan() {
		line := strings.TrimLeft(scanner.Text(), " ")
		prefix := line[0:5]
		switch prefix {
		case "<?xml":
		case "<smse":
			continue
		case "</sms":
			continue
		case "<sms ":
			eol := "> </sms>"
			endElem := "/>"
			if strings.Index(line, eol) == -1 {
				end := strings.Index(line, endElem)
				line = line[0:end] + eol
			}
		default:
			continue
		}

		xmlText += line + "\n"
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(xmlText)
	XmlJsonize(xmlText)
}
