package api

import (
	"github.com/gorilla/mux"

	"github.com/RedrikShuhartRed/finalTODO/api/handlers"
	"github.com/RedrikShuhartRed/finalTODO/config"
	"github.com/RedrikShuhartRed/finalTODO/db"
)

var RegisterTasksStoreRoutes = func(router *mux.Router, storage *db.Storage, cfg *config.Config) {
	handler := handlers.NewHandler(storage, cfg)
	router.HandleFunc("/api/nextdate", handler.GetNextDate).Methods("GET")
	router.HandleFunc("/api/task", handler.Auth(handler.AddNewTask)).Methods("POST")
	router.HandleFunc("/api/tasks", handler.Auth(handler.GetAllTasks)).Methods("GET")
	router.HandleFunc("/api/task", handler.Auth(handler.GetTasksById)).Methods("GET")
	router.HandleFunc("/api/task", handler.Auth(handler.UpdateTask)).Methods("PUT")
	router.HandleFunc("/api/task/done", handler.Auth(handler.DoneTask)).Methods("POST")
	router.HandleFunc("/api/task", handler.Auth(handler.DeleteTask)).Methods("DELETE")
	router.HandleFunc("/api/signin", handler.AuthorizationGetToken).Methods("POST")
}
