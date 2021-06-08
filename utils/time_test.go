package utils

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	var a = Time{}
	err := json.Unmarshal([]byte(`"0000-00-00 00:00:00"`), &a)
	if err != nil {
		//panic(err)
	}
	log.Println("valid", a.Valid, "unix", a.Time.Time.UTC().Unix())

	{
		var a Time
		assert.Equal(t, int64(0), a.Unix())
		assert.Equal(t, int64(0), a.UnixMilli())
		assert.Equal(t, int64(0), a.UnixNano())
	}
	{
		now := time.Now()
		a := TimeFrom(now)
		assert.Equal(t, now.Unix(), a.Unix())
		assert.Equal(t, now.UnixNano()/1e6, a.UnixMilli())
		assert.Equal(t, now.UnixNano(), a.UnixNano())
	}
}
