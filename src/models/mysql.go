package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/lucian0ramos/image-golang/src/config"

	"database/sql"
	"time"
)

var db *sql.DB

func init() {
	CreateConnection()
}

// CreateConnection create or return existent connection
func CreateConnection() {
	if GetConnection() != nil {
		return
	}

	if connection, err := sql.Open("mysql", config.URLDatabase()); err != nil {
		panic(err)
	} else {
		connection.SetConnMaxLifetime(time.Second * 300)
		connection.SetMaxOpenConns(50)
		connection.SetMaxIdleConns(50)
		db = connection
	}
}

// GetConnection return connection
func GetConnection() *sql.DB {
	return db
}

// CloseConnection close connection
func CloseConnection() {
	db.Close()
}
