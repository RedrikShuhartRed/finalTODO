package models

import (
	"log"
	"strings"
	"time"

	"github.com/RedrikShuhartRed/finalTODO/task_transfer"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

const (
	dateTimeFormat = "20060102"
)

func CheckDate(task *Task) (string, error) {
	now := time.Now()
	parseDate, err := time.Parse(dateTimeFormat, task.Date)

	if err != nil && (len(task.Date) == 0) || task.Date == "today" {
		task.Date = now.Format(dateTimeFormat)
		return task.Date, nil
	} else if err != nil && len(task.Date) != 0 {
		log.Printf("error %v", err)

	} else if err == nil && parseDate.Before(now) && len(task.Repeat) == 0 {
		task.Date = now.Format(dateTimeFormat)

	} else if parseDate.Year() == now.Year() && parseDate.Month() == now.Month() && parseDate.Day() == now.Day() {
		task.Date = now.Format(dateTimeFormat)

	} else if parseDate.After(time.Now()) && strings.HasPrefix(task.Repeat, "d") {
		task.Date = parseDate.Format(dateTimeFormat)

	} else if err == nil && parseDate.Before(now) && len(task.Repeat) != 0 {
		task.Date, err = task_transfer.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("error %v", err)
		}
	} else {
		task.Date, err = task_transfer.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("error %v", err)
		}
	}
	return task.Date, err

}
func CheckDoneDate(task *Task) (string, error) {
	now := time.Now()
	parseDate, err := time.Parse(dateTimeFormat, task.Date)
	if err != nil && (len(task.Date) == 0) || task.Date == "today" {
		task.Date = now.Format(dateTimeFormat)
	} else if err != nil && len(task.Date) != 0 {
		log.Printf("error %v", err)
	} else if err == nil && parseDate.Before(now) && len(task.Repeat) == 0 {
		task.Date = now.Format(dateTimeFormat)
	} else if err == nil && parseDate.Before(now) && len(task.Repeat) != 0 {
		task.Date, err = task_transfer.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("error %v", err)
		}
	} else {
		task.Date, err = task_transfer.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("error %v", err)

		}
	}
	return task.Date, err
}
