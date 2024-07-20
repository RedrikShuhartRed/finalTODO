package environment

import (
	"log"
	"os"
)

func LoadEnvPort() (port string) {

	port = os.Getenv("TODO_PORT")
	if port == "" {
		log.Printf("TODO_PORT environment variable not found, default is used")
		port = "7540"
	}
	log.Printf("environment variable is used: port = %s", port)
	return
}
func LoadEnvPortDBFile() (dbFile string) {

	dbFile = os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		log.Printf("dbFile environment variable not found, default is used")
		dbFile = "../cmd/scheduler.db"
	}
	log.Printf("environment variable is used: dbFile = %s", dbFile)
	return dbFile
}

func LoadEnvPassword() (password string) {
	password = os.Getenv("TODO_PASSWORD")

	return
}
func LoadEnvPasswordSalt() (passwordSalt string) {
	passwordSalt = os.Getenv("TODO_PASSWORDSALT")

	return
}
func LoadEnvTokenSalt() (tokedSalt string) {
	tokedSalt = os.Getenv("TODO_TOKENSALT")

	return
}
