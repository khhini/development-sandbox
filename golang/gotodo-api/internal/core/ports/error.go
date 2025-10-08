package ports

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrDuplicateID = errors.New("ID already exist")
)
