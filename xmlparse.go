package xml2json

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	// "github.com/davidwalter0/db/uuid"
	// "github.com/davidwalter0/go-cfg"
	// "github.com/davidwalter0/go-ftp"
	"github.com/davidwalter0/go-xml2json"
)

const (
	CharSixTeen = "\x00"
	CharZero    = "\x16"
)

// Byte2Struct xml and update a struct from the object
// func Byte2Struct(buffer *bytes.Buffer) {

// 	var v map[string]map[string]string
// 	err := json.Unmarshal(buffer.Bytes(), &v)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	parsed := false
// 	if a, ok := v["sms"]; ok {
// 		var writer SmsMessageMap
// 		for key, value := range a {
// 			// fmt.Println(key, reflect.Name())
// 			switch key {
// 			case "date":
// 				writer.Date = value
// 				parsed = true
// 			case "readable_date":
// 				writer.ReadableDate = value
// 				parsed = true
// 			case "address":
// 				writer.Address = value
// 				parsed = true
// 			case "subject":
// 				writer.Subject = value
// 				parsed = true
// 			case "contact_name":
// 				writer.ContactName = value
// 				parsed = true
// 			case "body":
// 				writer.Body = value
// 				parsed = true
// 			}
// 		}
// 		if parsed {
// 			parsed = false
// 			if _, ok := userDateMap[writer.Address]; !ok {
// 				userDateMap[writer.Address] = make(IntSMSMap)
// 			}
// 			date, _ := strconv.Atoi(writer.Date)
// 			if _, ok := userDateMap[writer.Address][date]; ok {
// 				userDateMap[writer.Address][date].Body += "\n" + writer.Body
// 			} else {
// 				userDateMap[writer.Address][date] = &writer
// 			}
// 		}
// 	}
// }

// Byte2Struct byte to struct conversion
func Byte2Struct(buffer *bytes.Buffer, dest interface{}) error {
	return json.Unmarshal(buffer.Bytes(), dest)
	// if err != nil {
	// 	log.Fatalf("json.Unmarshal failure: %v\n", err)
	// }
	// type Parseable struct {
	// 	SmsMessageMap `json:"sms"`
	// }
	// var sms Parseable
	// err := json.Unmarshal(buffer.Bytes(), &sms)
	// if err != nil {
	// 	log.Fatalf("error: xml element to json parse %v\n", err)
	// }
	// if _, ok := userDateMap[sms.SmsMessageMap.Address]; !ok {
	// 	userDateMap[sms.SmsMessageMap.Address] = make(IntSMSMap)
	// }
	// date, _ := strconv.Atoi(sms.SmsMessageMap.Date)
	// if _, ok := userDateMap[sms.SmsMessageMap.Address][date]; ok {
	// 	userDateMap[sms.SmsMessageMap.Address][date].Body += "\n" + sms.SmsMessageMap.Body
	// } else {
	// 	userDateMap[sms.SmsMessageMap.Address][date] = &sms.SmsMessageMap
	// }

}

// // Dump the data from the xml parsed file to stdout
// func Dump(file io.Writer) {
// 	addresses := make([]string, 0, len(userDateMap))
// 	for address := range userDateMap {
// 		addresses = append(addresses, address)
// 	}
// 	sort.Strings(addresses)
// 	fmt.Fprintln(file, "==========================================================")

// 	for _, address := range addresses {
// 		fmt.Fprintln(file, address)
// 		byDate := userDateMap[address]
// 		bySMSDates := make(BySMSDate, 0, len(byDate))
// 		fmt.Fprintln(file, "==========================================================")
// 		for _, smsMessage := range byDate {
// 			bySMSDates = append(bySMSDates, smsMessage)
// 		}
// 		sort.Sort(bySMSDates) //sort by key
// 		for _, smsMessage := range bySMSDates {
// 			fmt.Fprintln(file, smsMessage)
// 		}
// 	}
// }

// DumpParsedMessages the data from the xml parsed file to stdout
func DumpParsedMessages(file io.Writer, messages SmsMessages) {
	var byUserSMS map[string]map[int]*SmsMessage = make(map[string]map[int]*SmsMessage)

	for _, sms := range messages.Messages {
		tmpsms := sms
		date, _ := strconv.Atoi(sms.Date)
		address := sms.Address
		if _, ok := byUserSMS[address]; !ok {
			byUserSMS[address] = make(map[int]*SmsMessage)
		}
		if _, ok := byUserSMS[address][date]; !ok {
			byUserSMS[address][date] = &tmpsms
		} else {
			byUserSMS[address][date].Body += sms.Body
		}

	}
	addresses := make([]string, 0, len(byUserSMS))
	for address := range byUserSMS {
		addresses = append(addresses, address)
	}

	sort.Strings(addresses)

	for _, address := range addresses {
		fmt.Fprintln(file, "==========================================================")
		fmt.Println("Address", address)
		fmt.Fprintln(file, "==========================================================")

		dates := make([]int, 0, len(byUserSMS[address]))
		for date := range byUserSMS[address] {
			dates = append(dates, date)
		}
		sort.Ints(dates)
		for _, date := range dates {
			fmt.Fprintln(file, byUserSMS[address][date])
		}
	}
}

// NoOp do nothing placeholder handler
func NoOp(st interface{}) {
}

// SmsXMLFixUp preprocessor for xml string
func SmsXMLFixUp(xml string) (string, error) {
	line := strings.TrimLeft(xml, " ")
	prefix := line[0:5]
	// fmt.Printf("SmsXMLFixUp [%s] %s\n", prefix, line)
	switch prefix {
	case "<?xml", "<smse", "</sms":
		return "", nil
	case "<sms ":
		eol := "> </sms>"
		endElem := "/>"
		if strings.Index(line, eol) == -1 {
			end := strings.Index(line, endElem)
			line = line[0:end] + eol
		}
	default:
		return "", nil
	}
	return line, nil
}

type XMLFixUp func(xml string) (string, error)
type StructHandler func(st interface{})

// XMLParse assumes line oriented xml elements raw data is an xml
// document, dest is a struct unmarshal target, xmlFixup can be a pass
// through noop function returning the string sent to it, or a more
// complex fix function for broken xml, handler is a post process call
// which might write to a standard location, persistent store or a
// noop.
func XMLParse(rawData []byte,
	dest *SmsMessages,
	xmlFixUp XMLFixUp,
	handler StructHandler) (xmlText string) {

	var buffer *bytes.Buffer
	var err error
	var line string

	scanner := bufio.NewScanner(bytes.NewReader(rawData))
	for scanner.Scan() {
		line, err = xmlFixUp(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		if len(line) > 0 {
			var text = line
			xml := strings.NewReader(text)
			buffer, err = xml2json.Convert(xml)
			if err != nil {
				log.Fatalf("error: xml element to json parse %v\n", err)
			}
			var buffer2, err = xml2json.Convert(xml)

			if len(buffer.Bytes()) > 0 {
				var message SmsMessageUnmarshal

				if err = json.Unmarshal(buffer.Bytes(), &message); err != nil {
					fmt.Fprintf(os.Stderr, "error: %v\n", err)
					var i interface{}
					if err = json.Unmarshal(buffer2.Bytes(), &i); err != nil {
						fmt.Fprintf(os.Stderr, "error: %v\n", err)
					} else {
						fmt.Fprintf(os.Stderr, ">> %T %v\n", i, i)
						switch i.(type) {
						case map[string]interface{}:
							for k, v := range i.(map[string]interface{}) {
								fmt.Fprintf(os.Stderr, "i.(map[string]interface{}) %s %v %T\n", k, v, v)
							}
						case string:
							fmt.Fprintf(os.Stderr, "string %v\n", i)
						}
					}
					continue
				}
				handler(dest)
				dest.Messages = append(dest.Messages, message.SmsMessage)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

// XMLParseArray parse an xml document with multiple elements
func XMLParseArray(rawData []byte, dest interface{},
	xmlFixUp XMLFixUp,
	handler StructHandler) {

	var buffer *bytes.Buffer
	var err error
	var line string
	var xmlFixed string
	// fmt.Fprintf(os.Stderr, "%v", string(rawData))
	scanner := bufio.NewScanner(bytes.NewReader(rawData))
	for scanner.Scan() {
		line, err = xmlFixUp(scanner.Text())
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}
		if len(line) > 0 {
			xmlFixed += line
		}
	}

	xml := bytes.NewReader([]byte(xmlFixed))

	buffer, err = xml2json.Convert(xml)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	// parser error : xmlParseCharRef: invalid xmlChar value 16
	// parser error : xmlParseCharRef: invalid xmlChar value 0
	var text = buffer.String()
	// text = strings.Replace(buffer.String(), fmt.Sprintf("%c", 16), " ", -1)
	// text = strings.Replace(text, fmt.Sprintf("%c", 0), " ", -1)
	// fmt.Println(text)
	if err = json.Unmarshal([]byte(text), dest); err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
	}
	handler(dest)
}
