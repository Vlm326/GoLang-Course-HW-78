package domain

import "errors"

var (
	ErrInvalidRepository  = errors.New("invalid repository")
	ErrRepositoryNotFound = errors.New("repository not found")
)
