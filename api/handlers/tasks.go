package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	check "github.com/RedrikShuhartRed/finalTODO/api/checkQueryParam"
	"github.com/RedrikShuhartRed/finalTODO/api/middleware"
	"github.com/RedrikShuhartRed/finalTODO/config"
	"github.com/RedrikShuhartRed/finalTODO/db"
	"github.com/RedrikShuhartRed/finalTODO/models"
	"github.com/RedrikShuhartRed/finalTODO/task_transfer"
)

const (
	dateTimeFormat = "20060102"
	timeFormat     = "02.01.2006"
)

type Handler struct {
	storage *db.Storage
	cfg     *config.Config
}

var (
	errEmptyTitle = errors.New("error Decode request body, Task title is empty")
	Jerr          middleware.JsonErr
)

func NewHandler(storage *db.Storage, config *config.Config) *Handler {
	return &Handler{
		storage: storage,
		cfg:     config,
	}
}

func (h *Handler) GetNextDate(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	now := r.URL.Query().Get("now")
	repeat := r.URL.Query().Get("repeat")

	nowTime, err := time.Parse(dateTimeFormat, now)
	if err != nil {
		log.Printf("error time.Parse now %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusBadRequest)
		return

	}

	result, err := task_transfer.NextDate(nowTime, date, repeat)
	if err != nil {
		log.Printf("error task transfer %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(result)); err != nil {
		log.Printf("error w.Write %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AddNewTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("error Decode request body, %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(task.Title) == 0 {
		log.Printf("error %v", errEmptyTitle)
		Jerr.JsonError(w, errEmptyTitle.Error(), http.StatusBadRequest)
		return
	}

	task.Date, err = check.CheckDate(task)
	if err != nil {
		log.Printf("error %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	lastId, err := h.storage.AddNewTask(task)
	if err != nil {
		log.Printf("error insert into scheduler, %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]int64{
		"id": lastId,
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error Encode response, %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusInternalServerError)
		return

	}
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	var err error
	search := r.URL.Query().Get("search")

	tasks, err = check.CheckSearch(search, h.storage)
	if err != nil {
		log.Printf("error get all tasks: %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string][]models.Task{"tasks": tasks}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error Encode response, %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetTasksById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := check.CheckId(id)

	if err != nil {
		log.Printf("error get id task: %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.storage.GetTasksById(id)
	if err != nil {
		log.Printf("error get task from DB: %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Printf("error Encode response, %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task *models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("error Decode request body, %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = check.CheckId(task.ID)

	if err != nil {
		log.Printf("error get id task: %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(task.Title) == 0 {
		log.Printf("error %v", errEmptyTitle)
		Jerr.JsonError(w, errEmptyTitle.Error(), http.StatusBadRequest)
		return
	}

	task.Date, err = check.CheckDate(*task)
	if err != nil {
		log.Printf("error %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rowsAffected, err := h.storage.UpdateTask(*task)
	if err != nil {
		log.Printf("error getting rows affected: %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Printf("no rows affected: %v", err)
		Jerr.JsonError(w, "no rows affected", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(struct{}{})
	if err != nil {
		log.Printf("error Encode response, %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func (h *Handler) DoneTask(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	err := check.CheckId(id)
	if err != nil {
		log.Printf("error get id task: %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.storage.GetTasksById(id)
	if err != nil {
		log.Printf("error get task from DB: %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(task.Repeat) == 0 {
		h.DeleteTask(w, r)
		return
	}

	task.Date, err = check.CheckDoneDate(*task)
	if err != nil {
		log.Printf("error %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	rowsAffected, err := h.storage.UpdateTask(*task)
	if err != nil {
		log.Printf("error getting rows affected: %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		log.Printf("no rows affected: %v", err)
		Jerr.JsonError(w, "no rows affected", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(struct{}{})
	if err != nil {
		log.Printf("error Encode response, %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := check.CheckId(id)
	if err != nil {
		log.Printf("error get id task: %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	rowsAffected, err := h.storage.DeleteTask(id)
	if err != nil {
		log.Printf("error getting rows affected: %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		log.Printf("no rows affected: %v", err)
		Jerr.JsonError(w, "no rows affected", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(struct{}{})
	if err != nil {
		log.Printf("error Encode response, %v", err)
		Jerr.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
