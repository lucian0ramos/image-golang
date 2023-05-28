package config

import (
	"fmt"
	"os"
)

// Config Config
type Config interface {
	url() string
}

type DatabaseConfig struct {
	username string
	password string
	host     string
	port     int
	database string
	debug    bool
}

var database *DatabaseConfig

func init() {
	database = &DatabaseConfig{}
	database.username = os.Getenv("DB_USER")
	database.password = os.Getenv("DB_PSWD")
	database.database = os.Getenv("DB_NAME")
	database.host = os.Getenv("DB_HOST")
	database.port = 3306
	database.debug = true
}

func (db *DatabaseConfig) url() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true", db.username, db.password, db.host, db.port, db.database)
}

// URLDatabase URLDatabase
func URLDatabase() string {
	return database.url()
}
