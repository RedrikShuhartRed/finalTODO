package db

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	environment "github.com/RedrikShuhartRed/finalTODO/Environment"
	"github.com/RedrikShuhartRed/finalTODO/models"
	"github.com/RedrikShuhartRed/finalTODO/query"
	_ "modernc.org/sqlite"
)

var dbs *sql.DB

func CheckExistencesShedulerDB() (bool, string) {
	dbFile := environment.LoadEnvPortDBFile()
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

func AddNewTask(task *models.Task) (int64, error) {
	dbs := GetDB()
	res, err := dbs.Exec(query.AddNewTask,
		sql.Named("title", task.Title),
		sql.Named("date", task.Date),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		log.Printf("error insert into scheduler, %v", err)
		return 0, err

	}
	lastId, err := (res.LastInsertId())
	if err != nil {
		log.Printf("error insert into scheduler, %v", err)
		return 0, err
	}

	return lastId, nil

}

func GetAllTasks(search string) ([]models.Task, error) {
	dbs := GetDB()
	tasks := []models.Task{}
	limit := 20
	if search == "" {
		limit = 50
	}
	var rows *sql.Rows
	var err error
	timeFormat := "02.01.2006"

	if search == "" {
		rows, err = dbs.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit",
			sql.Named("limit", limit))
		if err != nil {
			log.Printf("error reading data from database: %v", err)
			return nil, err
		}
	} else {
		searchdate, err := time.Parse(timeFormat, search)
		if err != nil {
			rows, err = dbs.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :title OR comment LIKE :comment ORDER BY date LIMIT :limit",
				sql.Named("title", "%"+search+"%"),
				sql.Named("comment", "%"+search+"%"),
				sql.Named("limit", limit),
			)

		} else {
			correctsearchdate := searchdate.Format("20060102")
			rows, err = dbs.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date", sql.Named("date", correctsearchdate))
		}
		if err != nil {
			log.Printf("error reading data from database: %v", err)
			return nil, err
		}

	}
	defer rows.Close()

	for rows.Next() {
		task := models.Task{}
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		tasks = append(tasks, task)
	}
	if err != nil {
		log.Printf("error Scan data in Task: %v", err)

		return nil, err
	}
	return tasks, nil
}

func GetTasksById(id string) (*models.Task, error) {
	dbs := GetDB()
	task := &models.Task{}
	row := dbs.QueryRow("SELECT * FROM scheduler WHERE id = :id;",
		sql.Named("id", id),
	)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return task, nil
}

func UpdateTask(task *models.Task) (int64, error) {
	dbs := GetDB()

	result, err := dbs.Exec(`UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`,
		sql.Named("title", task.Title),
		sql.Named("date", task.Date),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)
	if err != nil {
		log.Printf("error insert into scheduler, %v", err)
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("error getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		log.Printf("no rows affected: %v", err)
	}
	return rowsAffected, nil
}

func DeleteTask(id string) (int64, error) {
	result, err := dbs.Exec(`DELETE FROM scheduler WHERE id = :id`,
		sql.Named("id", id))
	if err != nil {
		log.Printf("error delete task")

		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		log.Printf("task not found for id: %v", id)
	}
	return rowsAffected, err
}
