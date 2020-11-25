package repository

import (
	"errors"
)

var (
	ErrMissingParameter        = errors.New("missing parameter")
	ErrTapeAlreadyInitialized  = errors.New("tape is already initialized")
	ErrTapeNotInitialized      = errors.New("tape is not initialized")
	ErrDependencyAlreadyExists = errors.New("dependency already exists try update instead")
	ErrDependencyNotFound      = errors.New("dependency not found add instead")
	ErrFileMismatch            = errors.New("file mismatch")
)
