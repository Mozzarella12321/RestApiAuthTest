package postgresql

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"restapiauthtest/internal/storage"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
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

	_, err = postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not get driver: %s: %w", op, err)
	}

	migPath, err := getMigrationsPath()
	if err != nil {
		return nil, fmt.Errorf("could not get migrations path: %s: %w", op, err)
	}

	m, err := migrate.New("file://"+migPath, storagePath)
	if err != nil {
		return nil, fmt.Errorf("could not create migration instance: %s: %w", op, err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("error applying migrations: %s: %w", op, err)
	}

	return &Storage{db}, nil
}

// func New(storagePath string) (*Storage, error) {
// 	const op = "postgresql.New"
//
// 	db, err := sql.Open("postgres", storagePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not open storage: %s: %w", op, err)
// 	}
//
// 	_, err = db.Exec(`
// 	CREATE TABLE IF NOT EXISTS users (
// 		id SERIAL PRIMARY KEY,
// 		login TEXT,
// 		hash BYTEA,
// 		unsuccessful_logins INT
// 	);
// 	`)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not create table users: %s: %w", op, err)
// 	}
//
// 	_, err = db.Exec(`
// 	CREATE TABLE IF NOT EXISTS sessions (
// 		id SERIAL PRIMARY KEY,
// 		token INT,
// 		lifetime INT
// 	);
// 	`)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not create table sessions: %s: %w", op, err)
// 	}
//
// 	return &Storage{db}, nil
// }

func (s *Storage) CreateToken(token uuid.UUID) error {
	const op = "CreateToken"

	currentTime := time.Now()

	query := `
    INSERT INTO sessions (token, created_at) VALUES ($1, $2)
    `

	_, err := s.db.Exec(query, token, currentTime)
	if err != nil {
		return fmt.Errorf("failed to create token: %s: %w", op, err)
	}

	return err
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

func (s *Storage) GetToken(token uuid.UUID) error {
	const op = "postgresql.GetToken"
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM sessions WHERE token = $1)"
	err := s.db.QueryRow(query, token).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: could not execute query: %v", op, err)
	}

	if !exists {
		return storage.ErrNotFound
	}

	var created_At time.Time
	err = s.db.QueryRow(`SELECT created_at FROM sessions WHERE token = $1`, token).Scan(&created_At)
	if err != nil {
		return err
	}
	// Check if the difference is less than one hour
	if time.Since(created_At) > time.Hour {
		return storage.ErrExpired
	}
	return nil
}

func getMigrationsPath() (string, error) {
	const op = "postgresql.getMigrationsPath"
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("%s: could not get current working directory: %w", op, err)
	}

	migPath := path.Join(cwd, "..", "..", "..", "internal", "storage", "migrations")

	return migPath, nil
}
