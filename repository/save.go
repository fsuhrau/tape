package repository

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

func (r *Repository) Save() error {
	data, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(DIR, "config"), data, 0644)
}
