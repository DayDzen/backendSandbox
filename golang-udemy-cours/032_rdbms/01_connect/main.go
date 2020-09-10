package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "awsuser1:mypassword1@tcp(mydbinstance1.c2atw9g0okil.us-east-2.rds.amazonaws.com:3306)/test02?charset=utf8")
	check(err)
	defer db.Close()

	err = db.Ping()
	check(err)

	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	_, err = io.WriteString(w, "All is good")
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
