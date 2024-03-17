package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "postgresql.New"
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("could not open storage: %s: %w", op, err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login TEXT,
		password TEXT,
		salt TEXT,
		unsuccessful_logins INT
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("could not create table users: %s: %w", op, err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS sessions (
		id SERIAL PRIMARY KEY,
		token INT,
		lifetime INT
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("could not create table sessions: %s: %w", op, err)
	}

	return &Storage{db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) InsertLogin(login string) error {
	const op = "postgresql.InsertLogin"
	_, err := s.db.Exec(`INSERT INTO users (login) VALUES ($1)`, login)
	if err != nil {
		return fmt.Errorf("could not insert login: %s: %w", op, err)
	}

	return nil
}
