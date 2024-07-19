package main

import (
	"log"
	"net/http"
	"path/filepath"

	environment "github.com/RedrikShuhartRed/finalTODO/Environment"
	"github.com/RedrikShuhartRed/finalTODO/db"
	"github.com/RedrikShuhartRed/finalTODO/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := environment.LoadEnvPort()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	err := db.ConnectDB()
	if err != nil {
		log.Printf("Error connect DB, %v", err)
	}
	dbs := db.GetDB()
	defer db.CloseDB(dbs)

	webDir, err := filepath.Abs("../web")
	if err != nil {
		log.Fatalf("Failed to get absolute path for web directory: %v", err)
	}

	r := mux.NewRouter()
	FileServer := http.FileServer(http.Dir(webDir))

	r.HandleFunc("/api/nextdate", handlers.GetNextDate).Methods("GET")
	r.HandleFunc("/api/task", handlers.AddNewTask).Methods("POST")
	r.HandleFunc("/api/tasks", handlers.GetAllTasks).Methods("GET")
	r.HandleFunc("/api/task", handlers.GetTasksById).Methods("GET")
	r.HandleFunc("/api/task", handlers.UpdateTask).Methods("PUT")
	r.HandleFunc("/api/task/done", handlers.DoneTask).Methods("POST")
	r.HandleFunc("/api/task", handlers.DeleteTask).Methods("DELETE")
	//r.HandleFunc("/api/task", handlers.AuthorizationGetToken).Methods("POST")
	r.PathPrefix("/").Handler(FileServer)
	log.Printf("Starting server on port %s...\n", port)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Printf("Error starting Server, %v", err)
	}

}
