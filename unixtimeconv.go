package main

import (
	"fmt"
	"strconv"
	"time"
)

// UnixTimeParseString unix string time parse
func UnixTimeParseString(ut string) (time, nano string) {
	return ut[:10], ut[10:]
}

// UnixTimeParseStringInt unix string time parse
func UnixTimeParseStringInt(ut string) (time, nano int64) {
	t, n := ut[:10], ut[10:]
	tv, _ := strconv.ParseInt(t, 10, 64)
	nv, _ := strconv.ParseInt(n, 10, 64)
	return tv, nv
}

// UnixTimeStringWithNanoToPrintable unix string time parse
func UnixTimeStringWithNanoToPrintable(ut string) string {
	return fmt.Sprintf("%s", time.Unix(UnixTimeParseStringInt(ut)).Format(time.RFC1123))
}
