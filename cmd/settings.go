package cmd

import (
	"encoding/json"
	"os"
	"path"
)

// loadLastModeFromSettings loads the last UI selected mode from
// $HOME/.todo/settings.json. It returns the mode or an error.
func loadLastModeFromSettings(home string) (string, error) {
	fname := path.Join(home, ".todo", "settings.json")
	data, err := os.ReadFile(fname)
	if err != nil {
		return "", err
	}
	var s struct {
		Mode string `json:"mode"`
	}
	if err := json.Unmarshal(data, &s); err != nil {
		return "", err
	}
	if s.Mode == "" {
		s.Mode = "main"
	}
	return s.Mode, nil
}

// saveLastModeToSettings writes the provided mode to
// $HOME/.todo/settings.json.
func saveLastModeToSettings(home, mode string) error {
	dir := path.Join(home, ".todo")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	fname := path.Join(dir, "settings.json")
	data, err := json.Marshal(struct {
		Mode string `json:"mode"`
	}{Mode: mode})
	if err != nil {
		return err
	}
	return os.WriteFile(fname, data, 0644)
}
