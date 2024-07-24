package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/lpernett/godotenv"

	"github.com/RedrikShuhartRed/finalTODO/api"
	"github.com/RedrikShuhartRed/finalTODO/config"
	"github.com/RedrikShuhartRed/finalTODO/db"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("error load .env: %s", err)
	}
	cfg := config.NewConfig()

	log.Printf("Configuration loaded: %+v\n", cfg)
	port := cfg.Port

	storage, err := db.ConnectDB(cfg)
	if err != nil {
		log.Printf("Error connect DB, %v", err)
	}

	defer storage.CloseDB()
	webDir, err := filepath.Abs("./web")
	if err != nil {
		log.Fatalf("Failed to get absolute path for web directory: %v", err)
	}

	r := mux.NewRouter()
	FileServer := http.FileServer(http.Dir(webDir))
	api.RegisterTasksStoreRoutes(r, storage, cfg)

	r.PathPrefix("/").Handler(FileServer)
	log.Printf("Starting server on port %s...\n", port)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Printf("Error starting Server, %v", err)
	}

}
