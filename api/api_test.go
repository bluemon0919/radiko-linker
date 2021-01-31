package api

import (
	"log"
	"testing"
	"time"
)

func TestPreviousTime(t *testing.T) {
	cases := []struct {
		weekday       time.Weekday
		time          string
		expectHour    int
		expectMinute  int
		expectWeekday time.Weekday
	}{
		{weekday: time.Monday, time: "11:30", expectHour: 11, expectMinute: 30, expectWeekday: time.Monday},
		{weekday: time.Monday, time: "25:30", expectHour: 1, expectMinute: 30, expectWeekday: time.Tuesday},
		{weekday: time.Sunday, time: "24:00", expectHour: 0, expectMinute: 00, expectWeekday: time.Monday},
		{weekday: time.Sunday, time: "29:59", expectHour: 5, expectMinute: 59, expectWeekday: time.Monday},
		{weekday: time.Sunday, time: "06:00", expectHour: 6, expectMinute: 00, expectWeekday: time.Sunday},
	}

	for _, c := range cases {
		tExt, err := PreviousTime(c.weekday, c.time)
		if err != nil {
			log.Fatal(err)
		}
		tm := time.Time(tExt)
		if tm.Hour() != c.expectHour {
			log.Fatal("tm.Hour() is different from the expected value")
		}
		if tm.Minute() != c.expectMinute {
			log.Fatal("tm.Minute() is different from the expected value")
		}
		if tm.Weekday() != c.expectWeekday {
			log.Fatal("tm.Weekday() is different from the expected value")
		}
	}
}

func TestPreviousTimeFail(t *testing.T) {
	// 30時間表記対象の時間
	_, err := PreviousTime(time.Sunday, "30:00")
	if err == nil {
		log.Fatal("err == nil")
	}
}

func TestPreviousTimeURL(t *testing.T) {
	str := PreviousTimeURL("TBS", 0, "14:00")
	if str == "" {
		log.Fatal(str)
	}
}
