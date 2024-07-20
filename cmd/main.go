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
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("error load .env: %s", err)
	}
	port := environment.LoadEnvPort()
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
	r.PathPrefix("/").Handler(FileServer)
	log.Printf("Starting server on port %s...\n", port)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Printf("Error starting Server, %v", err)
	}

}
