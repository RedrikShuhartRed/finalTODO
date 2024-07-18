package query

const (
	CreateTable = `CREATE TABLE scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date CHAR(8), 
	title VARCHAR(256) NOT NULL DEFAULT "", comment TEXT , repeat VARCHAR(128) DEFAULT "" )`
	CreateIndexDate = `CREATE INDEX idx_date ON scheduler (date);`
	AddNewTask      = `INSERT INTO scheduler (title, date, comment, repeat) VALUES (:title, :date, :comment, :repeat)`
	UpdateTask      = `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
)
