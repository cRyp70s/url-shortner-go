package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	db *sql.DB
}

func (sqlitedb *SQLiteDB) Execute(statement string, values ...interface{}) (sql.Result, error) {
	stmt, err := sqlitedb.db.Prepare(statement)
	if err != nil {
		handleError(err)
	}
	res, err := stmt.Exec(values...)
	return res, err
}

func (sqlitedb *SQLiteDB) Query(query string, values ...interface{}) (*sql.Rows, error) {
	rows, err := sqlitedb.db.Query(query)
	return rows, err
}

func (sqlitedb *SQLiteDB) Close() {
	sqlitedb.db.Close()
}

func handleError(err error) {
	if err != nil {
		log.Printf("An Error Occured: %s", err)
	}
}

func NewDB(filepath string) *SQLiteDB {
	db, err := sql.Open("sqlite3", filepath)
	handleError(err)
	sqlitedb := SQLiteDB{db}
	sqlitedb.Execute("CREATE TABLE IF NOT EXISTS `urlmaps` ( `id` INTEGER PRIMARY KEY AUTOINCREMENT, `original_url` VARCHAR(256), `short_url` VARCHAR(256));")
	return &sqlitedb

}
