package auth

import (
	"fmt"

	"github.com/andskur/argon2-hashing"
	"github.com/google/uuid"
)

type Storage interface {
	SaveNewUser(login string, hash []byte) error
	GetUserData(login string) ([]byte, error)
	ExecuteQuery(query string, args ...interface{}) error
}

func Login(s Storage, login, password string) (token uuid.UUID, err error) {
	const op = "auth.Login"

	token = uuid.New()
	hash, err := s.GetUserData(login)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not get user data: %s: %w", op, err)
	}
	err = VerifyPassword(password, hash)
	if err == argon2.ErrMismatchedHashAndPassword {
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
	resetUnsuccessfulLogins(s, login)
	return token, nil
}
func RegisterNewUser(s Storage, login, password string) error {
	const op = "auth.RegisterNewUser"
	hash, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	err = s.SaveNewUser(login, hash)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func VerifyPassword(password string, hash []byte) error {
	const op = "auth.VerifyPassword"
	err := argon2.CompareHashAndPassword(hash, []byte(password))
	if err == argon2.ErrMismatchedHashAndPassword {
		return argon2.ErrMismatchedHashAndPassword
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
func HashPassword(password string) ([]byte, error) {
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