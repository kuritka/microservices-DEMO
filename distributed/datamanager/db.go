package datamanager


import (
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB


//this function is runned before rest of the functionality is called
func init(){
	var err error
	db, err = sql.Open("postgres", "postgres://test:password@localhost/distributed?sslmode=disable")

	if err != nil {
		panic(err.Error())
	}
}
