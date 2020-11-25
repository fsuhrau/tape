package repository

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func (r *Repository) Add(name, url string) error {

	var dependency *Dependency

	for i := range r.Dependencies {
		if r.Dependencies[i].Name == name {
			dependency = &r.Dependencies[i]
			break
		}
	}

	if dependency != nil {
		return DependencyAlreadyExists
	} else {
		dependency = &Dependency{
			Name: name,
			URL:  url,
			Type: DeterminateType(url),
		}

		tapeTmpDir := filepath.Join(os.TempDir(), "tape", name)
		os.MkdirAll(tapeTmpDir, os.ModePerm)

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
	}

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
		return DependencyNotFound
	} else {

		dependency.URL = url
		dependency.Type = DeterminateType(url)

		tapeTmpDir := filepath.Join(os.TempDir(), "tape", name)
		os.MkdirAll(tapeTmpDir, os.ModePerm)

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
	}

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
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", resp.Status)
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func FileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// hacky hacky we expect archives to be directories...
func DeterminateType(path string) Type {

	if strings.HasSuffix(path, ".zip") {
		return Directory
	}

	if strings.HasSuffix(path, ".tar") {
		return Directory
	}

	if strings.HasSuffix(path, ".tar.gz") {
		return Directory
	}

	return Executable
}
