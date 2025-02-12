package task_transfer

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	dateTimeFormat                  = "20060102"
	minRepeatDays                   = 1
	maxRepeatDays                   = 400
	maxDaysInMounth                 = 31
	minDaysInMounth                 = -2
	lastDayInMounthArgument         = -1
	PenultimatetDayInMounthArgument = -2
)

var (
	errLenRepeat = errors.New("еrror invalid repeat format")
)

func Contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
func LastDayInMonths(days []int, initialDate time.Time) int {
	t := time.Date(initialDate.Year(), initialDate.Month(), 32, 0, 0, 0, 0, time.UTC)
	var daysInMonth int
	for _, v := range days {
		switch v {
		case -1:
			daysInMonth = 32 - t.Day()
		case -2:
			daysInMonth = 32 - t.Day() - 1
		}
	}
	return daysInMonth
}

func NextDate(now time.Time, date string, repeat string) (string, error) {

	initialDate, err := time.Parse(dateTimeFormat, date)
	if err != nil {
		log.Printf("Error time.Parse date, %v:", err)
		return "", err
	}

	repeatSlice := strings.Split(repeat, " ")

	switch repeatSlice[0] {
	case "":
		return TransferEqualNill(now, initialDate)
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
	if err != nil || repeatDays < minRepeatDays || repeatDays > maxRepeatDays {
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
	if initialDate.Year() < now.Year() {
		initialDate = now
	}
	var realDays []int
	for _, days := range daysOfMounth {
		day, err := strconv.Atoi(days)
		if day > maxDaysInMounth || day < minDaysInMounth || err != nil {
			log.Printf("error invalid day format for months days repeat %v", err)
			return "", err
		}
		realDays = append(realDays, day)
		if day == lastDayInMounthArgument || day == PenultimatetDayInMounthArgument {
			day = LastDayInMonths(realDays, initialDate)
			realDays = append(realDays, day)
		}
	}

	var realMonths []int
	if len(repeatSlice) > 2 {
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
		_, month, day := initialDate.Date()

		if Contains(realDays, day) && (len(realMonths) == 0) || initialDate.After(now) && Contains(realMonths, int(month)) && Contains(realDays, day) {
			break
		}
	}
	return initialDate.Format(dateTimeFormat), nil
}

func TransferEqualNill(now time.Time, initialDate time.Time) (string, error) {
	if !initialDate.After(now) {
		return "", nil
	}
	return initialDate.Format(dateTimeFormat), nil
}
