package utils

import (
	"net/http"
	"time"

	"example.com/m/errs"
)

const dateLayout = "2006-01-02"

func ValidateDateRange(startStr, endStr string) error {
	start, err := time.Parse(dateLayout, startStr)
	if err != nil {
		return errs.New("invalid start date format", http.StatusBadRequest)
	}

	end, err := time.Parse(dateLayout, endStr)
	if err != nil {
		return errs.New("invalid end date format", http.StatusBadRequest)
	}

	if start.After(end) {
		return errs.New("start date must be before or equal to end date", http.StatusBadRequest)
	}

	return nil
}
