package repository

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Load() (*Repository, error) {
	repositoryConfigFile := filepath.Join(DIR, CONFIG)
	if _, err := os.Stat(repositoryConfigFile); os.IsNotExist(err) {
		return nil, TapeNotInitialized
	}

	repositoryData, err := ioutil.ReadFile(repositoryConfigFile)
	if err != nil {
		return nil, err
	}

	repository := &Repository{}
	if err := json.Unmarshal(repositoryData, repository); err != nil {
		return nil, err
	}

	return repository, nil
}
