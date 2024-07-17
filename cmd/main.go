package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/RedrikShuhartRed/finalTODO/db"
	"github.com/RedrikShuhartRed/finalTODO/handlers"
	"github.com/gorilla/mux"
)

func CreatePort() (port string) {
	port = os.Getenv("TODO_PORT")
	if port == "" {
		log.Printf("TODO_PORT environment variable not found, default port 7540 is used")
		port = "7540"
	}
	port = ":" + port

	return
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	err := db.ConnectDB()
	if err != nil {
		log.Printf("Error connect DB, %v", err)
	}
	dbs := db.GetDB()
	defer db.CloseDB(dbs)
	port := CreatePort()

	webDir, err := filepath.Abs("../web")
	if err != nil {
		log.Fatalf("Failed to get absolute path for web directory: %v", err)
	}

	r := mux.NewRouter()
	FileServer := http.FileServer(http.Dir(webDir))

	r.HandleFunc("/api/nextdate", handlers.GetNextDate).Methods("GET")
	r.HandleFunc("/api/task", handlers.AddNewTask).Methods("POST")
	r.HandleFunc("/api/tasks", handlers.GetAllTasks).Methods("GET")
	r.PathPrefix("/").Handler(FileServer)
	log.Printf("Starting server on port %s...\n", port)

	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Printf("Error starting Server, %v", err)
	}

}
