package config

import (
	"log"
	"os"
)

const (
	defaultPort         = "7540"
	defaultDbfile       = "./scheduler.db"
	defaultPassword     = "myPassword"
	defaultPasswordsalt = "kl4509dafh43589whfh"
	defaultTokensalt    = "klajglk54adgagsd"
)

type Config struct {
	Port         string
	DbFile       string
	Password     string
	PasswordSalt string
	TokenSalt    string
}

func NewConfig() *Config {

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
		log.Printf("TODO_PORT environment variable required but not set, used default")
	}

	dbfile := os.Getenv("TODO_DBFILE")
	if dbfile == "" {
		dbfile = defaultDbfile
		log.Printf("TODO_DBFILE environment variable required but not set, used default")
	}

	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		password = defaultPassword
		log.Printf("TODO_PASSWORD environment variable required but not set, used default")
	}

	passwordsalt := os.Getenv("TODO_PASSWORDSALT")
	if passwordsalt == "" {
		passwordsalt = defaultPasswordsalt
		log.Printf("TODO_PASSWORDSALT environment variable required but not se, used defaultt")
	}

	tokensalt := os.Getenv("TODO_TOKENSALT")
	if tokensalt == "" {
		tokensalt = defaultTokensalt
		log.Printf("TODO_TOKENSALT environment variable required but not set, used default")
	}

	return &Config{
		Port:         port,
		DbFile:       dbfile,
		Password:     password,
		PasswordSalt: passwordsalt,
		TokenSalt:    tokensalt,
	}
}
