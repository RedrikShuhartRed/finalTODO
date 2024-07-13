package tasktransfer_test

import (
	"log"
	"strconv"
	"strings"
	"time"
)

const dateTimeFormat = "20060102"

func TaskTransfer(now time.Time, date string, repeat string) (string, error) {
	initialDate, err := time.Parse(dateTimeFormat, date)
	if err != nil {
		log.Printf("Error timr.Parse date, %v:", err)
		return "", err
	}

	if repeat == "" {
		log.Printf("Error repeat rule is empty, %v:", err)
		return "", err
	}

	if repeat == "y" {
		for {
			initialDate = initialDate.AddDate(1, 0, 0)
			if initialDate.After(now) {
				break
			}
		}
		return initialDate.Format(dateTimeFormat), nil
	}

	if strings.HasPrefix(repeat, "d ") {
		repeatSlice := strings.Split(repeat, " ")
		if len(repeatSlice) != 2 {
			log.Printf("Error invalid repeat format for daily repeat %v:", err)
			return "", err
		}
		repeatDays, err := strconv.Atoi(repeatSlice[1])
		if err != nil || repeatDays < 1 || repeatDays > 400 {
			log.Printf("Error invslid repeaatDays interval: %v", err)
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

	if strings.HasPrefix(repeat, "w ") {
		repeatSlice := strings.Split(repeat, " ")
		if len(repeatSlice) != 2 {
			log.Printf("Error invalid repeat format for days repeat %v:", err)
			return "", err
		}
		daysOfWeek := strings.Split(repeatSlice[1], ",")

		var realDays []int
		for _, days := range daysOfWeek {
			day, err := strconv.Atoi(days)
			if err != nil || day < 1 || day > 7 {
				log.Printf("Error invalid day for week repeat %v", err)
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
			for _, realDay := range realDays {
				if dayOfWeek == realDay {
					return initialDate.Format(dateTimeFormat), nil
				}
			}
		}
	}

}
