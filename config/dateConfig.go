package config

import (
	"os"
	"time"
)

type DateTime string

const LayoutDate DateTime = "2006-01-02"
const LayoutDateTime DateTime = "2006-01-02 15:04:05"

func SetTimeZone() {
	tz := os.Getenv("APP_TIMEZONE")
	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic("Failed to get APP_TIMEZONE")
	}
	time.Local = loc
}
