package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/RedrikShuhartRed/finalTODO/db"
	"github.com/RedrikShuhartRed/finalTODO/models"
	"github.com/RedrikShuhartRed/finalTODO/query"
	"github.com/RedrikShuhartRed/finalTODO/task_transfer"
)

const (
	dateTimeFormat = "20060102"
)

var (
	errEmptyTitle = errors.New("error Decode request body, Task title is empty")
)

func jsonError(w http.ResponseWriter, message string) {
	resp := map[string]string{"error": message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(resp)
}

func GetNextDate(w http.ResponseWriter, r *http.Request) {

	date := r.URL.Query().Get("date")
	now := r.URL.Query().Get("now")
	repeat := r.URL.Query().Get("repeat")

	nowTime, err := time.Parse(dateTimeFormat, now)
	if err != nil {
		log.Printf("error time.Parse now %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	result, err := task_transfer.NextDate(nowTime, date, repeat)
	if err != nil {
		log.Printf("error time.Parse now %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(result)); err != nil {
		log.Printf("error time.Parse now %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func AddNewTask(w http.ResponseWriter, r *http.Request) {
	dbs := db.GetDB()
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	now := time.Now()
	if err != nil {
		log.Printf("error Decode request body, %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(task.Title) == 0 {
		log.Printf("error %v", errEmptyTitle)
		jsonError(w, errEmptyTitle.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parseDate, err := time.Parse(dateTimeFormat, task.Date)
	if err != nil && (len(task.Date) == 0) || task.Date == "today" {
		task.Date = now.Format(dateTimeFormat)
	} else if err != nil && len(task.Date) != 0 {
		log.Printf("error %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err == nil && parseDate.Before(now) && len(task.Repeat) == 0 {
		task.Date = now.Format(dateTimeFormat)
	} else if err == nil && parseDate.Before(now) && len(task.Repeat) != 0 {
		task.Date, err = task_transfer.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("error %v", err)
			jsonError(w, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		task.Date, err = task_transfer.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("error %v", err)
			jsonError(w, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	res, err := dbs.Exec(query.AddNewTask,
		sql.Named("title", task.Title),
		sql.Named("date", task.Date),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		log.Printf("error insert into scheduler, %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	lastId, err := (res.LastInsertId())
	if err != nil {
		log.Printf("error insert into scheduler, %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := map[string]int64{
		"id": lastId,
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error Encode response, %v", err)
		jsonError(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return

	}
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	dbs := db.GetDB()
	tasks := []models.Task{}
	search := r.URL.Query().Get("search")
	limit := 20
	if search == "" {
		limit = 50
	}

	var rows *sql.Rows
	var err error
	timeFormat := "02.01.2006"
	if search == "" {
		rows, err = dbs.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit",
			sql.Named("limit", limit))
		if err != nil {
			log.Printf("error reading data from database: %v", err)
			jsonError(w, "Error reading data from database")
			return
		}
	} else {
		searchdate, err := time.Parse(timeFormat, search)
		if err != nil {
			rows, err = dbs.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :title OR comment LIKE :comment ORDER BY date LIMIT :limit",
				sql.Named("title", "%"+search+"%"),
				sql.Named("comment", "%"+search+"%"),
				sql.Named("limit", limit),
			)

		} else {
			correctsearchdate := searchdate.Format("20060102")
			rows, err = dbs.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date", sql.Named("date", correctsearchdate))
		}
		if err != nil {
			log.Printf("error reading data from database: %v", err)
			jsonError(w, "Error reading data from database")
			return
		}

	}
	defer rows.Close()

	for rows.Next() {
		task := models.Task{}
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		tasks = append(tasks, task)
	}
	if err != nil {
		log.Printf("error Scan data in Task: %v", err)
		jsonError(w, "error Scan data in Task")
		return
	}

	response := map[string][]models.Task{"tasks": tasks}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

	}
}
