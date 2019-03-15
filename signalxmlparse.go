package signalstor // 	"github.com/davidwalter0/go-signalstor"

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

	xj "github.com/davidwalter0/go-xml2json"
)

const (
	CharSixTeen = "\x00"
	CharZero    = "\x16"
)

// Byte2Struct byte to struct conversion
func Byte2Struct(buffer *bytes.Buffer, dest interface{}) error {
	return json.Unmarshal(buffer.Bytes(), dest)
}

// DumpParsedMessages the data from the xml parsed file to stdout
func DumpParsedMessages(file io.Writer, messages SmsMessages) {
	var byUserSMS map[string]map[int]*SmsMessage = make(map[string]map[int]*SmsMessage)

	for _, sms := range messages.Messages {
		tmpsms := sms
		date, _ := strconv.Atoi(sms.Timestamp)
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
		fmt.Fprintln(file, "Address", address)
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
func NoOp(st interface{}) error {
	return nil
}

// SmsMessageValidate method, no empty messages
func SmsMessageValidate(msg interface{}) (err error) {
	switch message := msg.(type) {
	case SmsMessage:
		fmt.Printf("SmsMessage %v %T", message, message)
		if !message.IsValid() {
			err = fmt.Errorf("Invalid SmsMessage %v", message)
			fmt.Printf("SmsMessageValidate Error %v\n", err)
		}
	case SmsMessages:
		err = fmt.Errorf("SmsMessages %v %T", message, message)
	default:
		err = fmt.Errorf("Not SmsMessage %v %T", message, message)
	}
	return
}

// SmsXMLFixUp preprocessor for xml string
func SmsXMLFixUp(xml string) (string, error) {

	line := strings.TrimLeft(xml, " ")
	if len(line) < 5 {
		return "", nil
	}
	prefix := line[0:5]
	switch prefix {
	case "<?xml", "<smse", "</sms":
		return "", nil
	case "<sms ":
    var eol string
		var eol1 = "> </sms>"
    var eol2 = "></sms>"
    if strings.Index(line, eol1) >=0 {
      eol = eol1
    }
    if strings.Index(line, eol2) >=0 {
      eol = eol2
    }
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

// SmsXMLFixUpChan preprocessor for xml string
func SmsXMLFixUpChan(reader, writer chan string) {
	for xml := range reader {
		line := strings.TrimLeft(xml, " ")
		prefix := line[0:5]
		switch prefix {
		case "<?xml", "<smse", "</sms":
		case "<sms ":
			eol := "> </sms>"
			endElem := "/>"
			if strings.Index(line, eol) == -1 {
				end := strings.Index(line, endElem)
				line = line[0:end] + eol
			}
			writer <- line
		default:
			writer <- ""
		}
	}
}

// XMLFixUp repairs non-conforming xml for this import set
type XMLFixUp func(xml string) (string, error)

// ApplyStructRuleHandler cleanup handler
type ApplyStructRuleHandler func(st interface{}) error

// XMLParse assumes line oriented xml elements raw data is an xml
// document, dest is a struct unmarshal target, xmlFixup can be a pass
// through noop function returning the string sent to it, or a more
// complex fix function for broken xml, handler is a post process call
// which might write to a standard location, persistent store or a
// noop.
func XMLParse(rawData []byte,
	dest *SmsMessages,
	xmlFixUp XMLFixUp,
	handler ApplyStructRuleHandler) (xmlText string) {

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

			if buffer, err = xj.Convert(xml); err != nil {
				log.Printf("Error parsing buffer %s error %v\n", buffer.String(), err)
			}

			if err != nil {
				log.Fatalf("Error xml element to json parse %v\n", err)
			}

			if len(buffer.Bytes()) > 0 {
				var fixed []byte
				if fixed = KeyFix(buffer.Bytes()); fixed != nil {
					var message SmsMessageUnmarshal
					if err = json.Unmarshal(fixed, &message); err != nil {
						fmt.Fprintf(os.Stderr, "error: %v\n", err)
						continue
					}
					if err = handler(message.SmsMessage); err == nil {
						dest.Messages = append(dest.Messages, message.SmsMessage)
					}
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

// KeyFix repair of funkiness with json key being prefixed with a dash
// hyphen bytes in bytes out map representation
func KeyFix(raw []byte) []byte {
	var err error
	var iface interface{}

	if err = json.Unmarshal(raw, &iface); err != nil {
		panic(fmt.Sprintf("Error parsing (Unmarshaling) JSON: %v", err))
	}
	var m = keyFix(iface)
	if raw, err = json.Marshal(m); err != nil {
		panic(fmt.Sprintf("Error assembling (Marshaling) JSON: %v", err))
	}
	return raw
}

// KeyFix repair of funkiness with key being prefixed with a dash /
// hyphen
func keyFix(raw interface{}) (repaired map[string]map[string]string) {
	// JSON object parses into a map with string keys
	if kvMap, ok := raw.(map[string]interface{}); ok {
		// JSON object parses into a map with string keys
		if kvMap, ok = kvMap["sms"].(map[string]interface{}); ok {

			if len(kvMap) > 0 {
				repaired = make(map[string]map[string]string)
				repaired["sms"] = make(map[string]string)
			}
			for k, v := range kvMap {
				// hyphen from xml parse oddity
				if len(k) > 0 && string(k)[0] == '-' {
					repaired["sms"][k[1:]] = v.(string)
				} else {
					repaired["sms"][k] = v.(string)
				}
			}
		}
	}
	return
}

// xmlParseArray parse an xml document with multiple elements
// BUGS with invalid chars in xml conversion and others
func xmlParseArray(rawData []byte, dest interface{},
	xmlFixUp XMLFixUp,
	handler ApplyStructRuleHandler) error {

	var buffer *bytes.Buffer
	var err error
	var line string
	var xmlFixed string

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

	buffer, err = xj.Convert(xml)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	if err = json.Unmarshal(buffer.Bytes(), dest); err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
	}
	return handler(dest)
}

////////////////////////////////////////////////////////////////////////
// Line oriented implementation
////////////////////////////////////////////////////////////////////////

// GoChannelLineReader open, read, write line by line, run as go func
func GoChannelLineReader(filename string, writer chan *string) {
	var err error
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer close(writer)
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		writer <- &line
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// XMLParsePublish assumes line oriented xml elements raw data is an xml
// document, dest is a struct unmarshal target
func XMLParsePublish(filename string, messages chan *SmsMessage) {
	var xmlFixUp = SmsXMLFixUp
	var buffer *bytes.Buffer
	var err error
	var line string
	var scanner = make(chan *string)
	defer close(messages)

	go GoChannelLineReader(filename, scanner)

	for next := range scanner {
		line, err = xmlFixUp(*next)
		if err != nil {
			log.Fatal(err)
		}
		if len(line) > 0 {
			var text = line
			xml := strings.NewReader(text)
			buffer, err = xj.Convert(xml)
			if err != nil {
				continue
			}
			if len(buffer.Bytes()) > 0 {
				var message SmsMessageUnmarshal
				if err = json.Unmarshal(buffer.Bytes(), &message); err != nil {
					continue
				}
				// unbuffered channel to downstream
				messages <- &message.SmsMessage
			}
		}
	}
}

// DumpParsedMessagesSubscribe the data from the xml parsed file to
// file io.Writer reading objects one at a time from messages channel
// order by address then by date
func DumpParsedMessagesSubscribe(file io.Writer, messages chan *SmsMessage, done chan bool) {

	var byUserSMS map[string]map[int]*SmsMessage = make(map[string]map[int]*SmsMessage)

	for sms := range messages {
		tmpsms := *sms
		date, _ := strconv.Atoi(sms.Timestamp)
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
		fmt.Fprintln(file, "Address", address)
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
	done <- true
}
