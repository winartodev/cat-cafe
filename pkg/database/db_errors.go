package database

import (
	"log"
	"strings"
)

func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()

	log.Println(msg)

	return strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "duplicate key")
}
