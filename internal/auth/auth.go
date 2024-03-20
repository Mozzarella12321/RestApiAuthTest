package auth

import (
	"errors"
	"fmt"
	"restapiauthtest/internal/storage"

	"github.com/andskur/argon2-hashing"
	"github.com/google/uuid"
)

// var (
// 	ErrTooShortPassword = errors.New("password is too short")
// 	ErrTooShortLogin    = errors.New("login is too short")
// )

type Storage interface {
	SaveNewUser(login string, hash []byte) error
	GetUserData(login string) ([]byte, error)
	ExecuteQuery(query string, args ...interface{}) error
	CreateToken(token uuid.UUID) error
}

func Login(s Storage, login, password string) (token uuid.UUID, err error) {
	const op = "auth.Login"

	hash, err := s.GetUserData(login)
	if errors.Is(err, storage.ErrNotFound) {
		return uuid.Nil, fmt.Errorf("%s: user not found", op)
	}

	if errors.Is(err, storage.ErrBlocked) {
		return uuid.Nil, fmt.Errorf("%s: user is blocked", op)
	}

	if err != nil {
		return uuid.Nil, fmt.Errorf("could not get user data: %s: %w", op, err)
	}

	err = verifyPassword(password, hash)
	if errors.Is(err, storage.ErrWrongPassword) {

		query := `
			UPDATE users
			SET unsuccessful_logins = unsuccessful_logins + 1
			WHERE login = $1
    	`
		s.ExecuteQuery(query, login)

		return uuid.Nil, fmt.Errorf("%s: wrong password", op)
	}

	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	token = uuid.New()
	err = s.CreateToken(token)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	resetUnsuccessfulLogins(s, login)

	return token, nil
}
func RegisterNewUser(s Storage, login, password string) error {
	const op = "auth.RegisterNewUser"
	// if len(password) < 8 {
	// 	return ErrTooShortPassword
	// }
	// if len(login) < 4 {
	// 	return ErrTooShortLogin
	// }
	hash, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.SaveNewUser(login, hash)
	if errors.Is(err, storage.ErrExists) {
		return fmt.Errorf("%s: user already exists", op)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func verifyPassword(password string, hash []byte) error {
	const op = "auth.VerifyPassword"

	err := argon2.CompareHashAndPassword(hash, []byte(password))
	if err == argon2.ErrMismatchedHashAndPassword {
		return storage.ErrWrongPassword
	}

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
func hashPassword(password string) ([]byte, error) {
	const op = "auth.HashPassword"

	key, err := argon2.GenerateFromPassword([]byte(password), argon2.DefaultParams)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return key, err
}

func resetUnsuccessfulLogins(s Storage, login string) {

	query := `
		UPDATE users
		SET unsuccessful_logins = 0
		WHERE login = $1
	`

	s.ExecuteQuery(query, login)
}
