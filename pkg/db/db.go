package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Configuration struct for getting data from config.json file.
type Configuration struct {
	MYSQL_DATABASE string
	MYSQL_USER     string
	MYSQL_PASSWORD string
	MYSQL_SERVER   string
	MYSQL_PORT     string
}

var Db *sql.DB

func InitDB() {

	// reading json file
	configfile, err := os.Open("config.json")
	if err != nil {
		log.Panic("Error in opening config.json file ", err)
	}
	var config Configuration

	//decoding json file and checking for error
	decoder := json.NewDecoder(configfile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Panic("Error in decoding config.json file ", err)
	}
	// closing configuration file at end of program
	defer configfile.Close()

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true",
		config.MYSQL_USER, config.MYSQL_PASSWORD, config.MYSQL_SERVER, config.MYSQL_PORT, config.MYSQL_DATABASE)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Panic("Error in db connection ", err)
	}

	if err = db.Ping(); err != nil {
		log.Panic("Error in db ping ", err)
	}
	Db = db
}
