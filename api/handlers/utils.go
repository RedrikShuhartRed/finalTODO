package handlers

import (
	"errors"
	"log"
	"strconv"
)

var errEmptyId = errors.New("error Decode request body, Task id is empty")

func CheckId(id string) error {
	if id == "" {
		log.Printf("error get id task: id == \"\"")
		return errEmptyId
	}
	_, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("error get id task, id not int: %v", err)
		return err
	}
	return nil
}
