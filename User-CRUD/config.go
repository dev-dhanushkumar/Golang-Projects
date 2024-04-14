package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func connecDatabase() {
	db, err = sql.Open("mysql", "root:123hitesh@tcp(127.0.0.1:3306)/user")
	if err != nil {
		panic(err)
	}
	fmt.Println("Database Connected!")
}
