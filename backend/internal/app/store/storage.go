package storage

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // driver
)

type Storage struct {
	DatabaseURL    string
	db             *sql.DB
	userRepository *UserRepository
	postRepository *PostRepository
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

func (s *Storage) Post() *PostRepository {
	if s.postRepository == nil {
		s.postRepository = &PostRepository{
			storage: s,
		}
	}
	return s.postRepository
}