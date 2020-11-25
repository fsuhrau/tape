package repository

const (
	DIR     = ".tape"
	VERSION = "1.0"
	CONFIG  = "config"
)

type Type int

const (
	Executable Type = iota + 0
	Directory
)

type Repository struct {
	Version      string       `json:"version"`
	URL          string       `json:"url"`
	Dependencies Dependencies `json:"dependencies,omitempty"`
}

type Dependency struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Hash    string `json:"hash"`
	Type    Type   `json:"type"`
	Version string `json:"version,omitempty"`
}

type Dependencies []Dependency
