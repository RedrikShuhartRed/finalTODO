package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"

	"github.com/RedrikShuhartRed/finalTODO/config"
	"github.com/RedrikShuhartRed/finalTODO/models"
)

const (
	timeFormat = "02.01.2006"
	limit      = 20
)

type Storage struct {
	db *sql.DB
}

func CheckExistencesShedulerDB(cfg *config.Config) (bool, string) {
	dbFile := cfg.DbFile
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	return install, dbFile
}

func ConnectDB(cfg *config.Config) (*Storage, error) {
	install, dbFile := CheckExistencesShedulerDB(cfg)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Printf("Error connect to DB, %v", err)
		return nil, err
	}

	switch install {
	case true:
		_, err := db.Exec(`CREATE TABLE scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date CHAR(8), 
	title VARCHAR(256) NOT NULL DEFAULT "", comment TEXT , repeat VARCHAR(128) DEFAULT "" )`)
		if err != nil {
			log.Printf("Error create table in DB, %v", err)
			return nil, err
		}

		_, err = db.Exec(`CREATE INDEX idx_date ON scheduler (date);`)
		if err != nil {
			log.Printf("Error create index for date in DB, %v", err)
			return nil, err
		}
		log.Println("Database created successfully.")

	case false:
		log.Println("Database already exists.")
	}

	return &Storage{db: db}, nil
}

func (s Storage) CloseDB() {
	s.db.Close()
}

func (s *Storage) AddNewTask(task models.Task) (int64, error) {

	res, err := s.db.Exec(`INSERT INTO scheduler (title, date, comment, repeat) VALUES (:title, :date, :comment, :repeat)`,
		sql.Named("title", task.Title),
		sql.Named("date", task.Date),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		log.Printf("error insert into scheduler, %v", err)
		return 0, err

	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("error insert into scheduler, %v", err)
		return 0, err
	}
	return lastId, nil
}

func (s Storage) GetAllTasksWithoutSearch() ([]models.Task, error) {
	tasks := []models.Task{}
	rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit",
		sql.Named("limit", limit))
	if err != nil {
		log.Printf("error reading data from database: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		tasks = append(tasks, task)
	}
	if err != nil {
		log.Printf("error Scan data in Task: %v", err)

		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s Storage) GetAllTasksWithStringSearch(search string) ([]models.Task, error) {
	tasks := []models.Task{}
	rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :title OR comment LIKE :comment ORDER BY date LIMIT :limit",
		sql.Named("title", "%"+search+"%"),
		sql.Named("comment", "%"+search+"%"),
		sql.Named("limit", limit),
	)
	if err != nil {
		log.Printf("error reading data from database: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task

		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		tasks = append(tasks, task)
	}
	if err != nil {
		log.Printf("error Scan data in Task: %v", err)

		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s Storage) GetAllTasksWithDateSearch(search string) ([]models.Task, error) {
	tasks := []models.Task{}
	searchdate, err := time.Parse(timeFormat, search)
	if err != nil {
		return nil, err
	}
	correctsearchdate := searchdate.Format("20060102")
	rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date", sql.Named("date", correctsearchdate))
	if err != nil {
		log.Printf("error reading data from database: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var task models.Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		tasks = append(tasks, task)
	}
	if err != nil {
		log.Printf("error Scan data in Task: %v", err)

		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s Storage) GetTasksById(id string) (*models.Task, error) {
	task := &models.Task{}
	row := s.db.QueryRow("SELECT * FROM scheduler WHERE id = :id;",
		sql.Named("id", id),
	)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s Storage) UpdateTask(task models.Task) (int64, error) {
	result, err := s.db.Exec(`UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`,
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
		return 0, err
	}
	if rowsAffected == 0 {
		log.Printf("no rows affected: %v", err)
		return 0, err
	}
	return rowsAffected, nil
}

func (s Storage) DeleteTask(id string) (int64, error) {
	result, err := s.db.Exec(`DELETE FROM scheduler WHERE id = :id`,
		sql.Named("id", id))
	if err != nil {
		log.Printf("error delete task")

		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("error checking rows affected: %v", err)
		return 0, err
	}
	if rowsAffected == 0 {
		log.Printf("task not found for id: %v", id)
		return 0, err
	}
	return rowsAffected, err
}
