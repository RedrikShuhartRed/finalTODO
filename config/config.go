package config

import (
	"errors"
	"os"
)

type Config struct {
	Port         string
	DbFile       string
	Password     string
	PasswordSalt string
	TokenSalt    string
}

func NewConfig() (*Config, error) {

	port := os.Getenv("TODO_PORT")
	if port == "" {
		return nil, errors.New("TODO_PORT environment variable required but not set")
	}

	dbfile := os.Getenv("TODO_DBFILE")
	if dbfile == "" {
		return nil, errors.New("TODO_DBFILE environment variable required but not set")
	}

	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		return nil, errors.New("TODO_PASSWORD environment variable required but not set")
	}

	passwordsalt := os.Getenv("TODO_PASSWORDSALT")
	if passwordsalt == "" {
		return nil, errors.New("TODO_PASSWORDSALT environment variable required but not set")
	}

	tokensalt := os.Getenv("TODO_TOKENSALT")
	if tokensalt == "" {
		return nil, errors.New("TODO_TOKENSALT environment variable required but not set")
	}

	return &Config{
		Port:         port,
		DbFile:       dbfile,
		Password:     password,
		PasswordSalt: passwordsalt,
		TokenSalt:    tokensalt,
	}, nil
}
