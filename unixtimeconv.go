package signalstor // 	"github.com/davidwalter0/go-signalstor"

import (
	"fmt"
	"strconv"
	"time"
)

// UnixTime create a time.Time from unix time t
func UnixTime(t int64) time.Time {
	return time.Unix(t/1000, (t%1000)*1000000)
}

// UnixTimeWithMilli create a time.Time from unix time t ( second resolution ) +
// milli seconds
func UnixTimeWithMilli(t, m int64) time.Time {
	return time.Unix(t, m*1000000)
}

// UnixTimeParseString unix string time parse
func UnixTimeParseString(ut string) (time, nano string) {
	return ut[:10], ut[10:]
}

// UnixTimeMsResolutionStr2Int unix string time parse
func UnixTimeMsResolutionStr2Int(ut string) (tv, ms int64) {
	t, n := ut[:10], ut[10:]
	tv, _ = strconv.ParseInt(t, 10, 64)
	ms, _ = strconv.ParseInt(n, 10, 64)
	return tv, ms
}

// UnixTimeStringWithMsToPrintable unix string time parse
func UnixTimeStringWithMsToPrintable(ut string) string {
	return fmt.Sprintf("%s", time.Unix(UnixTimeMsResolutionStr2Int(ut)).Format(time.RFC1123))
}
