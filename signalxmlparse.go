package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/davidwalter0/db/uuid"
	"github.com/davidwalter0/go-cfg"
	"github.com/davidwalter0/go-ftp"
)

type IntSMSMap map[int]string
type UserDateMap map[string]IntSMSMap

type BySMSAddress []SMSItem

func (a BySMSAddress) Len() int           { return len(a) }
func (a BySMSAddress) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySMSAddress) Less(i, j int) bool { return a[i].Address < a[j].Address }

var userDateMap UserDateMap = make(UserDateMap)

// Ftp options to load from flags or env variables
type Ftp struct {
	Host     string
	Port     string
	User     string
	Password string
	Filename string
	Debug    bool `help:"dump configuration environment or flag parse result\n\t"`
}

var timeout uint = 5 // seconds
var writingToFile = false

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
		fmt.Errorf("%s\n", err)
	}

	jsonx, err := xj.Convert(xml)
	if err != nil {
		panic("That's embarrassing...")
	}

	if false {
		os.Exit(0)
	}

	var msinterface map[string][]interface{}
	err = json.Unmarshal(jsonx.Bytes(), &msinterface)

	if err != nil {
		_ = fmt.Errorf("%s\n", err)
	}

	for _, v := range msinterface {
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
				parsed = false
				if _, ok := userDateMap[w.Address]; !ok {
					userDateMap[w.Address] = make(IntSMSMap)
				}
				date, _ := strconv.Atoi(w.Date)
				if _, ok := userDateMap[w.Address][date]; ok {
					userDateMap[w.Address][date] += "\n" + w.Body
				} else {
					userDateMap[w.Address][date] = w.Body
				}
			}
		}
	}
	fmt.Println("==========================================================")
	names := make([]string, 0, len(userDateMap))
	for k, _ := range userDateMap {
		fmt.Println(k)
	}
	for name, _ := range userDateMap {
		names = append(names, name)
	}
	sort.Strings(names) //sort by key
	fmt.Println(names)
	fmt.Println("==========================================================")

	for _, name := range names {
		fmt.Println(name)
		byDate := userDateMap[name]
		dates := make([]int, 0, len(byDate))
		fmt.Println("==========================================================")
		for date := range byDate {
			dates = append(dates, date)
		}
		sort.Ints(dates) //sort by key
		for _, date := range dates {
			fmt.Println(date, ">", byDate[date])
		}
	}
}

var ftpcfg = &Ftp{}

func download() ([]byte, error) {
	var err error

	var connection *ftp.Connection
	var connectString = fmt.Sprintf("%s:%s", ftpcfg.Host, ftpcfg.Port)
	connection, err = ftp.Dial(connectString)
	defer connection.Logout()
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

	if writingToFile {
		err = connection.Get(ftpcfg.Filename, fmt.Sprintf("%s-%s", ftpcfg.Filename, uuid.GUID()), ftp.BINARY, timeout)
		if err != nil {
			fmt.Printf("error: download failed %v\n", err)
			os.Exit(-1)
		}
	}

	var result []byte
	result, err = connection.GetBuffer(ftpcfg.Filename, ftp.BINARY, timeout)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	return result, err
}

func xmlCleanup(xmlData []byte) (xmlText string) {
	scanner := bufio.NewScanner(bytes.NewReader(xmlData))
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
	return
}

func configure() {
	cfg.Parse(ftpcfg)
	if ftpcfg.Debug {
		fmt.Printf("%v %T\n", *ftpcfg, *ftpcfg)
		jsonText, _ := json.MarshalIndent(ftpcfg, "", "  ")
		fmt.Printf("\n%v\n", string(jsonText))
	}
}

func main() {
	fmt.Fprintln(os.Stderr, ftpcfg.Filename)
	configure()
	var err error
	var xmlData []byte

	xmlData, err = ioutil.ReadFile(ftpcfg.Filename)

	if err != nil {
		fmt.Printf("error: download failed %v\n", err)
		os.Exit(-1)
	}
	xmlText := xmlCleanup(xmlData)
	XmlJsonize(xmlText)
}
