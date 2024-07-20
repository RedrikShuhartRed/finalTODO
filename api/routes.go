package api

import (
	"github.com/RedrikShuhartRed/finalTODO/api/handlers"
	"github.com/gorilla/mux"
)

var RegisterTasksStoreRoutes = func(router *mux.Router) {
	router.HandleFunc("/api/nextdate", handlers.GetNextDate).Methods("GET")
	router.HandleFunc("/api/task", handlers.Auth(handlers.AddNewTask)).Methods("POST")
	router.HandleFunc("/api/tasks", handlers.Auth(handlers.GetAllTasks)).Methods("GET")
	router.HandleFunc("/api/task", handlers.Auth(handlers.GetTasksById)).Methods("GET")
	router.HandleFunc("/api/task", handlers.Auth(handlers.UpdateTask)).Methods("PUT")
	router.HandleFunc("/api/task/done", handlers.Auth(handlers.DoneTask)).Methods("POST")
	router.HandleFunc("/api/task", handlers.Auth(handlers.DeleteTask)).Methods("DELETE")
	router.HandleFunc("/api/signin", handlers.AuthorizationGetToken).Methods("POST")
}
