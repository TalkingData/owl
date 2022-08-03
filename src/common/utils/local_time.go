package utils

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	timestampFormat = "2006-01-02 15:04:05"
)

type LocalTime struct {
	time.Time
}

// NewLocalTimeByTimestamp 通过时间戳获取LocalTime
func NewLocalTimeByTimestamp(ts int64) *LocalTime {
	return &LocalTime{Time: time.Unix(ts, 0)}
}

func (t *LocalTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(timestampFormat))
	return []byte(formatted), nil
}

func (t *LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t == nil {
		return nil, nil
	} else if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *LocalTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalTime{Time: value}
		return nil
	}
	return fmt.Errorf("Can not convert %v to timestamp", v)
}

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timestampFormat+`"`, string(data), time.Local)
	*t = LocalTime{Time: now}
	return
}

func (t *LocalTime) String() string {
	return t.Format(timestampFormat)
}
