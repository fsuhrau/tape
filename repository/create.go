package repository

import (
	"os"
	"path/filepath"
)

func Create(repositoryURL string) error {
	if _, err := os.Stat(Dir); !os.IsNotExist(err) {
		return ErrTapeAlreadyInitialized
	}

	if _, err := os.Stat(filepath.Join(Dir, Config)); !os.IsNotExist(err) {
		return ErrTapeAlreadyInitialized
	}

	if err := os.MkdirAll(Dir, 0755); err != nil {
		return err
	}

	repo := Repository{
		Version: Version,
		URL:     repositoryURL,
	}

	return repo.Save()
}
