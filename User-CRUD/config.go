package main

import (
	"database/sql"
	"fmt"
)

var db *sql.DB
var err error

func connecDatabase() {
	db, err = sql.Open("mysql", "root:<password>@tcp(127.0.0.1:3306)/user")
	fmt.Println("Database Connected!")
}
