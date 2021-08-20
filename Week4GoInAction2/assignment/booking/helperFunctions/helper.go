/*
Package helperFunctions contains common functions
and error types used by the application.
*/
package booking

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrInvalidDate     = errors.New("Invalid Date Entered")
	ErrInvalidOption   = errors.New("Invalid option chosen. Please retry.")
	ErrInsufficientArg = errors.New("Missing information. Please check your input and retry.")
	ErrAlreadyTaken    = errors.New("This time slot is unavailable. Please retry.")
	ErrEmptyFields     = errors.New("Field(s) cannot be empty. Please navigate back and retry.")
)

//ConvertDateTime converts a date string (YYYYMMDD) and a time string (hhmm) into a time.Time object.
func ConvertDateTime(date, tm string) (time.Time, error) {
	if len(date) != 8 || len(tm) != 4 {
		return time.Time{}, errors.Wrap(ErrInvalidDate, "incorrect date or time")
	}
	yy, _ := strconv.Atoi(date[:4])
	mm, _ := strconv.Atoi(date[4:6])
	dd, _ := strconv.Atoi(date[6:])
	hh, _ := strconv.Atoi(tm[:2])
	m, _ := strconv.Atoi(tm[2:])
	return time.Date(yy, time.Month(mm), dd, hh, m, 0, 0, time.Local), nil
}
