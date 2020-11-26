package repository

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func TapeHome() string {
	homeDir, _ := homedir.Dir()
	return filepath.Join(homeDir, ".tape")
}

func LinksDir() string {
	return filepath.Join(Dir, "links")
}

func (r *Repository) Link() error {
	tapeSharedRepoDir := TapeHome()

	if _, err := os.Stat(tapeSharedRepoDir); os.IsNotExist(err) {
		if err := os.MkdirAll(tapeSharedRepoDir, 0755); err != nil {
			return err
		}
	}

	for _, dep := range r.Dependencies {
		needsDownload := false
		depDir := filepath.Join(tapeSharedRepoDir, dep.Hash)
		if _, err := os.Stat(depDir); os.IsNotExist(err) {
			needsDownload = true
		}

		tapeRepoDep := filepath.Join(depDir, dep.Name)
		if _, err := os.Stat(tapeRepoDep); os.IsNotExist(err) {
			needsDownload = true
		}

		if needsDownload {
			tapeTmpDir := filepath.Join(os.TempDir(), "tape", dep.Name)
			if err := os.MkdirAll(tapeTmpDir, 0755); err != nil {
				return err
			}

			tmpDownloadFileName := filepath.Join(tapeTmpDir, dep.Name)
			if err := HTTPDownload(dep.URL, tmpDownloadFileName); err != nil {
				return err
			}

			hash, err := FileHash(tmpDownloadFileName)
			if err != nil {
				return err
			}

			if dep.Hash != hash {
				return ErrFileMismatch
			}

			if err := os.MkdirAll(depDir, 0755); err != nil {
				return err
			}

			switch dep.Type {
			case Directory:
				if err := unzip(tmpDownloadFileName, tapeRepoDep); err != nil {
					return err
				}
			case Executable:
				if err := move(tmpDownloadFileName, tapeRepoDep); err != nil {
					return err
				}
			}
		}

		if err := dep.Link(); err != nil {
			return err
		}
	}
	return nil
}

func (d *Dependency) Link() error {
	repoLinks := LinksDir()
	linkSource := filepath.Join(TapeHome(), d.Hash, d.Name)

	if _, err := os.Stat(repoLinks); os.IsNotExist(err) {
		if err := os.MkdirAll(repoLinks, 0755); err != nil {
			return err
		}
	}
	linkName := filepath.Join(repoLinks, d.Name)
	if _, err := os.Stat(linkName); !os.IsNotExist(err) {
		if err := os.RemoveAll(linkName); err != nil {
			return err
		}
	}

	if err := os.Symlink(linkSource, linkName); err != nil {
		return err
	}

	return nil
}

func (d *Dependency) Unlink() error {
	repoLinks := LinksDir()
	linkName := filepath.Join(repoLinks, d.Name)
	if _, err := os.Stat(linkName); !os.IsNotExist(err) {
		if err := os.RemoveAll(linkName); err != nil {
			return err
		}
	}
	return nil
}

func unzip(zipFile, destination string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = r.Close()
	}()

	for _, f := range r.File {
		fpath := filepath.Join(destination, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(destination)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, 0755); err != nil {
				return err
			}
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		err := func() error {
			var buffer bytes.Buffer
			writer := bufio.NewWriter(&buffer)

			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer func() {
				_ = rc.Close()
			}()

			if _, err := io.Copy(writer, rc); err != nil {
				return err
			}

			if err := ioutil.WriteFile(fpath, buffer.Bytes(), 0755); err != nil {
				return err
			}
			return err
		}()

		if err != nil {
			return err
		}
	}
	return nil
}

func move(tmpFile, destinationFile string) error {
	return os.Rename(tmpFile, destinationFile)
}
