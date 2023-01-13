package oth

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	layout = "2006-01-02 15:04:05"
)

type TimeStamp time.Time

func (ts *TimeStamp) ToTime() time.Time {
	return time.Time(*ts)
}

func (ts *TimeStamp) String() string {
	return ts.ToTime().Format(layout)
}

func (ts *TimeStamp) MarshalJSON() ([]byte, error) {
	/*t := time.Time(*ts)
	return []byte(strconv.FormatInt(t.UnixNano()/1000000, 10)), nil*/

	return []byte(`"` + ts.String() + `"`), nil
}

func (ts *TimeStamp) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	/*millis, err := strconv.ParseInt(string(data), 10, 64)
	*ts = TimeStamp(time.Unix(0, millis*int64(time.Millisecond)))*/

	//t, err := time.Parse(layout, string(bytes.Trim(data, `"`)))
	shanghai, _ := time.LoadLocation("Asia/Shanghai")
	t, err := time.ParseInLocation(layout, string(bytes.Trim(data, `"`)), shanghai)
	*ts = TimeStamp(t)

	return err
}

func (ts *TimeStamp) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(*ts)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func (ts *TimeStamp) Scan(v interface{}) error {
	val, ok := v.(time.Time)
	if ok {
		*ts = TimeStamp(val)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
