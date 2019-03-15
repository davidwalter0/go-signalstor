// Parsing arbitrary JSON using interfaces in Go
// Demonstrates how to parse JSON with abritrary key names
// See https://blog.golang.org/json-and-go for more info on generic JSON parsing

package main

import (
	"encoding/json"
	"fmt"
)

var jsonBytes = []byte(`
{"sms": {"-sc_toa": "null", "-service_center": "null", "-toa": "null", "-locked": "0", "-protocol": "0", "-date": "1493139630014", "-type": "1", "-subject": "null", "-read": "1", "-status": "-1", "-address": "22000", "-contact_name": "null", "-readable_date": "Tue, 25 Apr 2017 13:00:30 EDT", "-body": "Account notification: The password for your Google Account watchtheriver@gmail.com was recently changed. google.com/password"}}`)

// Item struct; we want to create these from the JSON above
type KV struct {
	Key   string
	Value string
}

// Implement the String interface for pretty printing Items
func (kv KV) String() string {
	return fmt.Sprintf("KV: key: %s, value: %d", kv.Key, kv.Value)
}

func KeyFix(raw interface{}) (repaired map[string]map[string]string) {
	// JSON object parses into a map with string keys
	kvMap := raw.(map[string]interface{})
	// JSON object parses into a map with string keys
	kvMap = kvMap["sms"].(map[string]interface{})

	if len(kvMap) > 0 {
		repaired = make(map[string]map[string]string)
		repaired["sms"] = make(map[string]string)
	}

	for k, v := range kvMap {
		fmt.Println(k, v)
		if len(k) > 0 && string(k)[0] == '-' {
			repaired["sms"][k[1:]] = v.(string)
		} else {
			repaired["sms"][k] = v.(string)
		}
	}
	return
}

func main() {

	// Unmarshal using a generic interface
	var f interface{}
	err := json.Unmarshal(jsonBytes, &f)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
	}
	fmt.Println(KeyFix(f))
	var b []byte
	if b, err = json.Marshal(KeyFix(f)); err != nil {
		panic(err)
	}
	fmt.Println("b", string(b))
	// fmt.Println("jsonBytes", string(jsonBytes))
	// JSON object parses into a map with string keys
	kvMap := f.(map[string]interface{})
	// JSON object parses into a map with string keys
	kvMap = kvMap["sms"].(map[string]interface{})

	// Loop through the Items; we're not interested in the key, just the values
	for k, v := range kvMap {
		fmt.Println(k, v)
		// Use type assertions to ensure that the value's a JSON object
		// switch kv := v.(type) {
		// // The value is an Item, represented as a generic interface
		// // case KV:
		// // 	fmt.Println(kv)
		// // 	fmt.Println(kv.String())
		// case interface{}:
		// 	fmt.Println(k, v)
		// // Not a JSON object; handle the error
		// default:
		// 	fmt.Println("Expecting a JSON object; got something else")
		// }
	}

}
