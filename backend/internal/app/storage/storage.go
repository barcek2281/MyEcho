package storage

import (
	"database/sql"
	_ "github.com/lib/pq" // driver
	"log"
)

type Storage struct {
	DatabaseURL string
	db          *sql.DB
}

// New Config
func New(databaseURL string) *Storage {

	return &Storage{DatabaseURL: databaseURL}
}

//CREATE TABLE users (
//id bigserial not null primary key,
//email varchar not null unique,
//login varchar not null,
//password varchar not null
//);

func (s *Storage) Open() error {
	db, err := sql.Open("postgres", s.DatabaseURL)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
		return err
	}

	s.db = db
	return nil
}

func (s *Storage) Close() {
	s.db.Close()
}
