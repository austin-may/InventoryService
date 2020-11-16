package datamanager

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

var DbConn *sql.DB

func init() {
	var err error
	connString := fmt.Sprintf("sqlserver://austinmay@localhost/SQLEXPRESS?database=VitaminDB")
	DbConn, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
	}
	DbConn.SetMaxOpenConns(4)
	DbConn.SetMaxIdleConns(4)
	DbConn.SetConnMaxLifetime(1 * time.Millisecond)
	log.Printf("Connected!\n")

}
