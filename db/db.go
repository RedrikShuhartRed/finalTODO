package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/RedrikShuhartRed/finalTODO/query"
	_ "modernc.org/sqlite"
)

var dbs *sql.DB

func CheckExistencesShedulerDB() (bool, string) {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {

		appPath, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current working directory, %v", err)
		}

		dbFile = filepath.Join(appPath, "scheduler.db")
		log.Printf("TODO_DBFILE environment variable not found, using %s in project root directory", dbFile)
	}
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	return install, dbFile

}

func ConnectDB() error {

	install, dbFile := CheckExistencesShedulerDB()

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Printf("Error connect to DB, %v", err)
		return err
	}

	switch install {
	case true:
		_, err := db.Exec(query.CreateTable)
		if err != nil {
			log.Printf("Error create table in DB, %v", err)
			return err
		}

		_, err = db.Exec(query.CreateIndexDate)
		if err != nil {
			log.Printf("Error create index for date in DB, %v", err)
			return err
		}
		log.Println("Database created successfully.")

	case false:
		log.Println("Database already exists.")
	}
	dbs = db
	return nil
}

func GetDB() *sql.DB {
	return dbs
}

func CloseDB(db *sql.DB) error {
	err := db.Close()
	if err != nil {
		log.Printf("Error closing connection to database: %s", err)
		return err
	}
	return nil
}
