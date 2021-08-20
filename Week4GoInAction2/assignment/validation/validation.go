package validation

import (
	"errors"
	"regexp"
)

//Errors returned by failures to validate
var (
	ErrInvalidAlphaNum = errors.New("This field(s) must be alphanumerical (a-z, A-Z, 0-9 and _).")
	ErrInvalidAlphabet = errors.New("This field(s) must contain alphabets only (a-z, A-Z).")
	ErrInvalidDate     = errors.New("Invalid date. Date should be entered in YYYYMMDD format.")
	ErrInvalidTime     = errors.New("Invalid time. Time should be entered in 24-hour hhmm format.")
)

//IsAlphaNumeric tests if the input matches the character set [0-9A-Za-z_]
func IsAlphaNumeric(input string) bool {
	matched, _ := regexp.MatchString(`^(\w)+$`, input)
	return matched
}

//IsAlphabet tests if the input matches the character set [A-Za-z]
func IsAlphabet(input string) bool {
	//match names with or without spaces in between e.g. james, edward kim, noh im song
	matched, _ := regexp.MatchString(`^[A-Za-z]+( *[A-Za-z]*)( *[A-Za-z]*)$`, input)
	return matched
}

//IsDateString tests if input is of the date format YYYYMMDD
func IsDateString(input string) bool {
	if len(input) != 8 {
		return false
	}
	matched, _ := regexp.MatchString(`^([0-9]{4})(0?[1-9]|1[012])(0?[1-9]|[12][0-9]|3[01])$`, input)
	return matched
}

//IsTimeString tests if input is of the 24-hour format HHMM
func IsTimeString(input string) bool {
	if len(input) != 4 {
		return false
	}
	matched, _ := regexp.MatchString(`^([01]?[0-9]|2[0-3])[0-5][0-9]$`, input)
	return matched
}
