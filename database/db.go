package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mdapathy/url-shortener/config"
	"log"
)

func NewDBConfig(file string) *sql.DB {
	configuration := config.New(file)
	conn, err := Open(configuration.ConnectionUrl())
	if err != nil {
		log.Fatalf("could not connect to datbase %s", err)
	}
	return conn

}

// Creates a connection to database based on the given connection string
func Open(c string) (*sql.DB, error) {
	return sql.Open("postgres", c)

}
