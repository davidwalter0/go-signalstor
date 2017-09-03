// go run time.go unixtimeconv.go

package main

import (
	"fmt"
	// "os"
	"strconv"
	"time"
)

func main() {
	// func LoadLocation(name string) (*Location, error)
	// LoadLocation returns the Location with the given name.

	// If the name is "" or "UTC", LoadLocation returns UTC. If the name is "Local", LoadLocation returns Local.

	// Otherwise, the name is taken to be a location name corresponding to a file in the IANA Time Zone database, such as "America/New_York".
	// fmt.Println(time.LoadLocation("UTC"))
	// // fmt.Println(time.LoadLocation("PST8PDT"))
	// fmt.Println(time.LoadLocation("EST"))
	// fmt.Println(time.LoadLocation("EDT"))

	// loc, _ := time.LoadLocation("US/Eastern")
	// fmt.Println(loc)

	// date := "Tue, 25 Apr 2016 18:00:30 EDT"
	// date = "Tue, 25 Apr 2016 13:16:42 EDT"
	// // t, err := time.Parse(time.RFC1123, date)
	// t, err := time.ParseInLocation(time.RFC1123, date, loc)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(t.Unix())
	// fmt.Println(t)
	// tm := time.Unix(t.Unix(), 14)
	// fmt.Println("here", tm, fmt.Sprintf("\n%d%03d", t.Unix(), 14))
	// tm = time.Unix(t.Unix(), 0)
	// fmt.Println(tm.Format(time.RFC1123))
	// fmt.Println("here", tm, fmt.Sprintf("\n%d%03d", t.Unix(), 14))
	// if false {
	// 	os.Exit(0)
	// }
	loc, _ := time.LoadLocation("US/Eastern")
	date := "Tue, 25 Apr 2016 18:00:30 EDT"
	t, err := time.ParseInLocation(time.RFC1123, date, loc)
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

	date = "Tue, 25 Apr 2016 18:00:30 CDT"
	t, err = time.ParseInLocation(time.RFC1123, date, loc)
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
		t, n := ut[:10], ut[10:]
		nv, _ := strconv.ParseInt(n, 10, 64)
		tv, _ := strconv.ParseInt(t, 10, 64)
		fmt.Println("unix ", t, tv, n, nv)
		tm := time.Unix(tv, nv)
		fmt.Println("tm", tm)
		fmt.Printf("tm %v %T\n", tm, tm)
	}
	fmt.Println(UnixTimeParseString(ut))
	{
		t, n := UnixTimeParseStringInt(ut)
		fmt.Printf(">>>  %10d   :%03d\n", t, n)
		fmt.Printf("time %10d:ns:%03d\n", t, n)
		fmt.Println(UnixTimeStringWithNanoToPrintable(ut))
	}
}
