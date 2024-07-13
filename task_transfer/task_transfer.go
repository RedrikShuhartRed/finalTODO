package tasktransfer_test

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	dateTimeFormat = "20060102"
)

var (
	errLenRepeat   = errors.New("еrror invalid repeat format")
	errInvalidDate = errors.New("error invalid date format")
	errRepeatEmpty = errors.New("error repeat rule is empty")
)

func Contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func TransferForYear(now time.Time, initialDate time.Time) (string, error) {
	for {
		initialDate = initialDate.AddDate(1, 0, 0)
		if initialDate.After(now) {
			break
		}
	}
	return initialDate.Format(dateTimeFormat), nil
}

func TransferForDay(now time.Time, initialDate time.Time, repeatSlice []string) (string, error) {
	if len(repeatSlice) != 2 {
		log.Printf("error invalid repeat lenght: %v", errLenRepeat)
		return "", errLenRepeat
	}
	repeatDays, err := strconv.Atoi(repeatSlice[1])
	if err != nil || repeatDays < 1 || repeatDays > 400 {
		log.Printf("error invalid repeatDays interval: %v", err)
		return "", err
	}

	for {
		initialDate = initialDate.AddDate(0, 0, repeatDays)
		if initialDate.After(now) {
			break
		}
	}
	return initialDate.Format(dateTimeFormat), nil

}

func TransferForSpecifiedDayWeek(now time.Time, initialDate time.Time, repeatSlice []string) (string, error) {
	if len(repeatSlice) != 2 {
		log.Printf("error invalid repeat lenght, %v", errLenRepeat)
		return "", errLenRepeat
	}
	daysOfWeek := strings.Split(repeatSlice[1], ",")

	var realDays []int
	for _, days := range daysOfWeek {
		day, err := strconv.Atoi(days)
		if err != nil || day < 1 || day > 7 {
			log.Printf("error invalid day for week repeat %v", err)
			return "", err
		}
		realDays = append(realDays, day)
	}
	for {
		initialDate = initialDate.AddDate(0, 0, 1)
		dayOfWeek := int(initialDate.Weekday())
		if dayOfWeek == 0 {
			dayOfWeek = 7
		}

		if Contains(realDays, dayOfWeek) && initialDate.After(now) {
			return initialDate.Format(dateTimeFormat), nil
		}
	}
}

func TransferForSpecifiedDayMonth(now time.Time, initialDate time.Time, repeatSlice []string) (string, error) {
	if len(repeatSlice) != 2 && len(repeatSlice) != 3 {
		log.Printf("error invalid repeat lenght,%v", errLenRepeat)
		return "", errLenRepeat
	}
	daysOfMounth := strings.Split(repeatSlice[1], ",")

	var realDays []int
	for _, days := range daysOfMounth {
		day, err := strconv.Atoi(days)
		if day > 31 || day < -2 || err != nil {
			log.Printf("error invalid day format for months days repeat %v", err)
			return "", err
		}
		realDays = append(realDays, day)
	}

	var realMonths []int
	if len(repeatSlice) == 3 {
		monthsOfYear := strings.Split(repeatSlice[2], ",")
		for _, months := range monthsOfYear {
			month, err := strconv.Atoi(months)
			if month > 12 || month < 1 || err != nil {
				log.Printf("error invalid months format for months days repeat %v", err)
				return "", err
			}
			realMonths = append(realMonths, month)

		}

	}
	for {
		initialDate = initialDate.AddDate(0, 0, 1)
		_, Month, Day := initialDate.Date()
		if Contains(realDays, Day) && (len(realMonths) == 0) || Contains(realMonths, int(Month)) && initialDate.After(now) {
			break
		}
	}
	return initialDate.Format(dateTimeFormat), nil

}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	initialDate, err := time.Parse(dateTimeFormat, date)
	if err != nil {
		log.Printf("Error time.Parse date, %v:", err)
		return "", errInvalidDate
	}

	if repeat == "" {
		log.Printf("Error repeat rule is empty, %v:", errRepeatEmpty)
		return "", errRepeatEmpty
	}
	repeatSlice := strings.Split(repeat, " ")

	switch repeatSlice[0] {
	case "y":
		return TransferForYear(now, initialDate)
	case "d":
		return TransferForDay(now, initialDate, repeatSlice)
	case "w":
		return TransferForSpecifiedDayWeek(now, initialDate, repeatSlice)
	case "m":
		return TransferForSpecifiedDayMonth(now, initialDate, repeatSlice)
	default:
		log.Printf("error unsupported repeat format: %v", errLenRepeat)
		return "", errLenRepeat
	}

}
