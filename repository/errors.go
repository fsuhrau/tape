package repository

import "fmt"

var (
	MissingParameter        = fmt.Errorf("missing parameter")
	TapeAlreadyInitialized  = fmt.Errorf("tape is already initialized")
	TapeNotInitialized      = fmt.Errorf("tape is not initialized")
	DependencyAlreadyExists = fmt.Errorf("dependency already exists try update instead")
	DependencyNotFound      = fmt.Errorf("dependency not found add instead")
	FileMistmatch           = fmt.Errorf("file mismatch")
)
