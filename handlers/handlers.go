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

func jsonError(message string) string {
	return `{"error": "` + message + `"}`
}

func GetNextDate(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	date := queryParams.Get("date")
	now := queryParams.Get("now")
	repeat := queryParams.Get("repeat")

	nowTime, err := time.Parse(dateTimeFormat, now)
	if err != nil {
		log.Printf("error time.Parse now %v", err)
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return

	}

	result, err := task_transfer.NextDate(nowTime, date, repeat)
	if err != nil {
		log.Printf("error time.Parse now %v", err)
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(result)); err != nil {
		log.Printf("error time.Parse now %v", err)
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}
}

func AddNewTask(w http.ResponseWriter, r *http.Request) {
	dbs := db.GetDB()
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		log.Printf("error Decode request body, %v", err)
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		log.Printf("error %v", errEmptyTitle)
		http.Error(w, jsonError(errEmptyTitle.Error()), http.StatusBadRequest)
		return
	}

	now := time.Now()
	initialDate, err := task_transfer.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		log.Printf("error NextDate, %v", err)
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return
	}

	res, err := dbs.Exec(query.AddNewTask,
		sql.Named("title", task.Title),
		sql.Named("date", initialDate),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		log.Printf("error insert into scheduler, %v", err)
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
	lastId, err := (res.LastInsertId())
	if err != nil {
		log.Printf("error insert into scheduler, %v", err)
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]int64{
		"id": lastId,
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error Encode response, %v", err)
		http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
		return

	}
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	dbs := db.GetDB()
	tasks := []models.Task{}
	queryParams := r.URL.Query()
	search := queryParams.Get("search")
	limit := "20"
	if len(search) == 0 {
		rows, err := dbs.Query(`SELECT * FROM scheduler ORDER BY date LIMIT :limit`,
			sql.Named("limit", limit))
		if err != nil {
			log.Printf("error read data from database,%v", err)
			http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
			return
		}

		for rows.Next() {
			task := models.Task{}
			err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				log.Printf("error rows data from database,%v", err)
				http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
				return
			}
			tasks = append(tasks, task)

		}
		if err := rows.Err(); err != nil {
			log.Printf("error iterating over rows, %v", err)
			http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
			return
		}

		if err := rows.Close(); err != nil {
			log.Printf("error closing rows, %v", err)
			http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
			return
		}

		if len(tasks) == 0 {
			tasks = []models.Task{}
		}

		response := map[string][]models.Task{"tasks": tasks}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("error Encode response, %v", err)
			http.Error(w, jsonError(err.Error()), http.StatusBadRequest)
			return

		}
		w.WriteHeader(http.StatusOK)
	}
}
