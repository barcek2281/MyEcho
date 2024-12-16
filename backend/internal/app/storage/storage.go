package storage

import (
	"database/sql"
	_ "github.com/lib/pq" // driver
	"log"
)

type Storage struct {
	DatabaseURL    string
	db             *sql.DB
	userRepository *UserRepository
}

// New Config
func New(databaseURL string) *Storage {

	return &Storage{DatabaseURL: databaseURL}
}

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

func (s *Storage) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}
	s.userRepository = &UserRepository{
		storage: s,
	}
	return s.userRepository
}
