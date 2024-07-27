package checkqueryparam

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/RedrikShuhartRed/finalTODO/db"
	"github.com/RedrikShuhartRed/finalTODO/models"
	"github.com/RedrikShuhartRed/finalTODO/task_transfer"
)

const (
	dateTimeFormat = "20060102"
	timeFormat     = "02.01.2006"
)

var errEmptyId = errors.New("error Decode request body, Task id is empty")

func CheckSearch(search string, storage *db.Storage) ([]models.Task, error) {
	var tasks []models.Task
	var err error
	if search == "" {
		tasks, err = storage.GetAllTasksWithoutSearch()
	} else {
		_, err = time.Parse(timeFormat, search)
		if err == nil {
			tasks, err = storage.GetAllTasksWithDateSearch(search)
		} else {
			tasks, err = storage.GetAllTasksWithStringSearch(search)
		}
	}
	return tasks, err
}

func CheckId(id string) error {
	if id == "" {
		log.Printf("error get id task: id == \"\"")
		return errEmptyId
	}
	_, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("error get id task, id not int: %v", err)
		return err
	}
	return nil
}

func CheckDate(task models.Task) (string, error) {
	now := time.Now()
	parseDate, err := time.Parse(dateTimeFormat, task.Date)

	if err != nil {
		if (len(task.Date) == 0) || task.Date == "today" {
			task.Date = now.Format(dateTimeFormat)
			return task.Date, nil
		}
		return "", err
	}

	if parseDate.Before(now) && len(task.Repeat) == 0 {
		task.Date = now.Format(dateTimeFormat)
		return task.Date, nil
	}
	if parseDate.Year() == now.Year() && parseDate.Month() == now.Month() && parseDate.Day() == now.Day() {
		task.Date = now.Format(dateTimeFormat)
		return task.Date, nil
	}
	if parseDate.Before(now) && len(task.Repeat) != 0 {
		task.Date, err = task_transfer.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return "", err
		}
		return task.Date, nil
	}

	return task.Date, nil

}

func CheckDoneDate(task models.Task) (string, error) {
	now := time.Now()
	parseDate, err := time.Parse(dateTimeFormat, task.Date)

	if err != nil {
		if (len(task.Date) == 0) || task.Date == "today" {
			task.Date = now.Format(dateTimeFormat)
			return task.Date, nil
		}
		return "", err
	}

	if parseDate.Before(now) && len(task.Repeat) == 0 {
		task.Date = now.Format(dateTimeFormat)
		return task.Date, nil
	}

	if parseDate.Before(now) && len(task.Repeat) != 0 {
		task.Date, err = task_transfer.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return "", err
		}
		return task.Date, nil
	}

	task.Date, err = task_transfer.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return "", err
	}

	return task.Date, nil
}
