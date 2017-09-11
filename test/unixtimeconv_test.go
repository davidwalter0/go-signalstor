package xml2json

import (
	"fmt"
	"testing"

	"github.com/davidwalter0/xml2json"
)

var debug = false

type UnixDateSplit struct {
	tstr string
	nstr string
	t    int64
	n    int64
}

type Table struct {
	UnixDate    string
	RFC1123     string
	DefaultDate string
	UnixDateSplit
	UnixTime int64
}

// Table of test inputs and expected values
var Tables = []Table{
	{
		"1461621630014",
		"2016-04-25 18:00:30 -0400 EDT",
		"Mon, 25 Apr 2016 18:00:30 EDT",
		UnixDateSplit{
			"1461621630",
			"014",
			1461621630,
			14,
		},
		1461621630014,
	},
	{
		"1461621630014",
		"2016-04-25 18:00:30 -0500 CDT",
		"Mon, 25 Apr 2016 18:00:30 CDT",
		UnixDateSplit{
			"1461621630",
			"014",
			1461621630,
			14,
		},
		1461621630014,
	},
}

func Test_UnixTimeFunctions(t *testing.T) {
	var u UnixDateSplit
	var unixTime int64

	for _, itm := range Tables {
		u = itm.UnixDateSplit
		if debug {
			fmt.Println(u)
		}
		unixTime = xml2json.UnixTime(itm.UnixTime).Unix()*1000 + u.n
		if unixTime != itm.UnixTime {
			t.Errorf("UnixTime() was incorrect, want %d got: %d.", itm.UnixTime, unixTime)
		}

		unixTime = xml2json.UnixTime(itm.UnixTime).Unix()
		if unixTime != itm.UnixDateSplit.t {
			t.Errorf("UnixTime() was incorrect, want %d got: %d.", itm.UnixDateSplit.t, unixTime)
		}
		u = itm.UnixDateSplit
		unixTime = xml2json.UnixTimeWithMilli(u.t, u.n).Unix()*1000 + u.n
		if unixTime != itm.UnixTime {
			t.Errorf("UnixTime() was incorrect, want %d got: %d.", itm.UnixTime, unixTime)
		}
		if debug {
			fmt.Println(xml2json.UnixTimeParseString(itm.UnixDate))
		}
		{
			l, r := xml2json.UnixTimeMsResolutionStr2Int(itm.UnixDate)
			if debug {
				fmt.Printf(">>>  %10d   :%03d\n", l, r)
				fmt.Printf("time %10d:ms:%03d\n", l, r)
			}
			timeString := xml2json.UnixTimeStringWithMsToPrintable(itm.UnixDate)
			ln := len(timeString)
			if timeString[:ln-4] != itm.DefaultDate[:ln-4] {
				t.Errorf("UnixTime() was incorrect, want %s got: %s.", itm.DefaultDate[:ln-4], timeString[:ln-4])
			}
			if debug {
				fmt.Println(xml2json.UnixTimeStringWithMsToPrintable(itm.UnixDate))
			}
		}

	}

	/*
		est, _ := time.LoadLocation("US/Eastern")
		date := "Tue, 25 Apr 2016 18:00:30 EDT"
		t, err := time.ParseInLocation(time.RFC1123, date, est)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("ParseInLocation", t)
		{
			cst, _ := time.LoadLocation("US/Central")
			date := "Tue, 25 Apr 2016 18:00:30 EDT"
			t, err := time.ParseInLocation(time.RFC1123, date, cst)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("ParseInLocation", t)
		}

		t, err = time.Parse(time.RFC1123, date)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("parse          ", t)
		fmt.Println("format         ", t.Format(time.RFC1123))

		date = "Tue, 25 Apr 2016 18:00:30 CDT"
		t, err = time.ParseInLocation(time.RFC1123, date, est)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("ParseInLocation", t)

		date = "Tue, 25 Apr 2016 18:00:30"
		t, err = time.ParseInLocation(time.RFC1123, date, est)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("ParseInLocation", t)

		t, err = time.Parse(time.RFC1123, date)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("parse          ", t)
		fmt.Println("format         ", t.Format(time.RFC1123))
		ut := "1461621630014"
		{
			t, m := ut[:10], ut[10:]
			uti, _ := strconv.ParseInt(ut, 10, 64)
			mv, _ := strconv.ParseInt(m, 10, 64)
			tv, _ := strconv.ParseInt(t, 10, 64)
			fmt.Println("unix ", t, tv, m, mv)
			tm1 := UnixTime(uti)
			tm2 := UnixTimeWithMilli(tv, mv)
			fmt.Println("tm", tm1)
			fmt.Printf("tm %v %T\n", tm1, tm1)
			fmt.Println("tm", tm2)
			fmt.Printf("tm %v %T\n", tm2, tm2)
		}
		fmt.Println(UnixTimeParseString(ut))
		{
			t, m := UnixTimeMsResolutionStr2Int(ut)
			fmt.Printf(">>>  %10d   :%03d\n", t, m)
			fmt.Printf("time %10d:ms:%03d\n", t, m)
			fmt.Println(UnixTimeStringWithMsToPrintable(ut))
		}
	*/
}

func Test_SmsMessageString(t *testing.T) {

	expected := `Contact : name
Date    : Mon, 25 Apr 2016 18:00:30 EDT
Address : +15555555555
Message : More text
`
	var message xml2json.SmsMessage = xml2json.SmsMessage{
		ContactName:  "name",
		ReadableDate: "Mon, 25 Apr 2016 18:00:30 EDT",
		Address:      "+15555555555",
		Body:         "More text",
	}

	if expected != message.String() {
		t.Errorf("SmsMessage.String(), want\n%s\ngot\n%s\n", expected, message.String())
	}
}
