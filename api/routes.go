package api

import (
	"net/http"

	"github.com/RedrikShuhartRed/finalTODO/api/handlers"
	"github.com/RedrikShuhartRed/finalTODO/config"
	"github.com/RedrikShuhartRed/finalTODO/db"
	"github.com/gorilla/mux"
)

func addNewTaskHandler(storage *db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AddNewTask(w, r, storage)
	}
}
func getAllTasksHandler(storage *db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllTasks(w, r, storage)
	}
}
func getTaskByIdHandler(storage *db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.GetTasksById(w, r, storage)
	}
}
func UpdateTaskHandler(storage *db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateTask(w, r, storage)
	}
}
func DoneTaskHandler(storage *db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.DoneTask(w, r, storage)
	}
}
func DeleteTaskHandler(storage *db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteTask(w, r, storage)
	}
}
func authorizationGetTokenHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AuthorizationGetToken(w, r, cfg)
	}
}

var RegisterTasksStoreRoutes = func(router *mux.Router, storage *db.Storage, cfg *config.Config) {
	router.HandleFunc("/api/nextdate", handlers.GetNextDate).Methods("GET")
	router.HandleFunc("/api/task", handlers.Auth(addNewTaskHandler(storage), cfg)).Methods("POST")
	router.HandleFunc("/api/tasks", handlers.Auth(getAllTasksHandler(storage), cfg)).Methods("GET")
	router.HandleFunc("/api/task", handlers.Auth(getTaskByIdHandler(storage), cfg)).Methods("GET")
	router.HandleFunc("/api/task", handlers.Auth(UpdateTaskHandler(storage), cfg)).Methods("PUT")
	router.HandleFunc("/api/task/done", handlers.Auth(DoneTaskHandler(storage), cfg)).Methods("POST")
	router.HandleFunc("/api/task", handlers.Auth(DeleteTaskHandler(storage), cfg)).Methods("DELETE")
	router.HandleFunc("/api/signin", authorizationGetTokenHandler(cfg)).Methods("POST")
}
