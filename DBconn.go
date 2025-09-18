package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func initDB() *sql.DB {
	var err error
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/challange2016?charset=utf8&parseTime=true")
	if err != nil {
		panic(err)
	}
	return db
}
