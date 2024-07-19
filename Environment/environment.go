package environment

import (
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

func LoadEnvPort() (port string) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error load .env: %s", err)
	}
	port = os.Getenv("TODO_PORT")
	if port == "" {
		log.Printf("TODO_PORT environment variable not found, default is used")
		port = "7540"
	}
	log.Printf("environment variable is used: port = %s", port)
	return
}
func LoadEnvPortDBFile() (dbFile string) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error load .env: %s", err)
	}
	dbFile = os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		log.Printf("dbFile environment variable not found, default is used")
		dbFile = "../cmd/scheduler.db"
	}
	log.Printf("environment variable is used: dbFile = %s", dbFile)
	return dbFile
}
