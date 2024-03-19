package postgresql

import (
	"database/sql"
	"fmt"
	"restapiauthtest/internal/storage"

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
		hash BYTEA,
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

func (s *Storage) ExecuteQuery(query string, args ...interface{}) error {
	const op = "postgresql.ExecuteQuery"
	_, err := s.db.Exec(query, args...)

	if err != nil {
		return fmt.Errorf("could not execute query: %s: %w", op, err)
	}
	return nil
}
func (s *Storage) SaveNewUser(login string, hash []byte) error {
	const op = "postgresql.SaveNewUser"

	var count int

	err := s.db.QueryRow(`SELECT COUNT(*) FROM users WHERE login = $1`, login).Scan(&count)
	if err != nil {
		return fmt.Errorf("could not check user existence: %s: %w", op, err)
	}

	if count > 0 {
		return storage.ErrExists
	}

	_, err = s.db.Exec(`INSERT INTO users (login, hash, unsuccessful_logins) VALUES ($1, $2, $3)`, login, hash, 0)
	if err != nil {
		return fmt.Errorf("could not save new user: %s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUserData(login string) ([]byte, error) {
	const op = "postgresql.getUserData"

	var hash []byte
	var unsuccessful_logins int
	err := s.db.QueryRow(`SELECT hash, unsuccessful_logins FROM users WHERE login = $1`, login).Scan(&hash, &unsuccessful_logins)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not get user data: %s: %w", op, err)
	}

	if unsuccessful_logins >= 5 {
		return nil, storage.ErrBlocked
	}

	return hash, nil
}
