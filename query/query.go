package query

const (
	CreateTable = `CREATE TABLE scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date CHAR(8), 
	title VARCHAR(256) NOT NULL DEFAULT "", comment TEXT , repeat VARCHAR(128) DEFAULT "" )`
	CreateIndexDate = `CREATE INDEX idx_date ON scheduler (date);`
)
