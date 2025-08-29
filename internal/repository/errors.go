package repository

import "errors"

var (
	ErrDatabaseDSNNotFound = errors.New("database DSN not found in configuration")
	ErrRecordNotFound      = errors.New("record not found")
)