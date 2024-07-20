package main

import (
	"log"
	"net/http"
	"path/filepath"

	environment "github.com/RedrikShuhartRed/finalTODO/Environment"
	"github.com/RedrikShuhartRed/finalTODO/api"
	"github.com/RedrikShuhartRed/finalTODO/db"
	"github.com/gorilla/mux"
	"github.com/lpernett/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("error load .env: %s", err)
	}
	port := environment.LoadEnvPort()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	err = db.ConnectDB()
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
	api.RegisterTasksStoreRoutes(r)
	// r.HandleFunc("/api/nextdate", handlers.GetNextDate).Methods("GET")
	// r.HandleFunc("/api/task", handlers.Auth(handlers.AddNewTask)).Methods("POST")
	// r.HandleFunc("/api/tasks", handlers.Auth(handlers.GetAllTasks)).Methods("GET")
	// r.HandleFunc("/api/task", handlers.Auth(handlers.GetTasksById)).Methods("GET")
	// r.HandleFunc("/api/task", handlers.Auth(handlers.UpdateTask)).Methods("PUT")
	// r.HandleFunc("/api/task/done", handlers.Auth(handlers.DoneTask)).Methods("POST")
	// r.HandleFunc("/api/task", handlers.Auth(handlers.DeleteTask)).Methods("DELETE")
	// r.HandleFunc("/api/signin", handlers.AuthorizationGetToken).Methods("POST")
	r.PathPrefix("/").Handler(FileServer)
	log.Printf("Starting server on port %s...\n", port)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Printf("Error starting Server, %v", err)
	}

}
