package repository

import (
	"os"
	"path/filepath"
)

func Create(repositoryURL string) error {
	if _, err := os.Stat(DIR); !os.IsNotExist(err) {
		return TapeAlreadyInitialized
	}

	if _, err := os.Stat(filepath.Join(DIR, CONFIG)); !os.IsNotExist(err) {
		return TapeAlreadyInitialized
	}

	if err := os.MkdirAll(DIR, os.ModePerm); err != nil {
		return err
	}

	repo := Repository{
		Version: VERSION,
		URL:     repositoryURL,
	}

	return repo.Save()
}
