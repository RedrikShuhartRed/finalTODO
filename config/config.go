package config

import (
	"log"
	"os"
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
		port = "7540"
		log.Printf("TODO_PORT environment variable required but not set, used default")
	}

	dbfile := os.Getenv("TODO_DBFILE")
	if dbfile == "" {
		dbfile = "./scheduler.db"
		log.Printf("TODO_DBFILE environment variable required but not set, used default")
	}

	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		password = "myPassword"
		log.Printf("TODO_PASSWORD environment variable required but not set, used default")
	}

	passwordsalt := os.Getenv("TODO_PASSWORDSALT")
	if passwordsalt == "" {
		passwordsalt = "kl4509dafh43589whfh"
		log.Printf("TODO_PASSWORDSALT environment variable required but not se, used defaultt")
	}

	tokensalt := os.Getenv("TODO_TOKENSALT")
	if tokensalt == "" {
		tokensalt = "klajglk54adgagsd"
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
