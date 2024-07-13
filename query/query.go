package query

const (
	CreateTable = `CREATE TABLE scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date CHAR(8) NOT NULL, 
	title VARCHAR(256) NOT NULL , comment TEXT , repeat VARCHAR(128) )`
	CreateIndexDate = `CREATE INDEX idx_date ON scheduler (date);`
)
