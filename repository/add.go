package repository

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func (r *Repository) Add(name, url string) error {
	for i := range r.Dependencies {
		if r.Dependencies[i].Name == name {
			return ErrDependencyAlreadyExists
		}
	}

	var dependency = &Dependency{
		Name: name,
		URL:  url,
		Type: DeterminateType(url),
	}

	tapeTmpDir := filepath.Join(os.TempDir(), "tape", name)
	if err := os.MkdirAll(tapeTmpDir, 0755); err != nil {
		return err
	}

	tmpDownloadFileName := filepath.Join(tapeTmpDir, name)
	if err := HTTPDownload(url, tmpDownloadFileName); err != nil {
		return err
	}

	hash, err := FileHash(tmpDownloadFileName)
	if err != nil {
		return err
	}

	dependency.Hash = hash

	r.Dependencies = append(r.Dependencies, *dependency)

	return nil
}

func (r *Repository) Update(name, url string) error {
	var dependency *Dependency

	for i := range r.Dependencies {
		if r.Dependencies[i].Name == name {
			dependency = &r.Dependencies[i]
			break
		}
	}

	if dependency == nil {
		return ErrDependencyNotFound
	}

	dependency.URL = url
	dependency.Type = DeterminateType(url)

	tapeTmpDir := filepath.Join(os.TempDir(), "tape", name)
	if err := os.MkdirAll(tapeTmpDir, 0755); err != nil {
		return err
	}

	tmpDownloadFileName := filepath.Join(tapeTmpDir, name)

	if err := HTTPDownload(url, tmpDownloadFileName); err != nil {
		return err
	}

	hash, err := FileHash(tmpDownloadFileName)
	if err != nil {
		return err
	}

	dependency.Hash = hash

	r.Dependencies = append(r.Dependencies, *dependency)

	return nil
}

func HTTPDownload(url, path string) error {
	logrus.Infof("download: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req = AddAuthFromNetrc(url, req)

	// Get the data
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", resp.Status)
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func FileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// hacky hacky we expect archives to be directories...
func DeterminateType(path string) Type {
	switch filepath.Ext(path) {
	case ".zip", ".tar", ".gz":
		return Directory
	default:
		return Executable
	}
}
