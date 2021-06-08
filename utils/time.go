package utils

import (
	"encoding/json"
	"fmt"
	"time"

	null "gopkg.in/nullbio/null.v6"
)

func TimeFrom(t time.Time) Time {
	return Time{null.TimeFrom(t)}
}

func Now() Time {
	return Time{null.TimeFrom(time.Now())}
}

func NowInSecond() Time {
	now := time.Now()
	return Time{null.TimeFrom(time.Unix(now.Unix(), 0))}
}

const format = "2006-01-02 15:04:05"

const formatMs = "2006-01-02 15:04:05.999"

type Time struct {
	null.Time
}

// Unix returns t as a Unix time, the number of seconds elapsed
// since January 1, 1970 UTC. The result does not depend on the
// location associated with t.
// If time not valid, 0 returned
func (t Time) Unix() int64 {
	if !t.Valid {
		return 0
	}
	return t.Time.Time.Unix()
}

// UnixNano returns t.Time.Time.UnixNano().
// If time not valid, 0 returned
func (t Time) UnixNano() int64 {
	if !t.Valid {
		return 0
	}
	return t.Time.Time.UnixNano()
}

// UnixMilli returns t.Time.Time.UnixNano()/1e6.
// If time not valid, 0 returned
func (t Time) UnixMilli() int64 {
	if !t.Valid {
		return 0
	}
	return t.Time.Time.UnixNano() / 1e6
}

func (t Time) String() string {
	if !t.Valid {
		return string(null.NullBytes)
	}
	return t.Time.Time.Format(format)
}
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return null.NullBytes, nil
	}
	return json.Marshal(t.Time.Time.Format(format))
}
func (t *Time) UnmarshalJSON(b []byte) error {
	if b == nil || len(b) == 0 || string(b) == string(null.NullBytes) {
		t.Valid = false
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	loc, _ := time.LoadLocation("Local")
	val, err := time.ParseInLocation(format, s, loc)
	if err != nil {
		val, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}
	t.SetValid(val)
	return nil
}

func (t *Time) SetValid(v time.Time) {
	if v.Unix() <= 0 {
		return
	}
	t.Time.SetValid(v)
}

func (t *Time) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		t.SetValid(x)
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("null: cannot scan type %T into null.Time: %v", value, value)
	}
	return err
}

type TimeMs struct {
	Time
}

func (t TimeMs) String() string {
	if !t.Valid {
		return string(null.NullBytes)
	}
	return t.Time.Time.Time.Format(formatMs)
}
func (t TimeMs) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return null.NullBytes, nil
	}
	return json.Marshal(t.Time.Time.Time.Format(formatMs))
}
func (t *TimeMs) UnmarshalJSON(b []byte) error {
	if b == nil || len(b) == 0 || string(b) == string(null.NullBytes) {
		t.Valid = false
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	loc, _ := time.LoadLocation("Local")
	val, err := time.ParseInLocation(formatMs, s, loc)
	if err != nil {
		val, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}
	t.SetValid(val)
	return nil
}
