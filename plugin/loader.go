package plugin

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mattn/go-zglob"
	"github.com/sirupsen/logrus"
)

type Action struct {
	Parent      string `json:"parent"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Actions []Action

type Plugin struct {
	Executable  string
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Actions     string `json:"actions"`
}

type Loader struct {
	plugins map[string]Plugin
}

func (l *Loader) LoadPlugins() error {

	l.plugins = make(map[string]Plugin)

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	tapePlugins, err := zglob.Glob(filepath.Join(dir, "tape-*"))
	if err != nil {
		return err
	}

	for _, p := range tapePlugins {
		cmd := exec.Command(p, "actions")
		actionsJSON, err := cmd.CombinedOutput()
		if err != nil {
			logrus.Errorf("unsupported plugin: %s", p)
			continue
		}

		plugin := Plugin{Executable: p}
		if err := json.Unmarshal(actionsJSON, &plugin); err != nil {
			logrus.Errorf("unsupported plugin: %s", p)
			continue
		}

		l.plugins[plugin.Name] = plugin
	}

	return nil
}

func (l *Loader) List() {
	for _, v := range l.plugins {
		logrus.Infof("%s(%s) @ %s", v.Name, v.Version, v.Executable)
		logrus.Info(v.Description)
		logrus.Info("")
	}
}
